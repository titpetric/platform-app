package service

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"
	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"

	"github.com/titpetric/platform-app/pulse/storage"
	"github.com/titpetric/platform-app/pulse/view"
	"github.com/titpetric/platform-app/user"
	userstorage "github.com/titpetric/platform-app/user/storage"
)

// Handlers serves pulse HTTP endpoints.
type Handlers struct {
	storage     *storage.Storage
	userStorage *userstorage.UserStorage
	vuego       vuego.Template
	fs          fs.FS
}

// NewHandlers creates handlers backed by the given storage.
func NewHandlers(storage *storage.Storage, userStorage *userstorage.UserStorage) *Handlers {
	ofs := vuego.NewOverlayFS(view.FS, basecoat.FS)

	return &Handlers{
		fs:          ofs,
		storage:     storage,
		userStorage: userStorage,
		vuego:       vuego.NewFS(ofs),
	}
}

// Mount registers pulse routes on the router.
func (h *Handlers) Mount(r platform.Router) {
	r.Get("/assets/*", http.FileServer(http.FS(h.fs)).ServeHTTP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pulse", http.StatusFound)
	})
	r.Get("/pulse", h.IndexPage)
	r.Get("/pulse/{username}", h.UserPage)

	r.Group(func(r platform.Router) {
		r.Use(user.NewMiddleware(user.AuthHeader()))
		r.Post("/api/pulse/ingest", h.PostIngest)
	})
}

func (h *Handlers) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		ctx := r.Context()
		telemetry.CaptureError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// IndexPage serves the pulse index page.
func (h *Handlers) IndexPage(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.indexPage(w, r))
}

func (h *Handlers) indexPage(w http.ResponseWriter, r *http.Request) error {
	type userEntry struct {
		Username string `json:"username"`
		Count    int64  `json:"count"`
		Href     string `json:"href"`
	}

	type viewData struct {
		Menu  []any       `json:"menu"`
		Users []userEntry `json:"users"`
	}

	ctx := r.Context()

	users, err := h.userStorage.List(ctx)
	if err != nil {
		return fmt.Errorf("list users: %w", err)
	}

	counts, err := h.storage.ListUserCounts(ctx)
	if err != nil {
		return fmt.Errorf("list user counts: %w", err)
	}

	countByUser := make(map[string]int64, len(counts))
	for _, c := range counts {
		countByUser[c.UserID] = c.Count
	}

	entries := make([]userEntry, 0, len(users))
	for _, u := range users {
		entries = append(entries, userEntry{
			Username: u.Username,
			Count:    countByUser[u.ID],
			Href:     "/pulse/" + u.Username,
		})
	}

	data := viewData{
		Users: entries,
	}

	indexPage := vuego.View[viewData](h.vuego, "index.vuego", data)

	return indexPage.Render(ctx, w)
}

// UserPage serves the pulse user page.
func (h *Handlers) UserPage(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.userPage(w, r))
}

func (h *Handlers) userPage(w http.ResponseWriter, r *http.Request) error {
	type hourlyBar struct {
		Hour    int    `json:"hour"`
		Count   int64  `json:"count"`
		Percent int    `json:"percent"`
		Style   string `json:"style"`
		Tooltip string `json:"tooltip"`
		Label   string `json:"label"`
	}

	type hostSparkline struct {
		Hostname   string `json:"hostname"`
		Color      string `json:"color"`
		Points     string `json:"points"`
		Total      int64  `json:"total"`
		TotalLabel string `json:"totalLabel"`
		NumDays    int    `json:"numDays"`
		NumLabel   string `json:"numLabel"`
	}

	type viewData struct {
		Username   string          `json:"username"`
		FullName   string          `json:"fullName"`
		Hourly     []hourlyBar     `json:"hourly"`
		Sparklines []hostSparkline `json:"sparklines"`
		TotalCount int64           `json:"totalCount"`
	}

	ctx := r.Context()
	username := r.PathValue("username")

	user, err := h.userStorage.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	hourlyData, err := h.storage.GetUserHourly(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("get hourly data: %w", err)
	}

	dailyData, err := h.storage.GetUserDaily(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("get daily data: %w", err)
	}

	hosts, err := h.storage.GetUserHosts(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("get hosts: %w", err)
	}

	// Build hourly bars (24 hours, 0-23) with percentages for CSS
	hourlyMap := make(map[int]int64)
	var maxHourly int64
	for _, h := range hourlyData {
		hourlyMap[h.Hour] = h.Count
		if h.Count > maxHourly {
			maxHourly = h.Count
		}
	}
	hourly := make([]hourlyBar, 24)
	for i := 0; i < 24; i++ {
		count := hourlyMap[i]
		percent := 0
		if maxHourly > 0 && count > 0 {
			percent = int(count * 100 / maxHourly)
			if percent < 2 {
				percent = 2
			}
		}
		tooltip := fmt.Sprintf("%d:00 — %d keystrokes", i, count)
		var label string
		if i%3 == 0 {
			if i == 0 {
				label = "12a"
			} else if i < 12 {
				label = fmt.Sprintf("%da", i)
			} else if i == 12 {
				label = "12p"
			} else {
				label = fmt.Sprintf("%dp", i-12)
			}
		}
		hourly[i] = hourlyBar{
			Hour:    i,
			Count:   count,
			Percent: percent,
			Style:   fmt.Sprintf("height: %dpx", percent*180/100),
			Tooltip: tooltip,
			Label:   label,
		}
	}

	// Find date range from daily data (normalize stamps to just date)
	var minDate, maxDate string
	for _, d := range dailyData {
		stamp := d.Stamp
		if len(stamp) > 10 {
			stamp = stamp[:10]
		}
		if minDate == "" || stamp < minDate {
			minDate = stamp
		}
		if maxDate == "" || stamp > maxDate {
			maxDate = stamp
		}
	}

	// Build sparklines per host using data directly from dailyData
	colors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899", "#06b6d4", "#84cc16"}
	sparklines := make([]hostSparkline, 0, len(hosts))
	var totalCount int64

	// Build sorted date list and fill missing days
	var allDates []string
	if minDate != "" && maxDate != "" {
		startDate, _ := parseDate(minDate)
		endDate, _ := parseDate(maxDate)
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			allDates = append(allDates, d.Format("2006-01-02"))
		}
	}

	for i, host := range hosts {
		color := colors[i%len(colors)]
		var hostTotal int64
		var points []string
		numDays := 0
		label := "days"

		// Get counts for this host, keyed by the exact stamp from DB
		hostCounts := make(map[string]int64)
		for _, d := range dailyData {
			if d.Hostname == host {
				stamp := d.Stamp
				if len(stamp) > 10 {
					stamp = stamp[:10]
				}
				hostCounts[stamp] = d.Count
				hostTotal += d.Count
			}
		}

		// Build points for all dates
		for _, date := range allDates {
			points = append(points, fmt.Sprintf("%d", hostCounts[date]))
		}

		// Fallback to hourly data if fewer than 2 daily points
		if len(points) < 2 {
			hourlyByHost, err := h.storage.GetUserHourlyByHost(ctx, user.ID, host)
			if err != nil {
				return fmt.Errorf("get hourly by host: %w", err)
			}
			if len(hourlyByHost) >= 2 {
				points = points[:0]
				hostTotal = 0
				for _, h := range hourlyByHost {
					points = append(points, fmt.Sprintf("%d", h.Count))
					hostTotal += h.Count
				}
				label = "hours"
			}
		}

		numDays = len(points)
		totalCount += hostTotal

		sparklines = append(sparklines, hostSparkline{
			Hostname:   host,
			Color:      color,
			Points:     strings.Join(points, ","),
			Total:      hostTotal,
			TotalLabel: fmt.Sprintf("%d keystrokes", hostTotal),
			NumDays:    numDays,
			NumLabel:   label,
		})
	}

	data := viewData{
		Username:   user.Username,
		FullName:   user.FullName,
		Hourly:     hourly,
		Sparklines: sparklines,
		TotalCount: totalCount,
	}

	userPage := vuego.View[viewData](h.vuego, "user.vuego", data)

	return userPage.Render(ctx, w)
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// PostIngest handles pulse data ingestion requests.
func (h *Handlers) PostIngest(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.postIngest(w, r))
}

func (h *Handlers) postIngest(w http.ResponseWriter, r *http.Request) error {
	type ingestBody struct {
		Count    int64  `json:"count"`
		Hostname string `json:"hostname"`
	}

	body := ingestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return err
	}

	if body.Count <= 0 {
		return fmt.Errorf("count must be positive: %d", body.Count)
	}

	ctx := r.Context()
	if err := h.storage.Pulse(ctx, body.Count, body.Hostname); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
