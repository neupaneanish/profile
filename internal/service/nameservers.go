package service

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
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

	for i, n := range rows {
		res[i] = &profilev1.Nameservers{
			Id:        n.ID.String(),
			Ip:        n.Ip,
			IpType:    n.IpType,
			CreatedAt: timestamppb.New(n.CreatedAt),
			CreatedBy: n.CreatedBy.String(),
			UpdatedAt: timestamppb.New(n.UpdatedAt),
			UpdatedBy: n.UpdatedBy.String(),
		}
	}

	return &profilev1.NameserversResponse{Nameservers: res}, nil
}
