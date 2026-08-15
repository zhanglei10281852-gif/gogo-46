package retention

import (
	"testing"
	"time"

	"LogPilot/internal/config"
	"LogPilot/internal/model"
)

func TestCompactPropagatesFailedIncidentClosure(t *testing.T) {
	now := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	seen := now.AddDate(0, 0, -5)
	incident := model.Incident{
		ID: "incident-1", RuleID: "rule-1", Status: model.IncidentOpen,
		FirstSeen: seen, LastSeen: seen, OpenedAt: seen,
		UpdatedAt: now.Add(time.Hour), Count: 1, EventIDs: []string{"expired-event"},
	}
	compactor := New(config.Retention{EventDays: 1, IncidentDays: 90})
	compactor.Now = func() time.Time { return now }
	result, err := compactor.Compact(nil, []model.Incident{incident})
	if err == nil {
		t.Fatalf("Compact reported success after a rejected closure: %+v", result.Stats)
	}
}
