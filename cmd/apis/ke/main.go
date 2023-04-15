package main

import (
	"fmt"
	"log"
	"net"
	"os"

	v1alphaservice "github.com/nokamoto/2pf23/internal/service/ke/v1alpha"
	v1alphaapi "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// Run the ke server.
//
// GRPC_PORT: The port number to listen.
// DEBUG: Enable debug logging if set to non-empty.
func main() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if os.Getenv("DEBUG") != "" {
		cfg.Level.SetLevel(zap.DebugLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "9000"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	v1alphaapi.RegisterKeServiceServer(s, v1alphaservice.NewService(logger))

	logger.Info("server listening", zap.String("address", lis.Addr().String()))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
