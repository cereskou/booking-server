package db

import (
	"ditto/booking/models"

	"gorm.io/gorm"
)

//GetCasbinPolicies -
func (d *Database) GetCasbinPolicies(db *gorm.DB) ([]*models.CasbinRule, error) {
	db = d.ValidDB(db)

	result := make([]*models.CasbinRule, 0)
	rows, err := db.Table("casbin_rules").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.CasbinRule
		if err := db.ScanRows(rows, &record); err != nil {
			break
		}

		result = append(result, &record)
	}
	if len(result) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return result, nil
}
