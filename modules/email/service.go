package email

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/titpetric/platform-app/modules/email/model"
	"github.com/titpetric/platform-app/modules/email/smtp"
)

// EmailStorage interface for dependency injection
type EmailStorage interface {
	Create(ctx context.Context, email *model.Email) error
	Get(ctx context.Context, id string) (*model.Email, error)
	GetPending(ctx context.Context, limit int) ([]model.Email, error)
	Update(ctx context.Context, email *model.Email) error
}

// ServiceOptions configures the email service behavior
type ServiceOptions struct {
	// MaxRetries is the maximum number of retry attempts (default: 3)
	MaxRetries int
	// RetryDuration is the maximum time to attempt retries (default: 1 hour)
	RetryDuration time.Duration
	// TickerInterval is the interval for checking pending emails (default: 5 seconds)
	TickerInterval time.Duration
	// QueueCapacity is the in-memory queue capacity (default: 100)
	QueueCapacity int
	// Logger for structured logging
	Logger *slog.Logger
}

// DefaultServiceOptions returns sensible defaults
func DefaultServiceOptions() ServiceOptions {
	return ServiceOptions{
		MaxRetries:     3,
		RetryDuration:  time.Hour,
		TickerInterval: 5 * time.Second,
		QueueCapacity:  100,
		Logger:         slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

// Service manages email operations
type Service struct {
	emailStorage EmailStorage
	sender       smtp.Sender
	queue        chan *model.Email
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	ticker       *time.Ticker
	options      ServiceOptions
	logger       *slog.Logger
}

// NewService creates a new email service with optional smtp.Config
func NewService(emailStorage EmailStorage, smtpCfg *smtp.Config, opts ...ServiceOptions) *Service {
	// Set up options with defaults
	options := DefaultServiceOptions()
	if len(opts) > 0 {
		// Merge custom options with defaults
		custom := opts[0]
		if custom.MaxRetries != 0 {
			options.MaxRetries = custom.MaxRetries
		}
		if custom.RetryDuration != 0 {
			options.RetryDuration = custom.RetryDuration
		}
		if custom.TickerInterval != 0 {
			options.TickerInterval = custom.TickerInterval
		}
		if custom.QueueCapacity != 0 {
			options.QueueCapacity = custom.QueueCapacity
		}
		if custom.Logger != nil {
			options.Logger = custom.Logger
		}
	}

	// Create sender from config or environment
	var sender smtp.Sender
	if smtpCfg != nil {
		sender = smtp.NewSMTPSender(*smtpCfg)
	} else {
		sender = smtp.NewSMTPSender(smtp.ConfigFromEnv())
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &Service{
		emailStorage: emailStorage,
		sender:       sender,
		queue:        make(chan *model.Email, options.QueueCapacity),
		ctx:          ctx,
		cancel:       cancel,
		ticker:       time.NewTicker(options.TickerInterval),
		options:      options,
		logger:       options.Logger,
	}

	return service
}

// Start begins the email queue consumer and processes pending emails
func (s *Service) Start() {
	s.wg.Add(1)
	go s.processQueue()

	// Process any pending emails on startup
	s.logger.Info("email service started, checking for pending emails")
	s.processPendingEmails()
}

// Stop stops the email queue consumer
func (s *Service) Stop() {
	s.logger.Info("stopping email service")
	s.ticker.Stop()
	close(s.queue)
	s.cancel()
	s.wg.Wait()
}

// AddEmail adds an email to the queue and database
func (s *Service) AddEmail(ctx context.Context, email *model.Email) error {
	// Save to database
	if err := s.emailStorage.Create(ctx, email); err != nil {
		s.logger.Error("failed to create email in database",
			"email_id", email.ID,
			"recipient", email.Recipient,
			"error", err)
		return err
	}

	s.logger.Debug("email created and queued",
		"email_id", email.ID,
		"recipient", email.Recipient,
		"subject", email.Subject)

	// Add to queue
	select {
	case s.queue <- email:
	case <-s.ctx.Done():
		return context.Canceled
	default:
		// Queue full, will be picked up by ticker
		s.logger.Debug("queue full, email will be picked up by ticker",
			"email_id", email.ID)
	}

	return nil
}

// processQueue processes emails from the queue and from the database
func (s *Service) processQueue() {
	defer s.wg.Done()

	for {
		select {
		case email, ok := <-s.queue:
			if !ok {
				return
			}
			s.sendEmail(email)

		case <-s.ticker.C:
			s.processPendingEmails()

		case <-s.ctx.Done():
			return
		}
	}
}

// sendEmail sends a single email with retry logic
func (s *Service) sendEmail(email *model.Email) {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	// Check if we've exceeded retry limits
	if email.RetryCount >= s.options.MaxRetries {
		if email.CreatedAt.Add(s.options.RetryDuration).Before(time.Now()) {
			// Exceeded both retry count and duration
			email.Status = model.StatusFailed
			errMsg := "max retries exceeded and retry duration expired"
			email.Error = &errMsg
			s.logger.Error("email marked as failed - max retries exceeded",
				"email_id", email.ID,
				"recipient", email.Recipient,
				"retry_count", email.RetryCount)
			_ = s.emailStorage.Update(ctx, email)
			return
		}
	}

	// Attempt to send
	err := s.sender.Send(email.Recipient, email.Subject, email.Body)
	if err != nil {
		email.RetryCount++
		now := time.Now()
		email.LastRetry = &now
		errMsg := err.Error()
		email.LastError = &errMsg

		// Keep as pending if we can retry
		if email.RetryCount < s.options.MaxRetries {
			if email.CreatedAt.Add(s.options.RetryDuration).After(time.Now()) {
				email.Status = model.StatusPending
				s.logger.Warn("email send failed, will retry",
					"email_id", email.ID,
					"recipient", email.Recipient,
					"retry_count", email.RetryCount,
					"max_retries", s.options.MaxRetries,
					"error", err)
			} else {
				// Retry duration exceeded
				email.Status = model.StatusFailed
				email.Error = &errMsg
				s.logger.Error("email marked as failed - retry duration exceeded",
					"email_id", email.ID,
					"recipient", email.Recipient,
					"retry_count", email.RetryCount,
					"error", err)
			}
		} else {
			// Max retries exceeded
			email.Status = model.StatusFailed
			email.Error = &errMsg
			s.logger.Error("email marked as failed - max retries reached",
				"email_id", email.ID,
				"recipient", email.Recipient,
				"retry_count", email.RetryCount,
				"error", err)
		}
	} else {
		// Successfully sent
		email.Status = model.StatusSent
		now := time.Now()
		email.SentAt = &now
		s.logger.Info("email sent successfully",
			"email_id", email.ID,
			"recipient", email.Recipient,
			"subject", email.Subject,
			"retry_count", email.RetryCount)
	}

	// Update database
	if err := s.emailStorage.Update(ctx, email); err != nil {
		s.logger.Error("failed to update email status in database",
			"email_id", email.ID,
			"status", email.Status,
			"error", err)
	}
}

// processPendingEmails retrieves pending emails from the database and sends them
func (s *Service) processPendingEmails() {
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	emails, err := s.emailStorage.GetPending(ctx, 10)
	if err != nil {
		s.logger.Error("failed to retrieve pending emails from database",
			"error", err)
		return
	}

	if len(emails) > 0 {
		s.logger.Debug("processing pending emails",
			"count", len(emails))
		for i := range emails {
			s.sendEmail(&emails[i])
		}
	}
}
