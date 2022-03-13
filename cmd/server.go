package commands

import (
	"baal/configs"
	"baal/controllers"
	"baal/libs/logger"
	"baal/routers"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Ball server for localhost",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return nil
		}

		_, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			os.Setenv("PORT", args[0])
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		app := fx.New(
			fx.NopLogger,
			configs.Module,
			logger.Module,
			controllers.Module,
			routers.Module,
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

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serverStart(lc fx.Lifecycle, r *routers.Router, log *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s, port := r.Serve()
			log.Info(fmt.Sprintf("Server start on : %s port", port))

			go s.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Server shotdown")
			return nil
		},
	})
}
