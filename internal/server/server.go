package server

import (
	"context"
	"fmt"
	"github.com/Seann-Moser/BaseGoAPI/internal/configuration"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func Serve(ciCtx *cli.Context) error {
	config, err := configuration.LoadConfig()
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ciCtx.Context)

	go func() {
		osCall := <-c
		config.Logger.Info(fmt.Sprintf("system call:%+v", osCall))
		cancel()
	}()

	if err := NewEndpoints(ctx, config).StartServer(); err != nil {
		config.Logger.Error("failed to serve", zap.Error(err))
	}
	return nil
}
