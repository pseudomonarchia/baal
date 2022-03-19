package commands

import (
	"baal/configs"
	"baal/controllers"
	"baal/database"
	"baal/libs/logger"
	"baal/routers"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	port      int            = 7001
	mode      string         = "DEBUG"
	serverCmd *cobra.Command = &cobra.Command{
		Use:   "server",
		Short: "Run Ball server for localhost",
		RunE: func(cmd *cobra.Command, args []string) error {
			setENV()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			app := fx.New(
				fx.NopLogger,
				configs.Module,
				logger.Module,
				controllers.Module,
				routers.Module,
				database.Module,
				fx.Invoke(serverStart),
			)

			err := app.Start(ctx)
			defer app.Stop(ctx)
			if err != nil {
				return err
			}

			<-app.Done()
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&port, "port", "p", 7001, "Server listent on port")
	serverCmd.Flags().StringVarP(&mode, "mode", "m", "debug", "Use debug/release mode")
}

func serverStart(lc fx.Lifecycle, r *routers.Router, log *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s, port := r.Serve()
			log.Info(fmt.Sprintf("Server start on >>> %s port", port))

			go s.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Server shotdown")
			return nil
		},
	})
}

func setENV() {
	os.Setenv("PORT", strconv.Itoa(port))
	os.Setenv("MODE", mode)
}
