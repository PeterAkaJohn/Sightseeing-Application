package model

import "log"

//Location : when the user favorites a location this gets inserted in the database
type Location struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Position    string `json:"position"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Description string `json:"description"`
	OpenHours   string `json:"open_hours"`
	CloseHours  string `json:"close_hours"`
}

//GetLocation retrieves the location from the database with the use of the ID
func GetLocation(locationID int64) (*Location, error) {
	result := Location{}

	//Use database connection
	row := db.QueryRow(
		"SELECT id, name, position, address, city, state, postal_code, description, open_hours, close_hours "+
			"FROM location "+
			"WHERE id = $1", locationID)

	err := row.Scan(&result.ID, &result.Name, &result.Position, &result.Address, &result.City, &result.State, &result.PostalCode,
		&result.Description, &result.OpenHours, &result.CloseHours)

	return &result, err
}

//GetLocations : retrieves all the locations in the database
func GetLocations() ([]*Location, error) {
	result := []*Location{}

	rows, err := db.Query("SELECT *" +
		"FROM location")
	if err != nil {
		log.Print(err)
	} else {
		for rows.Next() {
			location := Location{}
			rows.Scan(&location.ID, &location.Name, &location.Position, &location.Address, &location.City,
				&location.PostalCode, &location.Description, &location.OpenHours, &location.CloseHours)
			result = append(result, &location)
		}
	}

	return result, err
}

//AddLocation : called when favoriting a location
func AddLocation(name string, position string, address string, city string, state string, postal_code string, description string, open_hours string, close_hours string) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO locations(name, position, address, city, state, postal_code, description, open_hours, close_hours) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	res, err := stmt.Exec(name, position, address, city, state, postal_code, description, open_hours, close_hours)
	if err != nil {
		log.Print(err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)

	return lastID, err
}
