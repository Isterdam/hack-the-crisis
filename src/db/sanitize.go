package db

import (
	"html"
)

// fields of type time and int should not be necessary to sanitize
// since they are parsed into structs from error. An invalid type would 
// already throw an error

type Sanitizer interface {
	Sanitize()
}

func (c *Company) Sanitize() {
	c.Name.String = html.EscapeString(c.Name.String)
	c.Adress.String = html.EscapeString(c.Adress.String)
	c.City.String = html.EscapeString(c.City.String)
	c.Country.String = html.EscapeString(c.Country.String)
	c.PostCode.String = html.EscapeString(c.PostCode.String)
	c.CFirstName.String = html.EscapeString(c.CFirstName.String)
	c.CLastName.String = html.EscapeString(c.CLastName.String)
	c.CNumber.String = html.EscapeString(c.CNumber.String)
	c.Email.String = html.EscapeString(c.Email.String)
	// does not sanitize passwords because maybe that would fuck em up?
	c.CEmail.String = html.EscapeString(c.CEmail.String)
}

func (b *Booking) Sanitize() {
	b.PhoneNumber.String = html.EscapeString(b.PhoneNumber.String)
	b.Code.String = html.EscapeString(b.Code.String)
	b.FirstName.String = html.EscapeString(b.FirstName.String)
	b.LastName.String = html.EscapeString(b.LastName.String)
	b.Visitee.String = html.EscapeString(b.Visitee.String)
	b.Message.String = html.EscapeString(b.Message.String)
	b.Status.String = html.EscapeString(b.Status.String)
}
