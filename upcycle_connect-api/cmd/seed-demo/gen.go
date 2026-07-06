package main

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

var rng = rand.New(rand.NewSource(42))

func newID() string {
	return uuid.New().String()
}

func pick[T any](list []T) T {
	return list[rng.Intn(len(list))]
}

func pickExcept(list []string, except string) string {
	for {
		v := pick(list)
		if v != except || len(list) == 1 {
			return v
		}
	}
}

func percent(p int) bool {
	return rng.Intn(100) < p
}

func randMoney(min, max float64) float64 {
	v := min + rng.Float64()*(max-min)
	return float64(int(v*100)) / 100
}

func randInt(min, max int) int {
	if max <= min {
		return min
	}
	return min + rng.Intn(max-min+1)
}

func daysAgo(days int) time.Time {
	return time.Now().AddDate(0, 0, -days)
}

func daysFromNow(days int) time.Time {
	return time.Now().AddDate(0, 0, days)
}

func sqlDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func sqlDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

type pairSet map[string]struct{}

func newPairSet() pairSet {
	return make(pairSet)
}

func fillPairs(target int, aList, bList []string, exclude pairSet) [][2]string {
	if exclude == nil {
		exclude = newPairSet()
	}
	var result [][2]string
	attempts := 0
	maxAttempts := target * 30
	for len(result) < target && attempts < maxAttempts {
		attempts++
		a := pick(aList)
		b := pick(bList)
		if !exclude.add(a, b) {
			continue
		}
		result = append(result, [2]string{a, b})
	}
	return result
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func (s pairSet) add(a, b string) bool {
	key := a + "|" + b
	if _, exists := s[key]; exists {
		return false
	}
	s[key] = struct{}{}
	return true
}
