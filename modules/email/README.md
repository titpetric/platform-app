# Email Module

Email service for sending transactional emails via SMTP with persistent queue, retry logic, and comprehensive logging.

## TODO List

- [x] Create modules/email directory structure with model/, storage/, and smtp/ subpackages
- [x] Create Email model with ID, recipient, subject, body, status, created_at fields
- [x] Create email storage (Create, Get, Update methods for database operations)
- [x] Create Sender interface and SMTP implementation in smtp/ subpackage
- [x] Create Service struct with Start() method for in-memory queue consumer
- [x] Implement email queue consumer that sends emails via SMTP
- [x] Add mailhog service to docker-compose.yml
- [x] Configure PLATFORM_EMAIL_* environment variables for SMTP
- [x] Create integration tests for email sending to mailhog
- [x] Create modules/email/README.md with TODO list and usage examples
- [x] Implement configurable retry logic with max retries and duration
- [x] Add structured logging for email operations
- [x] Test Start() processes pending emails on startup
- [x] Test database updates with sent status
- [x] Make smtp.Config optional parameter with env fallback

## Architecture

### Components

- **model/email.go**: Email data model with status tracking and retry counters
- **storage/email.go**: Database operations for persisting emails
- **smtp/sender.go**: SMTP sender interface and implementation with ConfigFromEnv()
- **service.go**: Email service with configurable options, in-memory queue, and background worker
- **handler.go**: Module integration with platform framework

### Email Lifecycle

1. Email is added via `Service.AddEmail(ctx, email)`
2. Email is saved to database with `pending` status and `retry_count = 0`
3. Email is queued in in-memory channel or picked up by ticker
4. Background worker sends email via SMTP
5. On success: status updated to `sent`, `sent_at` timestamp recorded, retry_count kept
6. On failure: status remains `pending`, `retry_count` incremented, `last_error` and `last_retry` recorded
7. Pending emails are periodically polled (configurable interval) and retried
8. After max retries or duration exceeded: status updated to `failed`
9. All operations are logged with structured logging

## Configuration

### ServiceOptions

Configure email service behavior via `ServiceOptions`:

```go
type ServiceOptions struct {
	MaxRetries     int           // Maximum retry attempts (default: 3)
	RetryDuration  time.Duration // Max time to retry (default: 1 hour)
	TickerInterval time.Duration // Check pending interval (default: 5 seconds)
	QueueCapacity  int           // In-memory queue size (default: 100)
	Logger         *slog.Logger  // Structured logger
}
```

### Environment Variables

```bash
PLATFORM_EMAIL_HOST=mailhog          # SMTP host (default: localhost)
PLATFORM_EMAIL_PORT=1025             # SMTP port (default: 1025)
PLATFORM_EMAIL_USERNAME=              # SMTP username (optional)
PLATFORM_EMAIL_PASSWORD=              # SMTP password (optional)
PLATFORM_EMAIL_FROM=noreply@example.com  # From email address
```

### SMTP Configuration

The `smtp.Config` is optional when creating a Service. If `nil`, configuration is loaded from environment variables:

```go
// Option 1: Explicit config
smtpConfig := &smtp.Config{
    Host: "smtp.example.com",
    Port: 587,
    From: "noreply@example.com",
}
service := NewService(storage, smtpConfig)

// Option 2: From environment
service := NewService(storage, nil)  // Uses PLATFORM_EMAIL_* env vars
```

### Development with Mailhog

Start the Docker environment which includes Mailhog:

```bash
task up
```

Mailhog services:
- SMTP: `localhost:1025`
- Web UI: `http://localhost:8025`

## Usage Examples

### Basic Usage

```go
package mypackage

import (
    "context"
    "github.com/titpetric/platform-app/modules/email"
    "github.com/titpetric/platform-app/modules/email/model"
)

func SendWelcomeEmail(ctx context.Context, emailService *email.Service, userEmail string) error {
    email := model.NewEmail(
        userEmail,
        "Welcome to Platform App",
        "Thank you for signing up! Your account is ready to use.",
    )
    
    return emailService.AddEmail(ctx, email)
}
```

### With Custom Options

```go
import (
    "log/slog"
    "time"
)

// Create service with custom retry behavior
service := NewService(
    emailStorage,
    nil,  // Use environment variables for SMTP
    ServiceOptions{
        MaxRetries:     5,
        RetryDuration:  2 * time.Hour,
        TickerInterval: 30 * time.Second,
        QueueCapacity:  200,
        Logger:         slog.New(slog.NewTextHandler(os.Stdout, nil)),
    },
)

service.Start()
defer service.Stop()
```

### In a Handler

```go
package user

import (
    "net/http"
    "github.com/titpetric/platform-app/modules/email/model"
)

func (h *Service) Register(w http.ResponseWriter, r *http.Request) {
    // ... registration logic ...
    
    // Send welcome email (non-blocking)
    email := model.NewEmail(
        user.Email,
        "Welcome!",
        "Your account has been created successfully.",
    )
    
    if err := h.EmailService.AddEmail(r.Context(), email); err != nil {
        h.Error(r, "Failed to queue welcome email", err)
    }
    
    // ... rest of handler ...
}
```

## Email Model

```go
type Email struct {
    ID         string     // ULID identifier
    Recipient  string     // Email address
    Subject    string     // Email subject
    Body       string     // Email body
    Status     string     // pending, sent, or failed (in-memory, not persisted)
    CreatedAt  time.Time  // Creation timestamp
    SentAt     *time.Time // Sent timestamp (if sent)
    Error      *string    // Final error (if failed)
    RetryCount int        // Number of retry attempts
    LastError  *string    // Last send attempt error
    LastRetry  *time.Time // Last retry attempt timestamp
}
```

## Email Statuses

- `pending`: Email is waiting to be sent (in email_queue table)
- `sent`: Email was successfully delivered to SMTP (moved to email_sent table)
- `failed`: Email exceeded max retries or duration (moved to email_failed table)

## Database Schema

### email_queue Table
Stores emails pending to be sent. Emails are removed upon successful send or failure.

```sql
CREATE TABLE email_queue (
    id TEXT PRIMARY KEY,           -- ULID identifier
    recipient TEXT NOT NULL,       -- Email address
    subject TEXT NOT NULL,         -- Email subject
    body TEXT NOT NULL,            -- Email body
    created_at DATETIME NOT NULL,  -- Creation timestamp
    retry_count INTEGER DEFAULT 0, -- Number of retry attempts
    last_error TEXT,               -- Last send attempt error
    last_retry DATETIME            -- Last retry attempt timestamp
);
```

### email_sent Table
Audit trail for successfully sent emails.

```sql
CREATE TABLE email_sent (
    id TEXT PRIMARY KEY,           -- ULID identifier
    recipient TEXT NOT NULL,       -- Email address
    subject TEXT NOT NULL,         -- Email subject
    sent_at DATETIME NOT NULL,     -- Sent timestamp
    retry_count INTEGER DEFAULT 0  -- Number of retry attempts
);
```

### email_failed Table
Audit trail for emails that failed to send.

```sql
CREATE TABLE email_failed (
    id TEXT PRIMARY KEY,           -- ULID identifier
    recipient TEXT NOT NULL,       -- Email address
    subject TEXT NOT NULL,         -- Email subject
    failed_at DATETIME NOT NULL,   -- Failure timestamp
    retry_count INTEGER DEFAULT 0, -- Number of retry attempts
    error TEXT                     -- Final error message
);
```

## Retry Logic

The service implements automatic retry with two limits:

1. **Max Retries**: Maximum number of send attempts (default: 3)
2. **Retry Duration**: Maximum time window to attempt retries from creation (default: 1 hour)

An email is marked as failed when either:
- `retry_count >= MaxRetries` AND `retry_duration` has passed, OR
- `retry_duration` has passed (regardless of retry count)

### Example Scenario

With default options:
- Email created at 09:00
- First attempt fails at 09:00 → retry_count=1, status=pending
- Ticker picks up at 09:05 → second attempt fails → retry_count=2, status=pending
- Ticker picks up at 09:10 → third attempt fails → retry_count=3, status=pending
- Ticker picks up at 09:15 → fourth attempt fails → FAILED (max retries reached)
- Email marked as failed with error message recorded

## Testing

### Run Integration Tests

```bash
task integration  # Runs all integration tests
go test -tags integration -v ./modules/email/...
```

### Tests Include

1. **Mailhog Connectivity** - Verifies SMTP connection
2. **Email Sending** - Multiple emails sent successfully
3. **Service Queue** - Service processes and sends emails
4. **Start() Processing** - Pending emails sent on startup
5. **Retry Logic** - Failed emails retried correctly
6. **Database Updates** - Status and retry counters persisted
7. **Environment Config** - Config loaded from env vars
8. **Logging** - Structured logging works without errors

## Logging

The service uses Go's standard `log/slog` package for structured logging:

- **INFO**: Successful operations (emails sent, service started)
- **WARN**: Retriable errors (send failures before max retries)
- **ERROR**: Permanent failures or fatal conditions

Example log output:
```
level=info msg="email sent successfully" email_id=01ARZ3NDEKTSV4RRFFQ69G5FAV recipient=user@example.com subject="Welcome!"
level=warn msg="email send failed, will retry" email_id=01ARZ3NDEKTSV4RRFFQ69G5FAV retry_count=1 max_retries=3 error="connection refused"
level=error msg="email marked as failed - max retries exceeded" email_id=01ARZ3NDEKTSV4RRFFQ69G5FAV retry_count=3
```

## Performance Considerations

- **In-memory queue capacity**: 100 emails (configurable)
- **Ticker interval**: 5 seconds (configurable)
- **Concurrent sends**: Sequential (one at a time per worker)
- **Context timeout**: 30 seconds per send operation
- **Database**: Uses connection pooling from platform package

For high-volume scenarios, consider:
- Increasing `QueueCapacity` in ServiceOptions
- Decreasing `TickerInterval` for faster retries
- Running multiple service instances (one per process)
- Implementing external queue (Redis, RabbitMQ) for cross-service coordination

## Error Handling

Emails that fail to send are logged and queued for retry. The `last_error` field captures the most recent SMTP error for debugging. Failed emails after exhausting retries are marked as `failed` with the final error stored.

Check the logs for:
- `connection refused`: Network/firewall issues
- `invalid credentials`: Authentication problems
- `no route to host`: DNS/routing issues
- `timeout`: Slow SMTP server
