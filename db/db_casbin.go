package db

import "ditto/booking/models"

//GetCasbinPolicies -
func (d *Database) GetCasbinPolicies() ([]*models.CasbinRule, error) {
	result := make([]*models.CasbinRule, 0)
	rows, err := d.DB().Table("casbin_rules").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.CasbinRule
		if err := d.DB().ScanRows(rows, &record); err != nil {
			break
		}

		result = append(result, &record)
	}

	return result, nil
}
