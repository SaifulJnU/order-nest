package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/order-nest/config"
	"github.com/order-nest/internal/delivery/http"
	appLogger "github.com/order-nest/pkg/logger"
	"github.com/spf13/cobra"
)

// serveCmd represents the CLI command to start the HTTP server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the HTTP server",
	Long:  "Starts the HTTP server for the Order Nest application",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a context that listens for termination signals (SIGINT, SIGTERM)
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		defer stop()

		// Bootstrap dependencies and get router setup
		router, err := http.Bootstrap(ctx)
		if err != nil {
			return fmt.Errorf("failed to bootstrap server: %w", err)
		}

		srv := http.NewHTTPServer(ctx, router)

		// Channel to capture server errors
		srvErr := make(chan error, 1)

		// Start the server in a separate goroutine
		go func() {
			appLogger.L().WithField("port", config.GetConfig().Port).Info("http server listening")
			srvErr <- srv.ListenAndServe()
		}()

		// Wait for either a server error or OS termination signal
		select {
		case err = <-srvErr:
			return err
		case <-ctx.Done():
			appLogger.L().Info("received termination signal, shutting down")
			stop()
		}

		// Gracefully shutdown the server
		return srv.Shutdown(context.Background())
	},
}

func init() {
	// Attach the serve command to the root command
	rootCmd.AddCommand(serveCmd)
}
