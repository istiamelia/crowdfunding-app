package helper

import (
	"campaign-service/models"
	"errors"
	"time"
)

func ValidateCampaign(c models.CampaignDB) error {
    if c.UserID == 0 {
        return errors.New("user ID is required")
    }
    if len(c.Title) < 5 || len(c.Title) > 100 {
        return errors.New("title must be between 5 and 100 characters")
    }
    if c.TargetAmount <= 0 {
        return errors.New("target amount must be positive")
    }
	if c.MinDonation <= 0 {
        return errors.New("minimal donation amount must be positive")
    }
    if c.Deadline.Before(time.Now()) {
        return errors.New("deadline must be in the future")
    }
    if c.MinDonation <= 0 || c.MinDonation > c.TargetAmount {
        return errors.New("min donation must be positive and less than target amount")
    }
    // Add more as needed, e.g. validate category is within accepted values
    return nil
}

func ValidateUpdateCampaign(c models.CampaignDB) error {
	if c.Title != "" && (len(c.Title) < 5 || len(c.Title) > 100) {
		return errors.New("title must be between 5 and 100 characters")
	}
	if c.TargetAmount != 0 && c.TargetAmount <=0 {
		return errors.New("target amount must be positive")
	}
	if !c.Deadline.IsZero() && c.Deadline.Before(time.Now()) {
		return errors.New("deadline must be a future date")
	}
	if c.MinDonation != 0{
		if c.MinDonation <= 0{
			return errors.New("minimum donation amount must be positive")
		}else if c.TargetAmount != 0 {
				if c.MinDonation > c.TargetAmount {
				return errors.New("minimum donation must be greater than 0 and less than or equal to target amount")
			}
		}
	}
    // Add more as needed, e.g. validate category is within accepted values
    return nil
}