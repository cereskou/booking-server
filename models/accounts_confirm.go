package models

//AccountConfirm -
type AccountConfirm struct {
	ID          int64  `json:"id" gorm:"column:id;primary_key"`         //id
	AccountID   int64  `json:"account_id" gorm:"column:account_id"`     //アカウントID
	Email       string `json:"email" gorm:"column:email"`               //Email
	ConfirmCode string `json:"confirm_code" gorm:"column:confirm_code"` //確認コード
}

// TableName sets the insert table name for this struct type
func (h *AccountConfirm) TableName() string {
	return "accounts_confirm"
}
