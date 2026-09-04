package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProgatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateExperience(
	ctx context.Context,
	req *gatewayProgatewayProfilev1.CreateExperienceRequest,
) (*gatewayProgatewayProfilev1.CreateExperienceResponse, error) {
	serviceName := "CreateExperience"
	userSession := utils.UserSessionContext(ctx)

	params := &repository.CreateExperienceParams{
		UserID:       userSession.UserID,
		Title:        req.GetTitle(),
		CompanyName:  req.GetCompanyName(),
		Location:     req.GetLocation(),
		LocationType: enum.LocationType(req.GetLocationType()),
		StartDate:    req.GetStartDate().AsTime(),
		EndDate:      timestampValue(req.GetEndDate()),
		Description:  stringValue(req.GetDescription()),
		CreatedBy:    userSession.UserID,
		UpdatedBy:    userSession.UserID,
	}

	id, err := s.cfg.Repository.CreateExperience(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Create Experience Failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}
	return &gatewayProgatewayProfilev1.CreateExperienceResponse{Id: id.String()}, nil
}
