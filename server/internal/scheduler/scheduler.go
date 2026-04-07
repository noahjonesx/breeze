package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/noahjonesx/breeze/server/internal/config"
	"github.com/noahjonesx/breeze/server/internal/notify"
	"github.com/noahjonesx/breeze/server/internal/store"
	"github.com/noahjonesx/breeze/server/internal/weather"
)

const (
	alertTwoHour   = "2hr"
	alertThirtyMin = "30min"

	// How wide a window around the target time we'll fire an alert.
	// 15-min poll interval → ±8 min catches every tick.
	alertWindow = 8 * time.Minute
)

type Scheduler struct {
	cfg   *config.Config
	store *store.Store
}

func New(cfg *config.Config, store *store.Store) *Scheduler {
	return &Scheduler{cfg: cfg, store: store}
}

func (s *Scheduler) Run() error {
	entries, err := weather.FetchForecast(s.cfg.Latitude, s.cfg.Longitude)
	if err != nil {
		return fmt.Errorf("fetch forecast: %w", err)
	}

	eventStart := findNextWindEvent(entries, s.cfg.ThresholdMPH)
	if eventStart == nil {
		log.Println("no upcoming wind events above threshold")
		return nil
	}

	log.Printf("next wind event at %s (%.0f mph threshold)", eventStart.Format(time.RFC3339), s.cfg.ThresholdMPH)

	now := time.Now()

	if inWindow(now, eventStart.Add(-2*time.Hour)) {
		err = s.maybeAlert(
			*eventStart,
			alertTwoHour,
			"Wind incoming in ~2 hours",
			fmt.Sprintf("Wind forecast to hit %.0f mph around %s. Close your gazebo curtains soon.",
				s.cfg.ThresholdMPH, eventStart.Format("3:04 PM")),
		)
		if err != nil {
			return err
		}
	}

	if inWindow(now, eventStart.Add(-30*time.Minute)) {
		err = s.maybeAlert(
			*eventStart,
			alertThirtyMin,
			"Close your curtains now",
			fmt.Sprintf("Wind forecast to hit %.0f mph in ~30 minutes.", s.cfg.ThresholdMPH),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func inWindow(now, target time.Time) bool {
	return now.After(target.Add(-alertWindow)) && now.Before(target.Add(alertWindow))
}

func findNextWindEvent(entries []weather.Entry, threshold float64) *time.Time {
	now := time.Now()
	for _, e := range entries {
		if e.Time.Before(now) {
			continue
		}
		if e.WindSpeed >= threshold {
			t := e.Time
			return &t
		}
	}
	return nil
}

func (s *Scheduler) maybeAlert(eventTime time.Time, alertType, title, body string) error {
	sent, err := s.store.AlreadySent(eventTime, alertType)
	if err != nil {
		return err
	}
	if sent {
		log.Printf("%s alert for %s already sent, skipping", alertType, eventTime.Format(time.RFC3339))
		return nil
	}

	log.Printf("sending %s alert for event at %s", alertType, eventTime.Format(time.RFC3339))

	if err := notify.Send(s.cfg.ExpoPushToken, title, body); err != nil {
		return fmt.Errorf("notify: %w", err)
	}

	return s.store.MarkSent(eventTime, alertType)
}
