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
	seedRollbackAll                = false
	seedRootCmd     *cobra.Command = &cobra.Command{
		Use:   "seed",
		Short: "Seed cli tools",
		Args:  cobra.MinimumNArgs(1),
	}
	seedCreateCmd *cobra.Command = &cobra.Command{
		Use:   "create",
		Short: "Create migration file",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opt, err := migrations.CreateFile(args[0], migrations.TargetSeeds)
			if err != nil {
				return err
			}

			fmt.Printf("[Baal CLI] Seed create to >>> %s\n", opt)
			return nil
		},
	}
	seedUpCmd *cobra.Command = &cobra.Command{
		Use:   "up",
		Short: "Seed SQL schema sync up",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.Setup()
			databaseName := db.Migrator().CurrentDatabase()
			fmt.Printf("[Baal CLI] Migrate Current Database >>> %s\n", databaseName)

			if err != nil {
				return err
			}

			return migrations.Migrate(db, migrations.TargetSeeds)
		},
	}
	seedDownCmd *cobra.Command = &cobra.Command{
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
				return migrations.RollbackAll(db, migrations.TargetSeeds, func() {})
			}

			return migrations.RollbackLast(db, migrations.TargetSeeds)
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
