package model

type ServiceMap struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"uniqueIndex"`
	Value string `json:"value" gorm:"uniqueIndex"`
	Auth  bool   `json:"bool" gorm:"default:false"`
}

