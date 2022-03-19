package commands

import (
	"baal/database"
	"baal/database/migrations"
	"fmt"

	"github.com/spf13/cobra"

	//nolint
	_ "baal/database/migrations/migration_file"
)

var (
	rollbackAll                   = false
	migrateRootCmd *cobra.Command = &cobra.Command{
		Use:   "migrate",
		Short: "Migration cli tools",
		Args:  cobra.MinimumNArgs(1),
	}
	migrateCreateCmd *cobra.Command = &cobra.Command{
		Use:   "create",
		Short: "Create migration file",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt, err := migrations.CreateFile(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("[Baal CLI] Migrate create to >>> %s\n", opt)
			return nil
		},
	}
	migrateUpCmd *cobra.Command = &cobra.Command{
		Use:   "up",
		Short: "Migrate SQL schema sync up",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.Setup()
			databaseName := db.Migrator().CurrentDatabase()
			fmt.Printf("[Baal CLI] Migrate Current Database >>> %s\n", databaseName)

			if err != nil {
				return err
			}

			return migrations.Migrate(db)
		},
	}
	migrateDownCmd *cobra.Command = &cobra.Command{
		Use:   "down",
		Short: "Migrate SQL schema sync down",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.Setup()
			databaseName := db.Migrator().CurrentDatabase()
			fmt.Printf("[Baal CLI] Migrate Current Database >>> %s\n", databaseName)
			if err != nil {
				return err
			}

			if rollbackAll {
				return migrations.RollbackAll(db, func() {})
			}

			return migrations.RollbackLast(db)
		},
	}
)

func init() {
	rootCmd.AddCommand(migrateRootCmd)
	migrateRootCmd.AddCommand(
		migrateCreateCmd,
		migrateUpCmd,
		migrateDownCmd,
	)

	migrateDownCmd.Flags().BoolVarP(&rollbackAll, "all", "", false, "Rollback all")
}
