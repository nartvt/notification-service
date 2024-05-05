package extsvc

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/indikay/notification-service/internal/conf"
	pbProfile "github.com/indikay/profile-service/api/profile/v1"
)

type ProfileService struct {
	profileClient pbProfile.ProfileServiceClient
}

func NewProfileService(serverConfig *conf.Data) *ProfileService {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(serverConfig.Profile.Addr),
		grpc.WithTimeout(serverConfig.Profile.Timeout.AsDuration()),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	)

	if err != nil {
		panic(err)
	}

	client := pbProfile.NewProfileServiceClient(conn)
	return &ProfileService{
		profileClient: client,
	}
}

func (s *ProfileService) GetProfileByEmail(ctx context.Context, email string) (*pbProfile.GetUserProfileResponse, error) {
	resp, err := s.profileClient.GetUserProfileInternal(ctx, &pbProfile.GetUserProfileInternalRequest{Email: email})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
