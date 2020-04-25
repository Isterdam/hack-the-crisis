package db

import (
	"fmt"
	"time"

	null "gopkg.in/guregu/null.v3"
)

func GetCompanyByID(db *DB, id int) (Company, error) {
	stmt := db.prepared["company/getByID"]
	comp := Company{}
	err := stmt.Get(&comp, id)

	return comp, err
}

func GetCompanyByIDNoPass(db *DB, id int) (Company, error) {
	stmt := db.prepared["company/getByID"]
	comp := Company{}
	err := stmt.Get(&comp, id)

	comp.Password = null.NewString("", false)
	return comp, err
}

func GetCompaniesVerifiedPublic(db *DB) ([]CompanyPublic, error) {
	stmt := db.prepared["company/get/verified"]
	comps := []CompanyPublic{}
	err := stmt.Select(&comps)

	return comps, err
}

func GetCompanies(db *DB) ([]Company, error) {
	stmt := db.prepared["company/get"]
	comps := []Company{}
	err := stmt.Select(&comps)

	return comps, err
}

func GetCompaniesPublic(db *DB) ([]CompanyPublic, error) {
	stmt := db.prepared["company/get"]
	comps := []CompanyPublic{}
	err := stmt.Select(&comps)

	return comps, err
}

func InsertCompany(db *DB, comp Company) error {
	stmt := db.prepared["company/insert"]
	_, err := stmt.Exec(comp.Name, comp.Adress, comp.City, comp.Country, comp.PostCode, comp.CFirstName, comp.CNumber, comp.CLastName, comp.Email, comp.Password, comp.Longitude, comp.Latitude, comp.CEmail)
	fmt.Println(err)
	return err
}

func UpdateCompany(db *DB, comp Company) (Company, error) {
	var newComp Company

	stmt := db.prepared["company/update/location"]
	err := stmt.QueryRowx(comp.ID, comp.Name, comp.Adress, comp.City, comp.Country, comp.PostCode, comp.Longitude, comp.Latitude, comp.CEmail).StructScan(&newComp)

	if comp.CFirstName.String == "" {
		return newComp, err
	}

	stmt = db.prepared["company/update/contact"]
	err = stmt.QueryRowx(comp.ID, comp.CFirstName, comp.CLastName, comp.CNumber, comp.CEmail).StructScan(&newComp)

	if comp.Password.String == "" {
		return newComp, err
	}

	stmt = db.prepared["company/update/password"]
	err = stmt.QueryRowx(comp.ID, comp.Password).StructScan(&newComp)

	return newComp, err
}

func GetCompanyByEmail(db *DB, email string) (Company, error) {
	stmt := db.prepared["company/login"]

	retComp := Company{}
	err := stmt.Get(&retComp, email)

	return retComp, err
}

func GetCompaniesWithinDistance(db *DB, dist Distance) ([]CompanyPublic, error) {
	stmt := db.prepared["company/distance"]

	comps := []CompanyPublic{}
	err := stmt.Select(&comps, dist.LatMin, dist.LatMax, dist.LonMin, dist.LonMax, dist.Latitude, dist.Longitude, dist.R)

	return comps, err
}

func GetCompanyAverageAvailability(db *DB, compID int, startTime time.Time, days int) ([]float64, error) {
	stmt := db.prepared["company/get/avgAvailability"]

	av := []float64{}

	err := stmt.Select(&av, compID, startTime, days)

	return av, err
}

func GetCompanySlotAvailability(db *DB, compID int, startTime time.Time, days int) ([]int, error) {
	stmt := db.prepared["company/get/slotAvailability"]

	av := []int{}

	err := stmt.Select(&av, compID, startTime, days)

	return av, err
}
