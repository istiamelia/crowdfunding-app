package models

import (
	"time"

	"gorm.io/gorm"
)

type CampaignDB struct {
    ID              string `gorm:"primaryKey"`
    UserID          int32
    Title           string
    Description     string
    TargetAmount    int32
    CollectedAmount int32
    Deadline        time.Time
    Status          string `gorm:"type:campaign_status;default:'active'"`
    Category        string
    MinDonation     int32
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt
}

// Target schema and table
func (CampaignDB) TableName() string {
    return "campaigns.campaigns"
}
