package main

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"yandex_GophKeeper_client/config"
	"yandex_GophKeeper_client/internal/app/cli"
	grpc_requesters "yandex_GophKeeper_client/internal/app/requesters/gRPC"
	"yandex_GophKeeper_client/internal/app/requesters/gRPC/proto"
	http_requesters "yandex_GophKeeper_client/internal/app/requesters/http"
)

func main() {
	conf := config.AppConfig{}
	err := conf.Configure()
	if err != nil {
		log.Fatal(err)
	}

	//build a zap logger
	level, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}
	zCfg := zap.Config{}
	if level == zap.NewAtomicLevelAt(zap.InfoLevel) {
		zCfg = zap.NewProductionConfig()
		zCfg.DisableCaller = true
	} else {
		zCfg = zap.NewDevelopmentConfig()
		zCfg.DisableStacktrace = true
	}
	zCfg.OutputPaths = []string{"gophKeeper.log"}
	zCfg.ErrorOutputPaths = []string{"gophKeeperErr.log"}

	logger, err := zCfg.Build()
	if err != nil {
		log.Fatalf("failed to build logger: %v", err)
	}
	sugar := logger.Sugar()
	sugar.Debug("logger initialized")

	//prepare an HTTP client
	var fullApiPath string
	if conf.UseHTTPS {
		fullApiPath = "https://"
	} else {
		sugar.Info("GophKeeper will use UNPROTECTED connection to the server. Please DONT USE REAL sensitive data.")
		fullApiPath = "http://"
	}
	fullApiPath += conf.APIAddress

	httpClient := http.DefaultClient

	httpRequester := http_requesters.NewRequester(fullApiPath, httpClient, "")
	//jwt will be set after user`s authorisation.

	//prepare gRPC client
	grpcClient := &grpc.ClientConn{}
	if conf.UseHTTPS {
		creds := credentials.NewClientTLSFromCert(nil, "")
		grpcClient, err = grpc.NewClient(conf.GRPCAddress, grpc.WithTransportCredentials(creds))
		if err != nil {
			sugar.Fatal("failed to initialize gRPC client", zap.Error(err))
		}
	} else {
		grpcClient, err = grpc.NewClient(conf.GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			sugar.Fatal("failed to initialize gRPC client", zap.Error(err))
		}
	}

	grpcRequester := grpc_requesters.NewGRPCRequester(proto.NewGophKeeperServiceClient(grpcClient), "", conf.MaxBinDataChunkSize, sugar)

	//build cli controller
	cliController := cli.NewCommandsController(conf, httpRequester, grpcRequester, sugar)

	//run
	sugar.Info("Starting app...")
	sugar.Sync()
	cliController.Run(context.Background())
}
