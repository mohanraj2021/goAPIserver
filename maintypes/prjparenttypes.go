package maintypes

import (
	"math/rand"
	"time"
)

type Logtype string

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

const (
	Local  Logtype = "local"
	Remote Logtype = "Remote"
	// Both   Logtype = "Both"
)

var RandNumber = random(1, 1000)
