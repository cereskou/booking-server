package models

//AccountConfirm -
type AccountConfirm struct {
	ID          int64  `gorm:"column:id;primary_key"` //id
	AccountID   int64  `gorm:"column:account_id"`     //アカウントID
	Email       string `gorm:"column:email"`          //Email
	ConfirmCode string `gorm:"column:confirm_code"`   //確認コード
}

// TableName sets the insert table name for this struct type
func (h *AccountConfirm) TableName() string {
	return "accounts_confirm"
}
