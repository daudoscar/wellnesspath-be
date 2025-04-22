package models

import "time"

type WorkoutPlanDay struct {
	ID        uint64                `gorm:"primaryKey;autoIncrement"`
	PlanID    uint64                `gorm:"not null"`
	DayNumber int                   `gorm:"not null"`
	Focus     string                `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time             `gorm:"autoCreateTime"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime"`
	Exercises []WorkoutPlanExercise `gorm:"foreignKey:DayID"`
}
