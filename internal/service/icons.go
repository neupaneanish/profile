package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func (s *RootProfileService) Icons(
	ctx context.Context,
	_ *rootProfilev1.IconsRequest,
) (*rootProfilev1.IconsResponse, error) {
	serviceName := "Icons"

	rows, err := s.cfg.Repository.Icons(ctx)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Failed to query icons", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Icon, len(rows))
	for i, p := range rows {
		res[i] = &profilev1.Icon{
			Id:   p.ID.String(),
			Name: p.Name,
			Icon: p.Icon,
		}
	}

	return &rootProfilev1.IconsResponse{Icons: res}, nil
}
