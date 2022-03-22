package cmd

import (
	"baal/database"
	"baal/database/migration"
	"fmt"

	"github.com/spf13/cobra"

	//nolint
	_ "baal/database/migration/seedfile"
)

var (
	seedRollbackAll = false
	seedRootCmd     = &cobra.Command{
		Use:   "seed",
		Short: "Seed cli tools",
		Args:  cobra.MinimumNArgs(1),
	}
	seedCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create migration file",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt, err := migration.CreateFile(args[0], migration.TargetSeeds)
			if err != nil {
				return err
			}

			fmt.Printf("[Baal CLI] Seed create to >>> %s\n", opt)
			return nil
		},
	}
	seedUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Seed SQL schema sync up",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.Setup()
			databaseName := db.Migrator().CurrentDatabase()
			fmt.Printf("[Baal CLI] Migrate Current Database >>> %s\n", databaseName)

			if err != nil {
				return err
			}

			return migration.Migrate(db, migration.TargetSeeds)
		},
	}
	seedDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Seed SQL schema sync down",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.Setup()
			databaseName := db.Migrator().CurrentDatabase()
			fmt.Printf("[Baal CLI] Migrate Current Database >>> %s\n", databaseName)
			if err != nil {
				return err
			}

			if seedRollbackAll {
				return migration.RollbackAll(db, migration.TargetSeeds, func() {})
			}

			return migration.RollbackLast(db, migration.TargetSeeds)
		},
	}
)

func init() {
	rootCmd.AddCommand(seedRootCmd)
	seedRootCmd.AddCommand(
		seedCreateCmd,
		seedUpCmd,
		seedDownCmd,
	)

	seedDownCmd.Flags().BoolVarP(&seedRollbackAll, "all", "", false, "Rollback all")
}
