package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/bobgromozeka/metrics/internal/server/db"
	grpcServer "github.com/bobgromozeka/metrics/internal/server/grpc"
	httpServer "github.com/bobgromozeka/metrics/internal/server/http"
	"github.com/bobgromozeka/metrics/internal/server/storage"
)

func Start(ctx context.Context, startupConfig StartupConfig) {
	s, storageStoppedChan := createStorage(ctx, startupConfig)

	privateKey, readErr := os.ReadFile(startupConfig.PrivateKeyPath)
	if readErr != nil {
		log.Fatalln(readErr)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := httpServer.Start(
			ctx, httpServer.Config{
				Addr:          startupConfig.HttpAddr,
				PrivateKey:    privateKey,
				TrustedSubnet: startupConfig.TrustedSubnet,
				HashKey:       startupConfig.HashKey,
			}, s,
		)
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := grpcServer.Start(
			ctx, grpcServer.Config{
				Addr:           startupConfig.GRPCAddr,
				TrustedSubnet:  startupConfig.TrustedSubnet,
				PrivateKeyPath: startupConfig.GRPCPrivateKeyPath,
				CertPath:       startupConfig.GRPCCertPath,
			}, s,
		)
		log.Fatalln(err)
	}()

	wg.Wait()            // Wait for both servers to stop
	<-storageStoppedChan // Wait for storage stop (in database case)
}

func createStorage(ctx context.Context, sc StartupConfig) (storage.Storage, <-chan struct{}) {
	var s storage.Storage
	stoppedChan := make(chan struct{}, 1)

	if sc.DatabaseDsn != "" {
		connErr := db.Connect(sc.DatabaseDsn)
		if connErr != nil {
			panic(connErr)
		}

		ddlErr := storage.Bootstrap(db.Connection())
		if ddlErr != nil {
			panic(ddlErr)
		}
		s = storage.NewPG(db.Connection())

		go func() {
			defer func() {
				stoppedChan <- struct{}{}
			}()
			<-ctx.Done()

			hardCtx, hardCancel := context.WithTimeout(context.Background(), time.Second*15)
			defer hardCancel()

			db.Connection().Close(hardCtx)
		}()

	} else {
		s = storage.NewMemory()
		s = storage.NewPersistenceStorage(
			s, storage.PersistenceSettings{
				Path:     sc.FileStoragePath,
				Interval: sc.StoreInterval,
				Restore:  sc.Restore,
			},
		)
		stoppedChan <- struct{}{} // We do not need to wait for any signals when using memo storage
	}

	return s, stoppedChan
}
