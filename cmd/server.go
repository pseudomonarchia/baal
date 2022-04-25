package cmd

import (
	"baal/config"
	"baal/controller"
	"baal/database"
	"baal/lib/logger"
	"baal/router"
	"baal/service"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	debug     = config.Global.DEBUG
	port      = config.Global.PORT
	https     = config.Global.HTTPS
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run Ball server for localhost",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := config.GlobalConf{
				DEBUG: debug,
				PORT:  port,
			}

			shotdown := make(chan os.Signal, 1)
			signal.Notify(shotdown, syscall.SIGINT, syscall.SIGTERM)
			config.Setup(conf)
			logger.Setup()

			db, err := database.New()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}

			services := service.New(db)
			controllers := controller.New(services)
			router := router.New(controllers)
			srv := router.Serve(port)
			logger.Log.Info(fmt.Sprintf("Server start on >>> %s port", strconv.Itoa(conf.PORT)))

			go srv.ListenAndServe()
			<-shotdown

			logger.Log.Info("Server shotdown")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVarP(&debug, "dev", "", debug, "Use debug/release mode")
	serverCmd.Flags().IntVarP(&port, "port", "p", port, "Server listent on port")
	serverCmd.Flags().BoolVarP(&https, "https", "", https, "Use https protocol for host")
}
