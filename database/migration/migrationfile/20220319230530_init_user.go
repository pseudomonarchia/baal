package migrationfile

import (
	"baal/database/migration"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type initUserModel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"size:30;not null"`
	Nickname  string `gorm:"size:30;not null"`
	Email     string `gorm:"size:50;not null"`
	Enable    bool   `gorm:"default:false;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (*initUserModel) TableName() string {
	return "user"
}

func init() {
	migration.SetMigration(&gormigrate.Migration{
		ID: "20220319230530_init_user",
		Migrate: func(db *gorm.DB) error {
			return db.Migrator().AutoMigrate(&initUserModel{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable((&initUserModel{}).TableName())
		},
	})
}
