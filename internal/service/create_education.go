package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateEducation(
	ctx context.Context, req *gatewayProfilev1.CreateEducationRequest,
) (*gatewayProfilev1.CreateEducationResponse, error) {
	serviceName := "CreateEducation"
	userSession := utils.UserSessionContext(ctx)

	params := &repository.CreateEducationParams{
		UserID:        userSession.UserID,
		School:        req.GetSchool(),
		Degree:        req.GetDegree(),
		Affiliation:   stringValue(req.GetAffiliation()),
		FieldOfStudy:  stringValue(req.GetFieldOfStudy()),
		Concentration: stringValue(req.GetConcentration()),
		StartDate:     req.GetStartDate().AsTime(),
		EndDate:       timestampValue(req.GetEndDate()),
		Address:       req.GetAddress(),
		Description:   stringValue(req.GetDescription()),
		CreatedBy:     userSession.UserID,
		UpdatedBy:     userSession.UserID,
	}

	id, err := s.cfg.Repository.CreateEducation(ctx, params)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Create Education Failed", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	return &gatewayProfilev1.CreateEducationResponse{Id: id.String()}, nil
}
