package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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

func GetCompaniesAvailability(db *DB, compIDs []int, week int) ([]Availabilty, error) {
	query := `SELECT coalesce(company_id, crossp.id::int) as comp_id, coalesce(day, crossp.d) as dow, coalesce( sum(booked) / sum(max) ::float, 0) as avg
	FROM 
	(
		SELECT company_id, date_part('dow', start_time) as day, max, booked from slots s
		LEFT JOIN bookings b ON s.id=slot_id
		WHERE company_id IN (?) AND date_part('week', start_time) = (?)
		GROUP BY company_id, start_time, max, booked
	) t 
	RIGHT JOIN (
		SELECT a.d, c.id 
		FROM company c
		CROSS JOIN ( VALUES (1), (2), (3), (4), (5), (6), (0)) a (d)
		WHERE c.id IN (?)
	) crossp ON crossp.d=t.day AND crossp.id::int=t.company_id::int
	GROUP BY comp_id, dow
	ORDER BY comp_id`

	query, args, err := sqlx.In(query, compIDs, week, compIDs)

	query = db.DB.Rebind(query)

	stmt, err := db.DB.Preparex(query)

	res := []struct {
		CompanyID null.Int   `db:"comp_id"`
		DayOfWeek null.Int   `db:"dow"`
		Average   null.Float `db:"avg"`
	}{} //, len(compIDs)*7) //[]CompanyAvailabilityAverage{}

	err = stmt.Select(&res, args...)

	query = `SELECT coalesce(t.company_id, crossp.id::int) as comp_id, coalesce(t.dow, crossp.d) as dow, coalesce(t.count, 0) as count
	FROM 
	(
		SELECT company_id, date_part('dow', start_time) as dow, count(id) as count from slots s
		WHERE company_id IN (?) 
		AND booked < max
		AND date_part('week', start_time) = (?)
		GROUP BY company_id, dow
	) t
	RIGHT JOIN (
		SELECT a.d, c.id 
		FROM company c
		CROSS JOIN ( VALUES (1), (2), (3), (4), (5), (6), (0)) a (d)
		WHERE c.id IN (?)
	) crossp ON crossp.d=t.dow AND crossp.id=t.company_id
	ORDER BY comp_id, dow`

	query, args, err = sqlx.In(query, compIDs, week, compIDs)

	query = db.DB.Rebind(query)

	stmt, err = db.DB.Preparex(query)

	ress := []struct {
		CompanyID null.Int `db:"comp_id"`
		DayOfWeek null.Int `db:"dow"`
		Count     null.Int `db:"count"`
	}{}

	err = stmt.Select(&ress, args...)

	av := make([]Availabilty, len(compIDs))

	count := 0

	for i := range av {
		av[i].CompanyID = int(res[count].CompanyID.Int64)
		av[i].DailyAvailable = make([]float64, 7)
		av[i].AvailableSlots = make([]int, 7)
		for j := 0; j < 7; j++ {
			day := res[count].DayOfWeek.Int64 - 1
			if day == -1 {
				day = 6
			}
			av[i].DailyAvailable[day] = res[count].Average.Float64
			av[i].AvailableSlots[day] = int(ress[count].Count.Int64)
			count++
		}
	}

	return av, err
}
