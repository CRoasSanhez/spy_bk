package core

import (
	"math/rand"
	"time"
)

// CustomRandomInt generates a custom random with the time as seed
func CustomRandomInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
