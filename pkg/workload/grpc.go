package workload

import (
	"context"
	"fmt"
	"net"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCService pairs a service descriptor with its implementation so they can
// be registered together with WithGRPC.
type GRPCService struct {
	Desc *grpc.ServiceDesc
	Impl any
}

// GRPCOptions configures the gRPC workload added by WithGRPC.
type GRPCOptions struct {
	Port             *uint // TCP port; defaults to 50051
	EnableReflection bool  // registers the gRPC server-reflection service
}

// WithGRPC adds a gRPC server to the workload group.  All provided services are
// registered before the server starts.  The server is stopped gracefully when
// any other workload in the group returns.
func (s *Server) WithGRPC(_ context.Context, options GRPCOptions, services ...GRPCService) *Server {
	grpcServer := grpc.NewServer()
	if options.EnableReflection {
		reflection.Register(grpcServer)
	}
	for _, service := range services {
		grpcServer.RegisterService(service.Desc, service)
	}

	return s.add(func() error { //nolint:contextcheck// this is by design of oklog/run
		lc := net.ListenConfig{}
		lis, err := lc.Listen(
			context.Background(),
			"tcp", fmt.Sprintf(":%d", lo.FromPtrOr(options.Port, 50051)),
		)
		if err != nil {
			return fmt.Errorf("failed to listen for gRPC: %w", err)
		}

		return grpcServer.Serve(lis)
	}, func(interruptReason error) {
		s.logger.Warn("Interrupting gRPC server", "err", interruptReason)

		grpcServer.GracefulStop()
	})
}
