package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"
	"ditto/booking/utils"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//CreateReservation - 予約作成
func (d *Database) CreateReservation(db *gorm.DB, logon *cx.Payload, sid int64, uid int64, vmap map[string]interface{}) (*models.Reservation, error) {
	db = d.ValidDB(db)

	data := models.Reservation{}
	data.ScheduleID = sid
	data.UserID = uid
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
		sql := "insert into reservation_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}

//GetReservation - 予約情報を取得
func (d *Database) GetReservation(db *gorm.DB, logon *cx.Payload, id int64) (*models.ReservationWithDetail, error) {
	db = d.ValidDB(db)

	//get
	sql := "select t.*,d.detail as detail from reservation t"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from reservation_detail d group by id) d on (d.id=t.id)"
	sql += " where t.id=?"

	data := models.ReservationWithDetail{}
	result := db.Raw(sql, id).Scan(&data)
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

//GetReservations - 予約情報を取得（複数）
func (d *Database) GetReservations(db *gorm.DB, logon *cx.Payload, schedid int64, userid int64) ([]*models.ReservationWithDetail, error) {
	db = d.ValidDB(db)

	//get
	sql := "select t.*,d.detail as reservation from menus t"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from reservation_detail d group by id) d on (d.id=t.id)"
	sql += " where (t.schedule_id=? or ?) and (t.user_id=? or ?)"

	es := (schedid == 0)
	eu := (userid == 0)

	result := make([]*models.ReservationWithDetail, 0)
	rows, err := db.Raw(sql, schedid, es, userid, eu).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.ReservationWithDetail
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

//UpdateReservation - 予約情報を更新します
func (d *Database) UpdateReservation(db *gorm.DB, logon *cx.Payload, id int64, data map[string]interface{}) error {
	db = d.ValidDB(db)

	//update
	if len(data) > 0 {
		//get values
		values := make([]string, 0)
		for k, v := range data {
			val := fmt.Sprintf("(%v,%q,\"%v\",%v)", id, k, v, logon.ID)

			values = append(values, val)
		}

		//insert or update
		sql := "insert into reservation_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return err
		}
	}

	return nil
}

//DeleteReservation - 予約を削除します
func (d *Database) DeleteReservation(db *gorm.DB, logon *cx.Payload, id int64) error {
	db = d.ValidDB(db)

	//delete classes_detail
	sql := "delete from reservation_detail where id=?"
	err := db.Exec(sql, id).Error
	if err != nil {
		return err
	}

	//delete classes
	sql = "delete from reservation where id=?"
	err = db.Exec(sql, id).Error
	if err != nil {
		return err
	}

	return nil
}

//EnabledReservation - 予約の利用可否（有効・無効）
func (d *Database) EnabledReservation(db *gorm.DB, logon *cx.Payload, id int64, status int) error {
	db = d.ValidDB(db)

	//delete detail
	sql := "update reservation set status=? where id=?"
	err := db.Exec(sql, status, id).Error
	if err != nil {
		return err
	}

	return nil
}
