package db

import (
	null "gopkg.in/guregu/null.v3"
)

type Slot struct {
	ID        null.Int  `db:"id" json:"id"`
	CompanyID null.Int  `db:"company_id" json:"company_id"`
	StartTime null.Time `db:"start_time" json:"start_time"`
	EndTime   null.Time `db:"end_time" json:"end_time"`
	MaxAmount null.Int  `db:"max" json:"max"`
	Booked    null.Int  `db:"booked" json:"booked"`
}

type Company struct {
	ID         null.Int    `db:"id" json:"id"`
	Name       null.String `db:"name" json:"name"`
	Adress     null.String `db:"adress" json:"adress"`
	City       null.String `db:"city" json:"city"`
	Country    null.String `db:"country" json:"country"`
	PostCode   null.String `db:"post_code" json:"post_code"`
	Longitude  null.Float  `db:"lon" json:"longitude"`
	Latitude   null.Float  `db:"lat" json:"latitude"`
	CFirstName null.String `db:"contact_firstname" json:"contact_firstname"`
	CLastName  null.String `db:"contact_lastname" json:"contact_lastname"`
	CNumber    null.String `db:"contact_number" json:"contact_number"`
	Verified   null.Bool   `db:"verified" json:"verified"`
	Email      null.String `db:"email" json:"email"`
	Password   null.String `db:"password" json:"password,omitempty"`
	CEmail     null.String `db:"contact_email" json:"contact_email"`
}

type CompanyPublic struct {
	ID         null.Int    `db:"id" json:"id"`
	Name       null.String `db:"name" json:"name"`
	Adress     null.String `db:"adress" json:"adress"`
	City       null.String `db:"city" json:"city"`
	Country    null.String `db:"country" json:"country"`
	PostCode   null.String `db:"post_code" json:"post_code"`
	Longitude  null.Float  `db:"lon" json:"longitude"`
	Latitude   null.Float  `db:"lat" json:"latitude"`
	CFirstName null.String `db:"contact_firstname" json:"-"`
	CLastName  null.String `db:"contact_lastname" json:"-"`
	CNumber    null.String `db:"contact_number" json:"-"`
	Verified   null.Bool   `db:"verified" json:"-"`
	Email      null.String `db:"email" json:"-"`
	Password   null.String `db:"password" json:"-"`
	CEmail     null.String `db:"contact_email" json:"-"`
	DistToUser null.Float  `db:"-" json:"dist_to_user"`
}

type Booking struct {
	ID          null.Int    `db:"id" json:"id"`
	SlotID      null.Int    `db:"slot_id" json:"slot_id"`
	PhoneNumber null.String `db:"phone_number" json:"phone_number"`
	Code        null.String `db:"code" json:"code"`
	FirstName   null.String `db:"first_name" json:"first_name"`
	LastName    null.String `db:"last_name" json:"last_name"`
	Visitee     null.String `db:"visitee" json:"visitee"`
	Message     null.String `db:"message" json:"message"`
	Status      null.String `db:"status" json:"status"`
}

type Distance struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Distance  int     `json:"distance"`
	LonMax    float64
	LonMin    float64
	LatMax    float64
	LatMin    float64
	R         float64
}

type Availabilty struct {
	CompanyID      int       `json:"id"`
	DailyAvailable []float64 `json:"availability_average"`
	AvailableSlots []int     `json:"available_slots`
}

type SearchQuery struct {
	String    null.String
	Longitude null.Float
	Latitude  null.Float
	Distance  uint64
	Limit     uint64
}
