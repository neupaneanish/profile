package service

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func stringValue(wrapper *wrapperspb.StringValue) *string {
	if wrapper == nil {
		return nil
	}
	value := wrapper.GetValue()
	return &value
}

func timestampValue(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func timestamppbValue(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func stringpbValue(s *string) *wrapperspb.StringValue {
	if s == nil {
		return nil
	}
	return wrapperspb.String(*s)
}
