package service

import (
	"context"
	"time"

	gatewayProgatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/utils"
)

func (s *GatewayProfileService) CreateExperience(
	ctx context.Context,
	req *gatewayProgatewayProfilev1.CreateExperienceRequest,
) (*gatewayProgatewayProfilev1.CreateExperienceResponse, error) {
	serviceName := "CreateExperience"
	userSession := utils.UserSessionContext(ctx)

	id, err := createUpdateExperience(
		ctx,
		"",
		userSession.UserID,
		userSession.UserID,
		req.GetExperience(),
		time.Time{},
		s.cfg.Repository,
		s.cfg.Logger,
		serviceName,
	)
	if err != nil {
		return nil, err
	}

	return &gatewayProgatewayProfilev1.CreateExperienceResponse{Id: id.String()}, nil
}
