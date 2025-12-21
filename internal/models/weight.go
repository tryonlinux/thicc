package models

import (
	"time"

	"github.com/tryonlinux/thicc/internal/database"
)

// Weight represents a weight entry
type Weight struct {
	ID     int
	Date   string
	Weight float64
	BMI    float64
}

// AddWeight adds a new weight entry
func AddWeight(db *database.DB, date string, weight float64, bmi float64) error {
	_, err := db.Exec(
		"INSERT INTO weights (date, weight, bmi) VALUES (?, ?, ?)",
		date, weight, bmi,
	)
	return err
}

// DeleteWeight deletes a weight entry by ID
func DeleteWeight(db *database.DB, id int) error {
	_, err := db.Exec("DELETE FROM weights WHERE id = ?", id)
	return err
}

// ModifyWeight updates a weight entry
func ModifyWeight(db *database.DB, id int, weight float64, bmi float64) error {
	_, err := db.Exec("UPDATE weights SET weight = ?, bmi = ? WHERE id = ?", weight, bmi, id)
	return err
}

// GetWeights retrieves the last N weight entries
func GetWeights(db *database.DB, limit int) ([]Weight, error) {
	query := "SELECT id, date, weight, bmi FROM weights ORDER BY date DESC, id DESC LIMIT ?"
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weights []Weight
	for rows.Next() {
		var w Weight
		if err := rows.Scan(&w.ID, &w.Date, &w.Weight, &w.BMI); err != nil {
			return nil, err
		}
		weights = append(weights, w)
	}

	return weights, rows.Err()
}

// GetWeightsBetweenDates retrieves weight entries between two dates
func GetWeightsBetweenDates(db *database.DB, startDate, endDate string) ([]Weight, error) {
	query := "SELECT id, date, weight, bmi FROM weights WHERE date >= ? AND date <= ? ORDER BY date DESC, id DESC"
	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weights []Weight
	for rows.Next() {
		var w Weight
		if err := rows.Scan(&w.ID, &w.Date, &w.Weight, &w.BMI); err != nil {
			return nil, err
		}
		weights = append(weights, w)
	}

	return weights, rows.Err()
}

// GetTodayDate returns today's date in YYYY-MM-DD format
// Note: "2006-01-02" is Go's reference time format (Jan 2, 2006 at 3:04:05 PM MST)
func GetTodayDate() string {
	return time.Now().Format("2006-01-02")
}
