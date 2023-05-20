package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/nokamoto/2pf23/internal/app"
	"github.com/nokamoto/2pf23/internal/app/ke"
	"github.com/nokamoto/2pf23/internal/ent"
	keinfra "github.com/nokamoto/2pf23/internal/infra/postgresql/ke"
	v1alphaservice "github.com/nokamoto/2pf23/internal/server/generated/ke/v1alpha"
	"github.com/nokamoto/2pf23/pkg/api/ke/v1alpha/kev1alphaconnect"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type runtime struct {
	app.ResourceIDGenerator
	*keinfra.Cluster
}

func setupLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if os.Getenv("DEBUG") != "" {
		cfg.Level.SetLevel(zap.DebugLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	return logger
}

func envOr(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func setupEnt(logger *zap.Logger) *ent.Client {
	host := envOr("POSTGRESQL_HOST", "postgresql.default.svc.cluster.local")
	user := envOr("POSTGRESQL_USER", "postgres")
	dbname := envOr("POSTGRESQL_DBNAME", "ke")
	password := envOr("POSTGRESQL_PASSWORD", "local")
	dataSource := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", host, user, dbname, password)
	logger.Debug("connecting to postgres", zap.String("host", host), zap.String("dbname", dbname), zap.String("user", user))
	client, err := ent.Open("postgres", dataSource)
	if err != nil {
		logger.Fatal("failed opening connection to postgres", zap.Error(err))
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		logger.Fatal("failed creating schema resources", zap.Error(err))
	}
	return client
}

// Run the ke server.
//
// GRPC_PORT: The port number to listen. If not set, 9000 is used.
// DEBUG: Enable debug logging if set to non-empty.
// POSTGRESQL_HOST: The host name of the postgresql server. If not set, postgresql.default.svc.cluster.local is used.
// POSTGRESQL_USER: The user name of the postgresql server. If not set, postgres is used.
// POSTGRESQL_DBNAME: The database name of the postgresql server. If not set, ke is used.
// POSTGRESQL_PASSWORD: The password of the postgresql server. If not set, local is used.
func main() {
	logger := setupLogger()
	defer func() {
		_ = logger.Sync()
	}()

	client := setupEnt(logger)
	defer client.Close()

	rt := &runtime{
		Cluster: keinfra.NewCluster(client),
	}
	app := ke.NewCluster(rt)
	svc := v1alphaservice.NewService(logger, app)

	port := envOr("GRPC_PORT", "9000")
	mux := http.NewServeMux()
	path, handler := kev1alphaconnect.NewKeServiceHandler(svc)
	mux.Handle(path, handler)
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", port),
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		logger.Fatal("failed to listen and serve", zap.Error(err))
	}
}
