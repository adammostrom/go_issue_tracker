package database

/* database/sql = Go standard lib for SQL abstraction
- defines interfaces like *sql.DB, *sql.Rows, Query, Exec, etc.
- Does not talk to postgres directly
*/

import (
	"database/sql"
	"servicetestess/models"
)

// shared connection pool, not a single connection
type DeviceStore struct {
	db *sql.DB
}

/*
Constructor function

This is the standard Go constructor pattern.
Go doesn’t have classes. This is how you “build” objects.

- Capital N → exported (visible outside the package) (in NewDeviceStore)
- Return type: pointer to DeviceStore
*/
func NewDeviceStore(db *sql.DB) *DeviceStore {
	// the db argument, which is a pointer to sql.DB, gets assigned to the DeviceStore struct db.
	// meaning that the sql.DB pointer is reachable via the DeviceStore db field.
	// & = address, returns pointer of type DeviceStore
	return &DeviceStore{db: db} // Field db gets the value of db
}

/*
(s *DeviceStore) = method reciever -> method belongs to DeviceStore.
"s" is like "this" or "self".
its a pointer so it operates on the real store

method GetAll (capital G = exported)
returns []models.Device -> slice of Device structs

  - []T means “slice of T”

  - models.Device means:

  - models → package name

  - Device → exported type inside that package

    So this returns:

    “A slice of Device structs defined in the models package, plus an error.”

    Basically a list of <Device>

    error -> Go explicit error handling

Go returns values, not exceptions.
*/
func (s *DeviceStore) GetAll() ([]models.Device, error) {
	// s.db -> access the db field from the struct DeviceStore which is a pointer to sql.DB
	rows, err := s.db.Query("SELECT device_id, serial_number, device_type FROM Devices ")

	// If anything goes wrong, err is not nil
	// Idiomatic Go, no error handling, errors checked immediately, function exits early.
	if err != nil {
		return nil, err
	}

	/*
		defer = schedule a function call to run when the current function returns.
		even if it returns early due to error later.

		Basically: “Run this function after the surrounding function finishes.”

		rows.Close() -> frees DB resources, it is required or connections will leak

		When GetAll() returns (for any reason), Go runs rows.Close().

		Its good because:
		- You declare cleanup next to acquisition
		- No try/finally
		- No forgotten cleanup in long functions
	*/
	defer rows.Close()

	/*
				- A slice named devices

				- Element type: models.Device

				- Initial value: nil

				Go slices are NOT arrays, they are descriptors over backing memory
				- Arrays = fixed size, size part of type
				- Slices = pointer to a backing array (dynamic array / List)
				-

				rows.Next() = advances cursor to next row, returns false if no more rows.

				var d models.Device = New Device struct
				- Allocates a Device struct
		     	- Zero-initializes all fields

				Scan = copies columnv values into Go variables, order must match SELECT order

				&d.Field = Pass pointer, scan writes into memory
				d.DeviceID → the field

				- &d.DeviceID → pointer to that field’s memory
				- Scan needs to write into variables.
				You’re saying:
					“Here is the memory location where you should put column 1.”

				devices = append(devices, d) = adds element to slice.
				slices does not mutate variables in place, it returns a new slice descriptor (or another one)
				I needs itself as argument to know where to return the new slice

				Doing this:
				append(devices, d)
				Will  discard the returned slice. Therefore correct is:
				slice (or other name) = append(slice, element)
	*/
	var devices []models.Device
	for rows.Next() {
		var d models.Device
		rows.Scan(&d.DeviceID, &d.SerialNumber, &d.DeviceType)
		devices = append(devices, d)
	}
	return devices, nil
}

// Add new device
func (s *DeviceStore) Add(serialNumber int64, deviceType string) error {
	_, err := s.db.Exec("INSERT INTO Devices(serial_number, device_type) VALUES($1, $2)", serialNumber, deviceType)
	return err
}

// Returns device based on serial number
/*
QueryRow:

does not return rows iterator

does not need Close()

executes immediately

Scan triggers the query

It’s a tiny, elegant shortcut for “exactly one row expected.”


*/
func (s *DeviceStore) GetBySerialNumber(serial_number int64) (models.Device, error) {

	var dev models.Device

	err := s.db.QueryRow(
		"SELECT device_id, serial_number, device_type FROM Devices WHERE serial_number = $1",
		serial_number).Scan(&dev.DeviceID, &dev.SerialNumber, &dev.DeviceType)

	if err != nil {
		return models.Device{}, err
	}

	return dev, nil
}
