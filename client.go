package main

import (
	"context"
	"fmt"
	"time"

	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/indikay/go-core/middleware/jwt"
	"github.com/indikay/notification-service/api/notifications"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	con, _ := kgrpc.DialInsecure(
		context.Background(),
		kgrpc.WithEndpoint("localhost:9000"), // wallet-service.stag.svc.cluster.local:9000
		// kgrpc.WithMiddleware(
		// 	jwt.Client(jwt.WithClaims(func() jwtlib.Claims {
		// 		return jwtlib.RegisteredClaims{Subject: "1eb1a218-62cf-43dc-9d15-f1f0318c7813", ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Hour))}
		// 	})),
		// ),
		kgrpc.WithTimeout(time.Second*60*5),
	)

	ctx := context.Background()
	ctx, _ = jwt.ClientGrpcAuth(ctx, jwt.WithClaims(func() jwtlib.Claims {
		return jwtlib.RegisteredClaims{Subject: "1eb1a218-62cf-43dc-9d15-f1f0318c7813", ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Hour))}
	}))

	notiCli := notifications.NewNotificationClient(con)
	notiResp, err := notiCli.GetNotificationSettings(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(notiResp)
	}

}
