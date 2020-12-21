package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

//CreateSchedule - スケジュール作成
func (d *Database) CreateSchedule(db *gorm.DB, logon *cx.Payload, sched *models.Schedule) (*models.Schedule, error) {
	db = d.ValidDB(db)

	data := models.Schedule{}
	err := copier.Copy(&data, sched)
	if err != nil {
		return nil, err
	}
	data.UpdateUser = logon.ID

	//Create
	result := db.Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil
}

//GetSchedule - スケジュール情報を取得
func (d *Database) GetSchedule(db *gorm.DB, logon *cx.Payload, id int64) (*models.Schedule, error) {
	db = d.ValidDB(db)

	//get
	// sql := "select t.*,d.detail as detail from menus t"
	// sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from menus_detail d group by id) d on (d.id=t.id)"
	// sql += " where t.id=? and t.tenant_id=?"
	sql := "select * from schedules where id=?"

	data := models.Schedule{}
	result := db.Raw(sql, id).Scan(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// data.Detail = "{" + data.Detail + "}"

	// err := utils.JSON.NewDecoder(strings.NewReader(data.Detail)).Decode(&data.Extra)
	// if err != nil {
	// 	return nil, err
	// }

	return &data, nil
}

//GetSchedules - スケジュール情報を取得（複数）
func (d *Database) GetSchedules(db *gorm.DB, logon *cx.Payload, id int64) ([]*models.Schedule, error) {
	db = d.ValidDB(db)

	//get
	// sql := "select t.*,d.detail as detail from menus t"
	// sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from menus_detail d group by id) d on (d.id=t.id)"
	// sql += " where t.tenant_id=?"
	sql := "select * from schedules where menu_id=?"
	result := make([]*models.Schedule, 0)
	rows, err := db.Raw(sql, id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.Schedule
		if err := db.ScanRows(rows, &record); err != nil {
			break
		}

		// record.Detail = "{" + record.Detail + "}"
		// err := utils.JSON.NewDecoder(strings.NewReader(record.Detail)).Decode(&record.Extra)
		// if err != nil {
		// 	return nil, err
		// }

		result = append(result, &record)
	}
	if len(result) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return result, nil
}

//UpdateSchedule - スケジュール情報を更新します
func (d *Database) UpdateSchedule(db *gorm.DB, logon *cx.Payload, id int64, data *models.Schedule) error {
	db = d.ValidDB(db)

	// //update
	// name := data["name"].(string)
	// if name != "" {
	// 	sql := "update menus set name=?,update_user=? where id=?"
	// 	err := db.Exec(sql, name, logon.ID, id).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if len(data) > 0 {
	// 	//get values
	// 	values := make([]string, 0)
	// 	for k, v := range data {
	// 		val := fmt.Sprintf("(%v,%q,\"%v\",%v)", id, k, v, logon.ID)

	// 		values = append(values, val)
	// 	}

	// 	//insert or update
	// 	sql := "insert into menus_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
	// 	err := db.Exec(sql).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

//DeleteSchedule - スケジュールを削除します
func (d *Database) DeleteSchedule(db *gorm.DB, logon *cx.Payload, id int64) error {
	db = d.ValidDB(db)

	//delete
	sql := "delete from schedules where id=?"
	err := db.Exec(sql, id).Error
	if err != nil {
		return err
	}

	return nil
}

//EnabledSchedule - スケジュールの利用可否（有効・無効）
func (d *Database) EnabledSchedule(db *gorm.DB, logon *cx.Payload, id int64, status int) error {
	db = d.ValidDB(db)

	//delete detail
	sql := "update schedules set status=? where id=?"
	err := db.Exec(sql, status, id).Error
	if err != nil {
		return err
	}

	return nil
}
