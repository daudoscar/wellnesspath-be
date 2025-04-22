package models

import "time"

type WorkoutPlan struct {
	ID        uint64           `gorm:"primaryKey;autoIncrement"`
	UserID    uint64           `gorm:"not null"`
	SplitType string           `gorm:"type:varchar(50);not null"`
	Goal      string           `gorm:"type:varchar(100);not null"`
	IsDeleted bool             `gorm:"default:false"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"`
	Days      []WorkoutPlanDay `gorm:"foreignKey:PlanID"`
}
