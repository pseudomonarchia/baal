package cmd

import (
	"baal/database"
	"baal/database/migration"
	"fmt"

	"github.com/spf13/cobra"

	//nolint
	_ "baal/database/migration/migrationfile"
)

var (
	migrateRollbackAll = false
	migrateRootCmd     = &cobra.Command{
		Use:   "migrate",
		Short: "Migration cli tools",
		Args:  cobra.MinimumNArgs(1),
	}
	migrateCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create migration file",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt, err := migration.CreateFile(args[0], migration.TargetMigrations)
			if err != nil {
				return err
			}

			fmt.Printf("[Baal CLI] Migrate create to >>> %s\n", opt)
			return nil
		},
	}
	migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Migrate SQL schema sync up",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New()
			databaseName := db.Migrator().CurrentDatabase()
			fmt.Printf("[Baal CLI] Migrate Current Database >>> %s\n", databaseName)

			if err != nil {
				return err
			}

			return migration.Migrate(db, migration.TargetMigrations)
		},
	}
	migrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Migrate SQL schema sync down",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New()
			databaseName := db.Migrator().CurrentDatabase()
			fmt.Printf("[Baal CLI] Migrate Current Database >>> %s\n", databaseName)
			if err != nil {
				return err
			}

			if migrateRollbackAll {
				return migration.RollbackAll(db, migration.TargetMigrations, func() {})
			}

			return migration.RollbackLast(db, migration.TargetMigrations)
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

	migrateDownCmd.Flags().BoolVarP(&migrateRollbackAll, "all", "", false, "Rollback all")
}
