package service

import (
	"campaign-service/gen/go/campaign/v1"
	"campaign-service/helper"
	"campaign-service/models"
	"campaign-service/mq"
	"campaign-service/repository"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CampaignService interface defines the Service methods for campaign operations
type CampaignService interface {
	CreateCampaign(ctx context.Context, req *campaign.CreateCampaignRequest) (*campaign.CreateCampaignResponse, error)
	GetCampaignByID(ctx context.Context, req *campaign.GetCampaignByIDRequest) (*campaign.GetCampaignByIDResponse, error)
	DeleteCampaignByID(ctx context.Context, req *campaign.DeleteCampaignByIDRequest) (*campaign.DeleteCampaignByIDResponse, error)
	UpdateCampaignByID(ctx context.Context, req *campaign.UpdateCampaignByIDRequest) (*campaign.UpdateCampaignByIDResponse, error)
	GetCampaignsByUserID(ctx context.Context, req *campaign.GetCampaignsByUserIDRequest) (*campaign.GetCampaignsByUserIDResponse, error)
	MarkCampaignCompleted(ctx context.Context) 
}

// campaignService is the struct implementation of CampaignService
type campaignService struct {
	campaign.UnimplementedCampaignServiceServer
	campaignRepo repository.CampaignRepository
}

// NewCampaignService initializes and returns a new campaignService instance with a given Campaign repository
func NewCampaignService(campaignRepo repository.CampaignRepository) *campaignService {
	return &campaignService{campaignRepo: campaignRepo}
}

func (s *campaignService) CreateCampaign(ctx context.Context, req *campaign.CreateCampaignRequest) (*campaign.CreateCampaignResponse, error) {
	// Create a new uuid
	uuid := uuid.New()

	// Prepare a struct for campaign
	campaignPayload := models.CampaignDB{
		ID:           uuid.String(),
		UserID:       req.UserId,
		Title:        req.Title,
		Description:  req.Description,
		TargetAmount: req.TargetAmount,
		Deadline:     req.Deadline.AsTime(),
		Category:     helper.MapCategoryDB(int32(req.Category)),
		MinDonation:  req.MinDonation,
	}

	err := helper.ValidateCampaign(campaignPayload)
	if err != nil {
		return nil, err
	}

	// Insert campaign to database
	campaignInterface, err := s.campaignRepo.CreateCampaign(campaignPayload)
	if err != nil {
		return nil, err
	}

	// Cast the campaignInterface type to models.CampaignDB
	createdCampaign, ok := campaignInterface.(models.CampaignDB)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to cast created campaign")
	}

	res := &campaign.CreateCampaignResponse{
		CreatedCampaign: []*campaign.Campaign{
			{
				Id:              createdCampaign.ID,
				UserId:          createdCampaign.UserID,
				Title:           createdCampaign.Title,
				Description:     createdCampaign.Description,
				TargetAmount:    createdCampaign.TargetAmount,
				CollectedAmount: createdCampaign.CollectedAmount,
				Deadline:        timestamppb.New(createdCampaign.Deadline),
				Status:          campaign.CampaignStatus(helper.MapStatusProto(createdCampaign.Status)),
				Category:        campaign.CampaignCategory(helper.MapCateogryProto(createdCampaign.Category)),
				MinDonation:     createdCampaign.MinDonation,
				CreatedAt:       timestamppb.New(createdCampaign.CreatedAt),
				UpdatedAt:       timestamppb.New(createdCampaign.UpdatedAt),
			},
		},
	}
	// Marshal to protobuf binary
	data, err := proto.Marshal(res)
	if err != nil {
		log.Error(err)
	}
	mq.PublishCampaign(data, "campaign.created")

	return res, nil
}

func (s *campaignService) GetCampaignByID(ctx context.Context, req *campaign.GetCampaignByIDRequest) (*campaign.GetCampaignByIDResponse, error) {
	// Get campaign by id
	campaignInterface, err := s.campaignRepo.GetCampaignByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Cast the campaignInterface type to models.CampaignDB
	getCampaign, ok := campaignInterface.(models.CampaignDB)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to cast campaign")
	}

	res := &campaign.GetCampaignByIDResponse{
		Campaign: []*campaign.Campaign{
			{
				Id:              getCampaign.ID,
				UserId:          getCampaign.UserID,
				Title:           getCampaign.Title,
				Description:     getCampaign.Description,
				TargetAmount:    getCampaign.TargetAmount,
				CollectedAmount: getCampaign.CollectedAmount,
				Deadline:        timestamppb.New(getCampaign.Deadline),
				Status:          campaign.CampaignStatus(helper.MapStatusProto(getCampaign.Status)),
				Category:        campaign.CampaignCategory(helper.MapCateogryProto(getCampaign.Category)),
				MinDonation:     getCampaign.MinDonation,
				CreatedAt:       timestamppb.New(getCampaign.CreatedAt),
				UpdatedAt:       timestamppb.New(getCampaign.UpdatedAt),
			},
		},
	}

	return res, nil
}

func (s *campaignService) DeleteCampaignByID(ctx context.Context, req *campaign.DeleteCampaignByIDRequest) (*campaign.DeleteCampaignByIDResponse, error) {
	// Delete campaign by id
	userId, err := s.campaignRepo.DeleteCampaignByID(req.Id)
	if err != nil {
		return nil, err
	}

	msg := campaign.Notification{Id: req.Id, UserId: userId}
	data, err := proto.Marshal(&msg)
	if err != nil {
		log.Error(err)
	}
	mq.PublishCampaign(data, "campaign.deleted")

	return &campaign.DeleteCampaignByIDResponse{
		DeleteResponse: &emptypb.Empty{},
	}, nil
}

func (s *campaignService) UpdateCampaignByID(ctx context.Context, req *campaign.UpdateCampaignByIDRequest) (*campaign.UpdateCampaignByIDResponse, error) {
	var deadline time.Time
	if req.Deadline != nil {
		deadline = req.Deadline.AsTime()
	}

	// Check for user id or campaign id
	if req.Id == "" || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "Please input valid User ID or Campaign ID")
	}
	// Prepare a struct for campaign
	campaignPayload := models.CampaignDB{
		Title:        req.Title,
		Description:  req.Description,
		TargetAmount: req.TargetAmount,
		Deadline:     deadline,
		Status:       helper.MapStatusDB(int32(req.Status)),
		Category:     helper.MapCategoryDB(int32(req.Category)),
		MinDonation:  req.MinDonation,
	}

	err := helper.ValidateUpdateCampaign(campaignPayload)
	if err != nil {
		return nil, err
	}

	// Check if user is trying to update status manually to completed
	if campaignPayload.Status == "completed" {
		return nil, status.Error(codes.PermissionDenied, "You cannot manually set status to COMPLETED")
	}

	// Update campaign by id
	campaignInterface, err := s.campaignRepo.UpdateCampaignByID(req.Id, req.UserId, campaignPayload)
	if err != nil {
		return nil, err
	}

	// Cast the campaignInterface type to models.CampaignDB
	updatedCampaign, ok := campaignInterface.(models.CampaignDB)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to cast updated campaign")
	}

	return &campaign.UpdateCampaignByIDResponse{
		UpdatedCampaign: []*campaign.Campaign{
			{
				Id:              updatedCampaign.ID,
				UserId:          updatedCampaign.UserID,
				Title:           updatedCampaign.Title,
				Description:     updatedCampaign.Description,
				TargetAmount:    updatedCampaign.TargetAmount,
				CollectedAmount: updatedCampaign.CollectedAmount,
				Deadline:        timestamppb.New(updatedCampaign.Deadline),
				Status:          campaign.CampaignStatus(helper.MapStatusProto(updatedCampaign.Status)),
				Category:        campaign.CampaignCategory(helper.MapCateogryProto(updatedCampaign.Category)),
				MinDonation:     updatedCampaign.MinDonation,
				CreatedAt:       timestamppb.New(updatedCampaign.CreatedAt),
				UpdatedAt:       timestamppb.New(updatedCampaign.UpdatedAt),
			},
		},
	}, nil
}

func (s *campaignService) GetCampaignsByUserID(ctx context.Context, req *campaign.GetCampaignsByUserIDRequest) (*campaign.GetCampaignsByUserIDResponse, error) {
	// Get campaign by user id
	campaignInterface, err := s.campaignRepo.GetCampaignsByUserID(req.UserId)
	if err != nil {
		return nil, err
	}

	// Cast the campaignInterface type to models.CampaignDB
	getCampaign, ok := campaignInterface.([]models.CampaignDB)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to cast campaign")
	}

	var campaignList []*campaign.Campaign
	for _, val := range getCampaign {
		campaignList = append(campaignList, &campaign.Campaign{
			Id:              val.ID,
			UserId:          val.UserID,
			Title:           val.Title,
			Description:     val.Description,
			TargetAmount:    val.TargetAmount,
			CollectedAmount: val.CollectedAmount,
			Deadline:        timestamppb.New(val.Deadline),
			Status:          campaign.CampaignStatus(helper.MapStatusProto(val.Status)),
			Category:        campaign.CampaignCategory(helper.MapCateogryProto(val.Category)),
			MinDonation:     val.MinDonation,
			CreatedAt:       timestamppb.New(val.CreatedAt),
			UpdatedAt:       timestamppb.New(val.UpdatedAt),
		},
		)
	}

	return &campaign.GetCampaignsByUserIDResponse{
		Campaign: campaignList,
	}, nil
}

func (s *campaignService) MarkCampaignCompleted(ctx context.Context) {
	campaignInterface, err := s.campaignRepo.UpdateCampaignToCompleted()
	if err != nil {
		log.Printf("%v",err)
	}

	// Cast the campaignInterface type to models.CampaignDB
	updatedCampaign, ok := campaignInterface.([]models.CampaignDB)
	if len(updatedCampaign)!=0 && !ok {
		log.Printf("failed to cast campaign")
	}


	for _, val := range updatedCampaign {
		log.Printf("Campaign ID: %s is updated, current status: %s", val.ID, val.Status)
	}
}
