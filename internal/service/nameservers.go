package service

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func (s *RootProfileService) NameServers(
	ctx context.Context,
	_ *profilev1.NameServersRequest,
) (*profilev1.NameServersResponse, error) {
	serviceName := "Nameservers"

	rows, err := s.cfg.Repository.NameServers(ctx)
	if err != nil {
		s.cfg.Logger.ErrorContext(ctx, "Failed to fetch nameservers", "service", serviceName, "error", err)
		return nil, errs.ErrInternalServer
	}

	res := make([]*profilev1.NameServers, len(rows))

	for i, n := range rows {
		res[i] = &profilev1.NameServers{
			Id:        n.ID.String(),
			Cname:     n.Cname,
			Domain:    n.Domain,
			Active:    n.Active,
			CreatedAt: timestamppb.New(n.CreatedAt),
			CreatedBy: n.CreatedBy.String(),
			UpdatedAt: timestamppb.New(n.UpdatedAt),
			UpdatedBy: n.UpdatedBy.String(),
		}
	}

	return &profilev1.NameServersResponse{Nameservers: res}, nil
}
