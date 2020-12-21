package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"
	"ditto/booking/utils"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//CreateMenu - メニュー作成
func (d *Database) CreateMenu(db *gorm.DB, logon *cx.Payload, tid int64, vmap map[string]interface{}) (*models.Menu, error) {
	db = d.ValidDB(db)

	data := models.Menu{}
	data.Name = vmap["name"].(string)
	data.TenantID = tid
	data.OwnerID = logon.ID
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
		sql := "insert into menus_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}

//GetMenu - メニュー情報を取得
func (d *Database) GetMenu(db *gorm.DB, logon *cx.Payload, id int64) (*models.MenuWithDetail, error) {
	db = d.ValidDB(db)

	//get facility and details
	sql := "select t.*,d.detail as detail from menus t"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from menus_detail d group by id) d on (d.id=t.id)"
	sql += " where t.id=? and t.tenant_id=?"

	data := models.MenuWithDetail{}
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

//GetMenus - 施設情報を取得（複数）
func (d *Database) GetMenus(db *gorm.DB, logon *cx.Payload, tid int64) ([]*models.MenuWithDetail, error) {
	db = d.ValidDB(db)

	//get facility and details
	sql := "select t.*,d.detail as detail from menus t"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from menus_detail d group by id) d on (d.id=t.id)"
	sql += " where t.tenant_id=?"

	result := make([]*models.MenuWithDetail, 0)
	rows, err := db.Raw(sql, tid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.MenuWithDetail
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

//UpdateMenu - メニュー情報を更新します
func (d *Database) UpdateMenu(db *gorm.DB, logon *cx.Payload, id int64, data map[string]interface{}) error {
	db = d.ValidDB(db)

	//update
	name := data["name"].(string)
	if name != "" {
		sql := "update menus set name=?,update_user=? where id=?"
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
		sql := "insert into menus_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return err
		}
	}

	return nil
}

//DeleteMenu - メニューを削除します
func (d *Database) DeleteMenu(db *gorm.DB, logon *cx.Payload, id int64) error {
	db = d.ValidDB(db)

	//delete detail
	sql := "delete from menus_detail where id=?"
	err := db.Exec(sql, id).Error
	if err != nil {
		return err
	}

	//delete master
	sql = "delete from menus where id=?"
	err = db.Exec(sql, id).Error
	if err != nil {
		return err
	}

	return nil
}

//EnabledMenu - メニューの利用可否（有効・無効）
func (d *Database) EnabledMenu(db *gorm.DB, logon *cx.Payload, id int64, status int) error {
	db = d.ValidDB(db)

	//delete detail
	sql := "update menus set status=? where id=?"
	err := db.Exec(sql, status, id).Error
	if err != nil {
		return err
	}

	return nil
}
