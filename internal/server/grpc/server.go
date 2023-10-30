package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	proto_interfaces "github.com/bobgromozeka/metrics/internal/proto-interfaces"
	"github.com/bobgromozeka/metrics/internal/server/storage"
)

type Config struct {
	Addr           string
	TrustedSubnet  string // TODO
	PrivateKeyPath string
	CertPath       string
}

func Start(ctx context.Context, c Config, stor storage.Storage) error {
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		return err
	}

	creds, err := credentials.NewServerTLSFromFile(c.CertPath, c.PrivateKeyPath)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	proto_interfaces.RegisterMetricsServer(grpcServer, NewMetricsService(stor))

	go func() {
		<-ctx.Done()
		fmt.Println("Stopping grpc server......")
		grpcServer.GracefulStop()
	}()

	fmt.Printf("Starting gRPC server on addr: [%s]......\n", c.Addr)

	return grpcServer.Serve(lis)
}
