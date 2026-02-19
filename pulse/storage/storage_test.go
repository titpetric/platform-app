package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/pulse/schema"
)

func newTestStorage(t *testing.T) *Storage {
	t.Helper()
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	err = Migrate(context.Background(), db, schema.Migrations())
	require.NoError(t, err)

	return NewStorage(db)
}

func seedDaily(t *testing.T, s *Storage, userID string, rows []struct {
	hostname string
	stamp    string
	count    int64
},
) {
	t.Helper()
	for _, r := range rows {
		_, err := s.db.Exec(
			`INSERT INTO pulse_daily (user_id, hostname, stamp, count) VALUES (?, ?, ?, ?)`,
			userID, r.hostname, r.stamp, r.count,
		)
		require.NoError(t, err)
	}
	for _, r := range rows {
		_, err := s.db.Exec(
			`INSERT OR IGNORE INTO pulse_hosts (user_id, hostname, created_at) VALUES (?, ?, CURRENT_TIMESTAMP)`,
			userID, r.hostname,
		)
		require.NoError(t, err)
	}
}

func TestGetUserDaily(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()
	userID := "TESTUSER"

	today := time.Now().UTC().Format("2006-01-02")
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")

	seedDaily(t, s, userID, []struct {
		hostname string
		stamp    string
		count    int64
	}{
		{"lab", today, 1000},
		{"lab", yesterday, 500},
		{"chronos", today, 200},
	})

	results, err := s.GetUserDaily(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, results, 3)

	for _, r := range results {
		t.Logf("hostname=%s stamp=%q count=%d", r.Hostname, r.Stamp, r.Count)
	}
}

func TestDailyBarChartLogic(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()
	userID := "TESTUSER"

	today := time.Now().UTC().Format("2006-01-02")
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")

	seedDaily(t, s, userID, []struct {
		hostname string
		stamp    string
		count    int64
	}{
		{"lab", today, 1000},
		{"lab", yesterday, 500},
		{"chronos", today, 200},
	})

	dailyData, err := s.GetUserDaily(ctx, userID)
	require.NoError(t, err)

	hosts, err := s.GetUserHosts(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, []string{"chronos", "lab"}, hosts)

	// Replicate handler logic: index by host
	hostDailyMap := make(map[string]map[string]int64)
	for _, d := range dailyData {
		if hostDailyMap[d.Hostname] == nil {
			hostDailyMap[d.Hostname] = make(map[string]int64)
		}
		hostDailyMap[d.Hostname][d.Stamp] = d.Count
	}

	// Build 30 date stamps
	now := time.Now().UTC()
	var dateStamps []string
	for i := 29; i >= 0; i-- {
		dateStamps = append(dateStamps, now.AddDate(0, 0, -i).Format("2006-01-02"))
	}

	// Verify lab totals and bar heights
	labCounts := hostDailyMap["lab"]
	require.NotNil(t, labCounts, "lab should have daily data")
	assert.Equal(t, int64(1000), labCounts[today])
	assert.Equal(t, int64(500), labCounts[yesterday])

	var totalCount int64
	for _, host := range hosts {
		counts := hostDailyMap[host]
		for _, stamp := range dateStamps {
			totalCount += counts[stamp]
		}
	}
	assert.Equal(t, int64(1700), totalCount)

	// Verify percentage calculation
	var maxCount int64
	for _, stamp := range dateStamps {
		if c := labCounts[stamp]; c > maxCount {
			maxCount = c
		}
	}
	assert.Equal(t, int64(1000), maxCount)
	assert.Equal(t, "height: 100%", fmt.Sprintf("height: %d%%", labCounts[today]*100/maxCount))
	assert.Equal(t, "height: 50%", fmt.Sprintf("height: %d%%", labCounts[yesterday]*100/maxCount))
}
