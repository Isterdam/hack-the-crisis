package db

import "gopkg.in/guregu/null.v3"

type Slot struct {
	ID        null.Int  `db:"id"`
	CompanyID null.Int  `db:"company_id"`
	StartTime null.Time `db:"start_time"`
	EndTime   null.Time `db:"end_time"`
	MaxAmount null.Int  `db:"max"`
	Day       null.Int  `db:"day"`
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
	ID          null.Int    `db:"id"`
	SlotID      null.Int    `db:"slot_id"`
	PhoneNumber null.String `db:"phone_number"`
	Code        null.String `db:"code"`
	FirstName   null.String `db:"first_name"`
	LastName    null.String `db:"last_name"`
}
