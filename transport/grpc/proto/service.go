package proto

import (
	"context"
	"log/slog"
	"strconv"

	pb "server/transport/grpc/proto/com.opswat.mem.fusion.account"
	"server/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn    *grpc.ClientConn
	account *pb.AccountServiceClient
	Url     string
}

func (c *GrpcClient) Connect() {
	conn, err := grpc.NewClient(c.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.FailOnError(err, "Failed to connect to fusion account grpc server")
	slog.Info("Connected to fusion account api")

	grpcClient := pb.NewAccountServiceClient(conn)

	c.conn = conn
	c.account = &grpcClient
}

func (c *GrpcClient) Close() {
	utils.Close(c.conn)
}

func (client GrpcClient) GetDataRetentionMillis(ctx context.Context, accountId string) (string, error) {
	dto, err := (*client.account).FindByID(ctx, &pb.AccountIDParam{Value: accountId})
	slog.Info("Received response from grpc server", "dto", dto)
	if err != nil {
		return "", err
	}
	if dto == nil {
		return "2592000000", nil
	}
	accConfig := dto.Config
	if accConfig == nil {
		return "2592000000", nil
	}
	dataRetentionInDay := accConfig.DataRetention
	if dataRetentionInDay < 1 {
		return "2592000000", nil
	}
	millis := int64(dataRetentionInDay) * 86_400_000
	return strconv.FormatInt(millis, 10), nil
}
