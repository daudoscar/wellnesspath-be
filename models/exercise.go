package models

import "time"

type Exercise struct {
	ID                     uint64    `gorm:"primaryKey;autoIncrement"`
	Name                   string    `gorm:"type:varchar(255);not null"`
	BodyPart               string    `gorm:"type:varchar(100);not null"`
	Difficulty             string    `gorm:"type:varchar(50);not null"`
	Category               string    `gorm:"type:varchar(100);not null"`
	ExerciseType           string    `gorm:"type:varchar(50);not null"`
	GoalTag                string    `gorm:"type:varchar(100);not null"`
	Description            string    `gorm:"type:text"`
	StepByStepInstructions string    `gorm:"type:text"`
	Equipment              string    `gorm:"type:varchar(255)"`
	IsDeleted              bool      `gorm:"default:false"`
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`
}
