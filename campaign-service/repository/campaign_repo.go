package repository

import (
	"campaign-service/models"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// CampaignRepository defines methods for interacting with campaign-related data in the database.
type CampaignRepository interface {
	CreateCampaign(campaign models.CampaignDB) (interface{}, error) 
	GetCampaignByID(campaignID string) (interface{}, error) 
	DeleteCampaignByID(id string) (int32, error) 
	UpdateCampaignByID(id string, userID int32,campaign models.CampaignDB) (interface{},error) 
	GetCampaignsByUserID(userID int32) (interface{}, error) 
}

// campaignRepository is the concrete implementation of CampaignRepository.
// It uses a gorm database connection to perform operations.
type campaignRepository struct {
	db *gorm.DB
}

// Constructor NewCampaignRepository creates and returns a new instance of campaignRepository,
// injecting the gorm database connection.
func NewCampaignRepository(db *gorm.DB) CampaignRepository {
	return &campaignRepository{db: db}
}

func (r *campaignRepository) CreateCampaign(campaign models.CampaignDB) (interface{}, error) {
	// Insert campaign to campaigns schema and campaigns table
	result := r.db.Create(&campaign)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.Internal, "Failed to create a campaign")
	}

	return campaign, nil
}

func (r *campaignRepository) GetCampaignByID(id string) (interface{}, error) {
	var campaign models.CampaignDB
	// Get campaign by id where deleted_at != nil
	if err := r.db.First(&campaign, "id=?",id).Error; err != nil{
		return nil, status.Error(codes.NotFound,"Campaign not found")
	}
	return campaign,nil
}

func (r *campaignRepository) DeleteCampaignByID(id string) (int32, error) {
	// Check if campaign exist in table
	retreivedCampaign, err := r.GetCampaignByID(id)
	if err != nil {
		return 0, err
	}

	// Cast to campaign.db
	castedCampaign, ok := retreivedCampaign.(models.CampaignDB)
	if !ok{
		return 0, status.Error(codes.Internal,"Failed to cast campaign")
	}

	// Delete data and update status to "cancelled"
	var campaign models.CampaignDB
	if err := r.db.Model(&campaign).Where("id=?", id).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"status":     "cancelled", // Ensure this matches your enum string
	}).Error; err != nil{
		return 0, status.Error(codes.Internal,"Error deleting campaign")
	}

	return castedCampaign.UserID, nil
}

func (r *campaignRepository) UpdateCampaignByID(id string, userID int32,campaign models.CampaignDB) (interface{},error) {
	// Check if campaign exist in table
	retreivedCampaign, err := r.GetCampaignByID(id)
	if err != nil {
		return nil,err
	}
	// Cast to models.campaignDB
	castedCampaign, ok := retreivedCampaign.(models.CampaignDB)
	if !ok{
		return nil, status.Error(codes.Internal,"Failed to cast campaign")
	}

	// check if current status cancelled or completed
	if castedCampaign.Status == "cancelled" || castedCampaign.Status == "completed"{
		return nil, status.Errorf(codes.FailedPrecondition,"campaign status is %v",castedCampaign.Status)
	}
	// Update data
	result := r.db.Model(campaign).Where("id=? AND user_id=?", id, userID).Updates(campaign) 
	if result.Error != nil{
		return nil, status.Error(codes.Internal,"Error updating campaign")
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.FailedPrecondition, "No changes were made")
	}

	updatedCampaign, err := r.GetCampaignByID(id)
	if err != nil {
		return nil,err
	}
	return updatedCampaign, nil
}

func (r *campaignRepository) GetCampaignsByUserID(userID int32) (interface{}, error) {
	var campaign []models.CampaignDB
	// Get campaign by id where deleted_at != nil
	if err := r.db.Where("user_id=?",userID).Find(&campaign).Error; err != nil{
		return nil, status.Error(codes.NotFound,"Campaign not found")
	}
	return campaign,nil
}




