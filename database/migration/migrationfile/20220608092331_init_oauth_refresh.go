package migrationfile

import (
	"baal/database/migration"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type initOAuthRefreshModel struct {
	OAuthUID  string    `gorm:"primaryKey;not null;column:oauth_uid"`
	IP        string    `gorm:"size:30;not null"`
	Token     string    `gorm:"unique"`
	IssuedAt  time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (*initOAuthRefreshModel) TableName() string {
	return "oauth_refresh"
}

func init() {
	migration.SetMigration(&gormigrate.Migration{
		ID: "20220608092331_init_oauth_refresh",
		Migrate: func(db *gorm.DB) error {
			return db.Migrator().AutoMigrate(&initOAuthRefreshModel{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable((&initOAuthRefreshModel{}).TableName())
		},
	})
}
