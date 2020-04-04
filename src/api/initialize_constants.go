package api

import (
	"github.com/Isterdam/hack-the-crisis-backend/src/db"

	"io/ioutil"
	"strings"
	"sync"
	"time"
)

var JWTkey string
var PhoneUser string
var PhonePass string

// use this hacky solution to confirm bookings - see public for further info
var Confirmed map[string]db.Booking

// initialize very secret constants from local environment
func Initialize_constants() {
	JWTkeyTemp, _ := ioutil.ReadFile("secretJWTKey.txt")
	JWTkey = string(JWTkeyTemp)

	phoneNumTemp, _ := ioutil.ReadFile("phoneNum.txt")
	temp := strings.Split(string(phoneNumTemp), "-")
	PhoneUser = temp[0]
	PhonePass = temp[1]

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
