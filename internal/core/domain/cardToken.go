package domain

type CardToken struct {
	Token string `gorm:"primaryKey" json:"token"`
}
