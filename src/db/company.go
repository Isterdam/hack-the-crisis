package db

func GetCompanyByID(db *DB, id int) (Company, error) {
	stmt := db.prepared["company/getByID"]
	comp := Company{}
	err := stmt.Get(&comp, 1)

	return comp, err
}

func GetCompanies(db *DB) (Company, error) {
	stmt := db.prepared["company/get"]
	comp := Company{}
	err := stmt.Get(&comp)

	return comp, err
}
