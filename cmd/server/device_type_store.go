package database

import (
	"database/sql"
	"servicetestess/models"
)

type DeviceTypeStore struct {
	db *sql.DB
}

// Constructor function

func NewDeviceTypeStore(db *sql.DB) *DeviceTypeStore {
	return &DeviceTypeStore{db: db}
}

// func <bind to type struct > function name <no arguments> return slice of DeviceType (List<DeviceType>)
func (self *DeviceTypeStore) GetAll() ([]models.DeviceType, error) {

	// Sql query
	rows, err := self.db.Query("SELECT device_type FROM DeviceTypes ")

	if err != nil {
		// models.DeviceType returns as nil, due to failed operation, err returned as error
		return nil, err
	}
	defer rows.Close()

	var device_types []models.DeviceType
	for rows.Next() {
		var dt models.DeviceType
		rows.Scan(&dt.DeviceType)
		device_types = append(device_types, dt)
	}
	// Error returns as nil due to successful operation
	return device_types, nil

}
