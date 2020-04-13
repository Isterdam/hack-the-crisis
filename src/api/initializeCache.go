package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"

	"sync"
	"time"
)

// use this hacky solution to confirm bookings and companies - see public for further info
var ConfirmedBookings map[string]db.Booking
var ConfirmedCompanies map[string]db.Company

func InitializeCache() {
	ConfirmedBookings = make(map[string]db.Booking)
	ConfirmedCompanies = make(map[string]db.Company)

	go func(conf map[string]db.Booking) { // clear map every hour to avoid congestion
		var mutex = &sync.Mutex{}
		for {
			time.Sleep(60 * time.Minute)
			mutex.Lock()
			ConfirmedBookings = make(map[string]db.Booking)
			mutex.Unlock()
		}
	}(ConfirmedBookings)

	go func(conf map[string]db.Company) {
		var mutex = &sync.Mutex{}
		for {
			time.Sleep(60 * time.Minute)
			mutex.Lock()
			ConfirmedCompanies = make(map[string]db.Company)
			mutex.Unlock()
		}
	}(ConfirmedCompanies)
}
