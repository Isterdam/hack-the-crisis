package db

import "fmt"

func GetCompanyByID(db *DB, id int) (Company, error) {
	stmt := db.prepared["company/getByID"]
	comp := Company{}
	err := stmt.Get(&comp, id)

	return comp, err
}

func GetCompanies(db *DB) ([]Company, error) {
	stmt := db.prepared["company/get"]
	comps := []Company{}
	err := stmt.Select(&comps)

	return comps, err
}

func InsertCompany(db *DB, comp Company) error {
	stmt := db.prepared["company/insert"]
	fmt.Printf("%#v", comp.CNumber)
	_, err := stmt.Exec(comp.Name, comp.Adress, comp.City, comp.Country, comp.PostCode, comp.CFirstName, comp.CNumber, comp.CLastName, comp.Email, comp.Password)
	return err
}

func UpdateCompany(db *DB, comp Company) (Company, error) {
	var newComp Company

	stmt := db.prepared["company/update/location"]
	err := stmt.QueryRowx(comp.ID, comp.Name, comp.Adress, comp.City, comp.Country, comp.PostCode).StructScan(&newComp)

	if comp.CFirstName.String == "" {
		return newComp, err
	}

	stmt = db.prepared["company/update/contact"]
	err = stmt.QueryRowx(comp.ID, comp.CFirstName, comp.CLastName, comp.CNumber).StructScan(&newComp)

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
