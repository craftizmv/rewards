package main

import (
	"context"
	"github.com/craftizmv/rewards/config"
	consumers2 "github.com/craftizmv/rewards/internal/app/handlers/consumers"
	"github.com/craftizmv/rewards/internal/app/usecase"
	"github.com/craftizmv/rewards/internal/data/infrastructure/cache"
	"github.com/craftizmv/rewards/internal/data/infrastructure/database"
	"github.com/craftizmv/rewards/internal/data/infrastructure/external/proxies"
	"github.com/craftizmv/rewards/internal/data/infrastructure/external/service"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue/consumers"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue/events"
	"github.com/craftizmv/rewards/internal/data/infrastructure/queue/publisher"
	repository_impl "github.com/craftizmv/rewards/internal/data/repository-impl"
	"github.com/craftizmv/rewards/internal/domain/entities"
	"github.com/craftizmv/rewards/pkg/logger"
	"github.com/craftizmv/rewards/server"
	"go.uber.org/zap"
	"time"
)

func main() {

	appCtx := context.Background()

	// Initialising the config
	cfg := config.GetConfig()

	// init logger
	log := logger.InitLogger(cfg.Logger)

	// init concrete cache
	redisCache := cache.NewRedisCache[entities.Order](cfg.CacheCfg, log)

	// init DB
	postgresDB, err := database.NewPostgresDB(cfg.DBCfg, log)
	if err != nil {
		log.Error("Could not initialize PostgreSQL:", err)
		// Handle the error appropriately (e.g., exit the application)
		return
	}
	defer func() {
		if err := postgresDB.Close(); err != nil {
			log.Error("Error closing PostgreSQL connection:", err)
		}
	}()

	// Checking whether postgres is working or not.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var currentTime string
	err = postgresDB.QueryRowContext(ctx, "SELECT NOW()").Scan(&currentTime)
	if err != nil {
		log.Error("Failed to execute query:", err)
	} else {
		log.Info("Current time from PostgreSQL:", currentTime)
	}

	rewardRepo := repository_impl.NewPostgresRewardRepository(postgresDB)
	campaignProxy := proxies.NewCampaignProxy()
	inventoryProxy := proxies.NewInventoryProxy()
	mailer := service.NewSomeConcreteMailer()
	emailProxy := proxies.NewEmailProxy(mailer, log)

	shipper := service.NewLogiDeli()
	shippingProxy := proxies.NewShippingProxy(shipper)
	userProxy := proxies.NewUserProxy()
	orderProxy := proxies.NewOrderProxy()

	rewardProxies := &usecase.RewardProxies{
		CampaignProxy:  campaignProxy,
		InventoryProxy: inventoryProxy,
		EmailProxy:     emailProxy,
		ShippingProxy:  shippingProxy,
		UserProxy:      userProxy,
		OrderProxy:     orderProxy,
	}
	rewardUseCase := usecase.NewRewardUseCaseImpl(redisCache, rewardRepo, log, rewardProxies)

	// start the echo server.
	server.NewEchoServer(cfg.EchoCfg, log, rewardUseCase).Start()

	// init RabbitMQ
	conn, err := queue.NewRabbitMQConn(cfg.Rabbitmq, appCtx)
	if err != nil {
		log.Error("Failed to create RabbitMQ connection", "err", zap.Error(err))
		panic(err)
	}
	defer conn.Close()

	// init the consumer for RabbitMQ
	// TODO-MV : May be pass a producer to reproduce the message.

	// we need to inject the obj having usecase repo - which can talk to redis and DB layer via repository.
	// using below object and interfaces consumer should be able to talk to usecase layer via the inversion of control

	// create rabbitMQ publisher which will help with re-allocating when an order is cancelled.
	pub := publisher.NewPublisher(ctx, cfg.Rabbitmq, conn, log)

	eligibleOrder := queue.OrderDeliveryBase{
		Ctx:          ctx,
		Log:          log,
		Cfg:          cfg,
		ConnRabbitmq: conn,
		GiftUseCases: rewardUseCase,
		Publisher:    pub,
	}

	allocateOrderConsumer := consumers.NewOrderConfirmedConsumer[*queue.OrderDeliveryBase](ctx, cfg.Rabbitmq, conn, log, consumers2.HandleAllocateReward)
	allocateOrderFromBufferConsumer := consumers.NewOrderConfirmedBufferConsumer[*queue.OrderDeliveryBase](ctx, cfg.Rabbitmq, conn, log, consumers2.HandleAllocateFromBufferReward)
	cancelOrderConsumer := consumers.NewOrderCancelledConsumer[*queue.OrderDeliveryBase](ctx, cfg.Rabbitmq, conn, log, consumers2.HandleCancelReward)
	go func() {
		e := allocateOrderConsumer.ConsumeMessage(events.AllocateReward{}, &eligibleOrder)
		if e != nil {
			log.Error("Failed to consume order:", "err", e)
		}
	}()

	go func() {
		// event on this queue will also be AllocateReward event
		e := allocateOrderFromBufferConsumer.ConsumeMessage(events.AllocateReward{}, &eligibleOrder)
		if e != nil {
			log.Error("Failed to consume order:", "err", e)
		}
	}()

	go func() {
		e := cancelOrderConsumer.ConsumeMessage(events.RevokeReward{}, &eligibleOrder)
		if e != nil {
			log.Error("Failed to consume order:", "err", e)
		}
	}()
}
