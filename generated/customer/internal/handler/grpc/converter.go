package grpc

import (
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	pb "github.com/makoto-developer/go_microservice_example/proto/customer-service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func domainGenderToProto(gender *domain.Gender) pb.Gender {
	if gender == nil {
		return pb.Gender_GENDER_UNSPECIFIED
	}
	switch *gender {
	case domain.GenderMale:
		return pb.Gender_MALE
	case domain.GenderFemale:
		return pb.Gender_FEMALE
	case domain.GenderOther:
		return pb.Gender_OTHER
	default:
		return pb.Gender_GENDER_UNSPECIFIED
	}
}

func protoGenderToDomain(gender pb.Gender) *domain.Gender {
	if gender == pb.Gender_GENDER_UNSPECIFIED {
		return nil
	}
	var g domain.Gender
	switch gender {
	case pb.Gender_MALE:
		g = domain.GenderMale
	case pb.Gender_FEMALE:
		g = domain.GenderFemale
	case pb.Gender_OTHER:
		g = domain.GenderOther
	default:
		return nil
	}
	return &g
}

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func parseDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func timestampProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func parseOptionalUUID(s string) (*uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func formatOptionalString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
