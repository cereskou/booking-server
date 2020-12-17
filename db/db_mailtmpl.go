package db

import (
	"ditto/booking/models"

	"gorm.io/gorm"
)

//GetMailTemplate -
func (d *Database) GetMailTemplate(db *gorm.DB, tenantid int64, mailid string) (*models.MailTemplate, error) {
	db = d.ValidDB(db)

	var record models.MailTemplate
	result := db.Table("mails_template").
		Where("tenant_id=? and mail_id=? and enabled=1", tenantid, mailid).
		First(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return &record, nil
}
