package db

import "gopkg.in/guregu/null.v3"

type Slot struct {
	ID        null.Int  `db:"id" json:"id"`
	CompanyID null.Int  `db:"company_id" json:"company_id"`
	StartTime null.Time `db:"start_time" json:"start_time"`
	EndTime   null.Time `db:"end_time" json:"end_time"`
	MaxAmount null.Int  `db:"max" json:"max"`
	Day       null.Int  `db:"day" json:"day"`
}

type Company struct {
	ID         null.Int    `db:"id"`
	Name       null.String `db:"name"`
	Adress     null.String `db:"adress"`
	City       null.String `db:"city"`
	Country    null.String `db:"country"`
	PostCode   null.String `db:"post_code"`
	CFirstName null.String `db:"contact_firstname"`
	CLastName  null.String `db:"contact_lastname"`
	CNumber    null.String `db:"contact_number"`
	Verified   null.Bool   `db:"verified"`
	Email      null.String `db:"email"`
	Password   null.String `db:"password"`
}

type Booking struct {
	ID          null.Int    `db:"id" json:"id"`
	SlotID      null.Int    `db:"slot_id" json:"slot_id"`
	PhoneNumber null.String `db:"phone_number" json:"phone_number"`
	Code        null.String `db:"code" json:"code"`
	FirstName   null.String `db:"first_name" json:"first_name"`
	LastName    null.String `db:"last_name" json:"last_name"`
}
