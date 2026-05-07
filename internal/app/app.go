package app

import (
	"context"
	"fmt"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/shutdown"
	shoppinglistpb "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/config"
	remoteServices "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/repository/grpc"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/repository/postgres"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/transport/amqp"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/transport/dependencies/service"
	shoppingList "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"time"
)

func Run(cfg *config.Config) {
	log.InitWithService("shopping-list", *cfg.LogsPath, *cfg.Environment == config.EnvDev)
	cfg.Print()

	ctx := context.Background()

	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		log.LogFatal(ctx, log.Event{
			Event:     "app.startup.failed",
			Message:   "service startup failed",
			Component: "app",
		}, err)
		return
	}

	repository := postgres.NewRepository(db, cfg.ShoppingList)
	grpcServices, err := remoteServices.NewRepository(cfg)
	if err != nil {
		log.LogFatal(ctx, log.Event{
			Event:     "app.startup.failed",
			Message:   "service startup failed",
			Component: "app",
		}, err)
		return
	}

	shoppingListService, err := service.New(ctx, cfg, repository, grpcServices)
	if err != nil {
		log.LogFatal(ctx, log.Event{
			Event:     "app.startup.failed",
			Message:   "service startup failed",
			Component: "app",
		}, err)
		return
	}

	var mqServer *amqp.Server = nil
	if len(*cfg.Amqp.Host) > 0 {
		mqServer, err = amqp.NewServer(cfg.Amqp, shoppingListService.MQ)
		if err != nil {
			return
		}
		if err := mqServer.Start(); err != nil {
			log.LogFatal(ctx, log.Event{
				Event:     "app.startup.failed",
				Message:   "service startup failed",
				Component: "app",
			}, err)
			return
		}
		log.Log(ctx, log.Event{
			Event:     "mq.server.initialized",
			Message:   "mq server initialized",
			Component: log.ComponentAMQP,
		})
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *cfg.Port))
	if err != nil {
		log.LogFatal(ctx, log.Event{
			Event:     "app.startup.failed",
			Message:   "service startup failed",
			Component: "app",
		}, err)
		return
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			log.UnaryServerInterceptor(),
		),
	)

	healthServer := health.NewServer()
	shoppingListServer := shoppingList.NewServer(*shoppingListService, *cfg.ShoppingList.CheckSubscription)

	go monitorHealthChecking(db, healthServer)

	shoppinglistpb.RegisterShoppingListServiceServer(grpcServer, shoppingListServer)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.LogError(ctx, log.Event{
				Event:     "grpc.server.failed",
				Message:   "error occurred while running grpc server",
				Component: log.ComponentGRPC,
			}, err)
		} else {
			log.Log(ctx, log.Event{
				Event:     "grpc.server.started",
				Message:   "grpc server started",
				Component: log.ComponentGRPC,
			})
		}
	}()

	wait := shutdown.Graceful(ctx, 5*time.Second, map[string]shutdown.Operation{
		"grpc-server": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
		"database": func(ctx context.Context) error {
			return db.Close()
		},
		"mq": func(ctx context.Context) error {
			if mqServer == nil {
				return nil
			}
			return mqServer.Stop()
		},
	})
	<-wait
}
