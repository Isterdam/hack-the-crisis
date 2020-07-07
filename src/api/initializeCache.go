package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	jwt "github.com/dgrijalva/jwt-go"

	"sync"
	"time"
	"os"
)

// use this hacky solution to confirm bookings and companies - see public for further info
var ConfirmedBookings map[string]db.Booking
var ConfirmedCompanies map[string]db.Company
var ConfirmedHashes map[string]string

func InitializeCache() {
	ConfirmedBookings = make(map[string]db.Booking)
	ConfirmedCompanies = make(map[string]db.Company)
	ConfirmedHashes = make(map[string]string)

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

	go func(conf map[string]string) {
		var mutex = &sync.Mutex{}
		for {
			time.Sleep(15 * time.Minute)
			mutex.Lock()

			jwtKey := []byte(os.Getenv("JWTKEY"))

			loc, _ := time.LoadLocation("Europe/Stockholm")

			for k, v := range ConfirmedHashes {
				token, _ := jwt.ParseWithClaims(v, &PassClaims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtKey, nil
				})
				
				if payload, ok := token.Claims.(*PassClaims); ok {
					expiration := payload.StandardClaims.ExpiresAt

					if expiration <= time.Now().In(loc).Unix() {
						delete(ConfirmedHashes, k)
					}
				}
			}
			
			mutex.Unlock()
		}
	}(ConfirmedHashes)
}
