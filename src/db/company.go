package db

func GetCompanyByID(db *DB, id int) (Company, error) {
	stmt := db.prepared["company/getByID"]
	comp := Company{}
	err := stmt.Get(&comp, 1)

	return comp, err
}

func GetCompanies(db *DB) ([]Company, error) {
	stmt := db.prepared["company/get"]
	comps := []Company{}
	err := stmt.Select(&comps)

	return comps, err
}
