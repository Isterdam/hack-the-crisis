package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"

	"sync"
	"time"
)

// use this hacky solution to confirm bookings - see public for further info
var Confirmed map[string]db.Booking

// initialize very secret constants from local environment
func Initialize_constants() {
	Confirmed = make(map[string]db.Booking)

	go func(conf map[string]db.Booking) { // clear map every hour to avoid congestion
		var mutex = &sync.Mutex{}
		for {
			time.Sleep(60 * time.Minute)
			mutex.Lock()
			Confirmed = make(map[string]db.Booking)
			mutex.Unlock()
		}
	}(Confirmed)
}
