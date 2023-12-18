package db

import (
	"encoding/json"
	"os"

	"github.com/Lutz-Pfannenschmidt/ResponsePlan/internal/db/models"
)

type Database struct {
	Data map[string]models.Scan
}

func NewDatabase() *Database {
	return &Database{
		Data: make(map[string]models.Scan),
	}
}

func (db *Database) SaveToFile(filepath string) error {
	data, err := json.Marshal(db.Data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) LoadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &db.Data)
	if err != nil {
		return err
	}

	return nil
}
