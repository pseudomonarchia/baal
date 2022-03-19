package migrationfile

import (
	"baal/database/migrations"
	"baal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations.SetMigration(&gormigrate.Migration{
		ID: "20220319230530_init_user",
		Migrate: func(db *gorm.DB) error {
			model := &models.UserSchema{}
			return db.Migrator().AutoMigrate(model)
		},
		Rollback: func(db *gorm.DB) error {
			model := &models.UserSchema{}
			return db.Migrator().DropTable(model.TableName())
		},
	})
}
