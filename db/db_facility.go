package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"
	"ditto/booking/utils"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//CreateFacility - 施設を作成
func (d *Database) CreateFacility(db *gorm.DB, logon *cx.Payload, tid int64, vmap map[string]interface{}) (*models.Facility, error) {
	db = d.ValidDB(db)

	data := models.Facility{}
	data.Name = vmap["name"].(string)
	data.TenantID = tid
	data.UpdateUser = logon.ID
	//Create
	result := db.Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	//詳細情報
	if len(vmap) > 1 {
		//get values
		values := make([]string, 0)
		for k, v := range vmap {
			val := fmt.Sprintf("(%v,%q,\"%v\",%v)", data.ID, k, v, logon.ID)

			values = append(values, val)
		}

		//insert update
		sql := "insert into facility_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}

//GetFacility - 施設情報を取得
func (d *Database) GetFacility(db *gorm.DB, logon *cx.Payload, id int64) (*models.FacilityWithDetail, error) {
	db = d.ValidDB(db)

	//get facility and details
	sql := "select t.*,d.detail as detail from facility t"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from facility_detail d group by id) d on (d.id=t.id)"
	sql += " where t.id=? and t.tenant_id=?"

	data := models.FacilityWithDetail{}
	result := db.Raw(sql, id, logon.Tenant).Scan(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	data.Detail = "{" + data.Detail + "}"

	err := utils.JSON.NewDecoder(strings.NewReader(data.Detail)).Decode(&data.Extra)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

//GetFacilities - 施設情報を取得（複数）
func (d *Database) GetFacilities(db *gorm.DB, logon *cx.Payload, tid int64) ([]*models.FacilityWithDetail, error) {
	db = d.ValidDB(db)

	//get facility and details
	sql := "select t.*,d.detail as detail from facility t"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from facility_detail d group by id) d on (d.id=t.id)"
	sql += " where t.tenant_id=?"

	result := make([]*models.FacilityWithDetail, 0)
	rows, err := db.Raw(sql, tid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.FacilityWithDetail
		if err := db.ScanRows(rows, &record); err != nil {
			break
		}

		record.Detail = "{" + record.Detail + "}"
		err := utils.JSON.NewDecoder(strings.NewReader(record.Detail)).Decode(&record.Extra)
		if err != nil {
			return nil, err
		}

		result = append(result, &record)
	}
	if len(result) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return result, nil
}

//UpdateFacility - 施設情報を更新します
func (d *Database) UpdateFacility(db *gorm.DB, logon *cx.Payload, id int64, data map[string]interface{}) error {
	db = d.ValidDB(db)

	//update facility
	name := data["name"].(string)
	if name != "" {
		sql := "update facility set name=?,update_user=? where id=?"
		err := db.Exec(sql, name, logon.ID, id).Error
		if err != nil {
			return err
		}
	}

	if len(data) > 0 {
		//get values
		values := make([]string, 0)
		for k, v := range data {
			val := fmt.Sprintf("(%v,%q,\"%v\",%v)", id, k, v, logon.ID)

			values = append(values, val)
		}

		//insert or update
		sql := "insert into facility_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return err
		}
	}

	return nil
}

//DeleteFacility - 施設を削除します
func (d *Database) DeleteFacility(db *gorm.DB, logon *cx.Payload, cid int64) error {
	db = d.ValidDB(db)

	//delete detail
	sql := "delete from facility_detail where id=?"
	err := db.Exec(sql, cid).Error
	if err != nil {
		return err
	}

	//delete master
	sql = "delete from facility where id=?"
	err = db.Exec(sql, cid).Error
	if err != nil {
		return err
	}

	return nil
}

//EnabledFacility - 施設の利用可否（有効・無効）
func (d *Database) EnabledFacility(db *gorm.DB, logon *cx.Payload, cid int64, status int) error {
	db = d.ValidDB(db)

	//delete detail
	sql := "update facility set status=? where id=?"
	err := db.Exec(sql, status, cid).Error
	if err != nil {
		return err
	}

	return nil
}
