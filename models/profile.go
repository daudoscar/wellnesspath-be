package models

import "time"

type Profile struct {
	ID                 uint64 `gorm:"primaryKey;autoIncrement"`
	UserID             uint64 `gorm:"not null;unique"`
	SplitType          string `gorm:"type:varchar(50);not null"`
	Intensity          string `gorm:"type:varchar(50);not null"`
	TargetWeight       float64
	BMI                float64
	BMICategory        string `gorm:"type:varchar(50)"`
	Frequency          int
	DurationPerSession int
	Goal               string    `gorm:"type:varchar(100)"`
	EquipmentJSON      string    `gorm:"type:text"` // manually marshal/unmarshal []string
	IsDeleted          bool      `gorm:"default:false"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}
