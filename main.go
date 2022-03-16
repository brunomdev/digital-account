package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/brunomdev/digital-account/app/api"
	"github.com/brunomdev/digital-account/config"
	"github.com/brunomdev/digital-account/domain"
	"github.com/brunomdev/digital-account/domain/account"
	"github.com/brunomdev/digital-account/domain/operationtype"
	"github.com/brunomdev/digital-account/domain/transaction"
	"github.com/brunomdev/digital-account/infra/log"
	repo "github.com/brunomdev/digital-account/infra/mysql"
	"github.com/brunomdev/digital-account/infra/newrelic"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(ctx, "unable to load configuration", err)
	}

	newRelic, err := newrelic.NewNewRelic(cfg)
	if err != nil {
		log.Fatal(ctx, "unable to connect to newrelic", err)
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBDatabase)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(ctx, "unable connect with database", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(ctx, "unable to get mysql driver", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.DBDatabase,
		driver,
	)
	if err != nil {
		log.Fatal(ctx, "unable to define initiate migrate", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(ctx, "unable to migrate database", err)
	}

	accountRepo := repo.NewAccountRepository(db)
	accountSvc := account.NewService(accountRepo)
	opTypeRepo := repo.NewOperationTypeRepository(db)
	opTypeSvc := operationtype.NewService(opTypeRepo)
	transactionRepo := repo.NewTransactionRepository(db)
	transactionSvc := transaction.NewService(transactionRepo, accountSvc, opTypeSvc)

	service := &domain.Service{
		Account:       accountSvc,
		OperationType: opTypeSvc,
		Transaction:   transactionSvc,
	}

	srv, err := api.NewServer(
		api.WithConfig(cfg),
		api.WithService(service),
		api.WithNewRelic(newRelic),
	)
	if err != nil {
		log.Fatal(ctx, "new server: ", err)
	}

	<-ctx.Done()

	stop()

	log.Info(ctx, "shutting down gracefully")

	err = srv.Close()
	if err != nil {
		log.Error(ctx, "forced server to shutdown: ", err)
	}

	err = db.Close()
	if err != nil {
		log.Error(ctx, "forced db to shutdown: ", err)
	}

	newRelic.Shutdown(time.Second * 10)

	log.Info(ctx, "exiting")

	err = log.Close()
	if err != nil {
		fmt.Printf("forced log to shutdown: %v", err)
	}
}
