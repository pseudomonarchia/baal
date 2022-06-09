package migrationfile

import (
	"baal/database/migration"
	"baal/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type initOAuthTokenModel struct {
	UID       string              `gorm:"primaryKey"`
	UserID    uint                `gorm:"not null"`
	Provider  model.OAuthProvider `gorm:"size:10"`
	Use       bool                `gorm:"default:false"`
	TokenInfo datatypes.JSON      `gorm:"not null"`
}

func (*initOAuthTokenModel) TableName() string {
	return "oauth_token"
}

func init() {
	migration.SetMigration(&gormigrate.Migration{
		ID: "20220326195125_init_oauth_token",
		Migrate: func(db *gorm.DB) error {
			return db.Migrator().AutoMigrate(&initOAuthTokenModel{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable((&initOAuthTokenModel{}).TableName())
		},
	})
}
