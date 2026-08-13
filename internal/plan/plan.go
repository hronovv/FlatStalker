package plan

import (
	"fmt"
	"time"
)

const (
	Free = "free"
	Plus = "plus"
	Pro  = "pro"
)

// Intervals holds poll intervals for each user plan.
type Intervals struct {
	Free time.Duration
	Plus time.Duration
	Pro  time.Duration
}

func (i Intervals) For(name string) time.Duration {
	switch name {
	case Plus:
		return i.Plus
	case Pro:
		return i.Pro
	default:
		return i.Free
	}
}

func Normalize(name string) string {
	switch name {
	case Plus, Pro:
		return name
	default:
		return Free
	}
}

func Label(name string) string {
	switch Normalize(name) {
	case Plus:
		return "PLUS"
	case Pro:
		return "PRO"
	default:
		return "FREE"
	}
}

func Rank(name string) int {
	switch Normalize(name) {
	case Pro:
		return 2
	case Plus:
		return 1
	default:
		return 0
	}
}

func Effective(name string, expiresAt *time.Time, now time.Time) string {
	p := Normalize(name)
	if p == Free {
		return Free
	}
	if expiresAt == nil || expiresAt.After(now) {
		return p
	}
	return Free
}

func FormatDaysRU(n int) string {
	if n < 1 {
		n = 1
	}
	return fmt.Sprintf("%d %s", n, russianDays(n))
}

func russianDays(n int) string {
	n = n % 100
	if n >= 11 && n <= 14 {
		return "дней"
	}
	switch n % 10 {
	case 1:
		return "день"
	case 2, 3, 4:
		return "дня"
	default:
		return "дней"
	}
}

func LinkLimit(name string) int {
	switch Normalize(name) {
	case Plus:
		return 3
	case Pro:
		return 5
	default:
		return 1
	}
}

// FormatIntervalRU returns a short Russian phrase like "каждые 5 минут".
func FormatIntervalRU(d time.Duration) string {
	if d < time.Minute {
		secs := int(d.Round(time.Second) / time.Second)
		if secs <= 0 {
			secs = 1
		}
		return fmt.Sprintf("каждые %d %s", secs, russianSeconds(secs))
	}
	mins := int(d.Round(time.Minute) / time.Minute)
	if mins <= 0 {
		mins = 1
	}
	return fmt.Sprintf("каждые %d %s", mins, russianMinutes(mins))
}

func russianSeconds(n int) string {
	n = n % 100
	if n >= 11 && n <= 14 {
		return "секунд"
	}
	switch n % 10 {
	case 1:
		return "секунду"
	case 2, 3, 4:
		return "секунды"
	default:
		return "секунд"
	}
}

func russianMinutes(n int) string {
	n = n % 100
	if n >= 11 && n <= 14 {
		return "минут"
	}
	switch n % 10 {
	case 1:
		return "минуту"
	case 2, 3, 4:
		return "минуты"
	default:
		return "минут"
	}
}
