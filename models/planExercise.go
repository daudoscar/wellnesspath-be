package models

import "time"

type WorkoutPlanExercise struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	DayID      uint64    `gorm:"not null"`
	ExerciseID uint64    `gorm:"not null"`
	Order      int       `gorm:"not null"`
	Reps       int       `gorm:"not null"`
	Sets       int       `gorm:"not null"`
	Note       string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
