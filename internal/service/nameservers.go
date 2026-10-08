package service

import (
	"context"

	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func (s *RootProfileService) Nameservers(
	ctx context.Context,
	_ *profilev1.NameserversRequest,
) (*profilev1.NameserversResponse, error) {
	serviceName := "Nameservers"

	rows, err := s.cfg.Repository.Nameservers(ctx)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Failed to fetch nameservers", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.Nameservers, len(rows))

	for i, row := range rows {
		res[i] = &profilev1.Nameservers{
			Id:     row.ID.String(),
			Server: row.Server,
		}
	}

	return &profilev1.NameserversResponse{Nameservers: res}, nil
}
