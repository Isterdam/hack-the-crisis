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
	ID         null.Int    `db:"id" json:"id"`
	Name       null.String `db:"name" json:"name"`
	Adress     null.String `db:"adress" json:"adress"`
	City       null.String `db:"city" json:"city"`
	Country    null.String `db:"country" json:"country"`
	PostCode   null.String `db:"post_code" json:"post_code"`
	CFirstName null.String `db:"contact_firstname" json:"contact_firstname"`
	CLastName  null.String `db:"contact_lastname" json:"contact_lastname"`
	CNumber    null.String `db:"contact_number" json:"contact_number"`
	Verified   null.Bool   `db:"verified" json:"verified"`
	Email      null.String `db:"email" json:"email"`
	Password   null.String `db:"password" json:"password"`
}

type Booking struct {
	ID          null.Int    `db:"id"`
	SlotID      null.Int    `db:"slot_id"`
	PhoneNumber null.String `db:"phone_number"`
	Code        null.String `db:"code"`
	FirstName   null.String `db:"first_name"`
	LastName    null.String `db:"last_name"`
}
