package migrations

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	//nolint
	_ "baal/database"
)

var (
	dir, _                                  = os.Getwd()
	migrations      []*gormigrate.Migration = []*gormigrate.Migration{}
	option          *gormigrate.Options     = &gormigrate.Options{}
	timestampFormat string                  = "20060102150405"
)

// CreateFile generate `name.go` file and init to migration
func CreateFile(filename string) (output string, err error) {
	templateBuf, err := ioutil.ReadFile(path.Join(dir, "/database/migrations/template.txt"))
	if err != nil {
		return "", err
	}

	id := time.Now().Format(timestampFormat)
	filename = fmt.Sprintf("%s_%s", id, filename)
	output = path.Join(dir, "/database/migrations/migration_file", fmt.Sprintf("%s.go", filename))
	code := fmt.Sprintf(string(templateBuf), filename)

	f, err := os.Create(output)
	if err != nil {
		return "", err
	}

	defer f.Close()
	f.WriteString(code)

	return output, nil
}

// SetMigration to migrations list
func SetMigration(m *gormigrate.Migration) {
	migrations = append(migrations, m)
}

// Migrate sync all migration to database
func Migrate(db *gorm.DB) error {
	migrate := gormigrate.New(db, option, migrations)
	if err := migrate.Migrate(); err != nil {
		return err
	}

	return nil
}

// RollbackAll rollback all migration to database
func RollbackAll(db *gorm.DB, cb func()) error {
	migrate := gormigrate.New(db, option, migrations)
	list := make([]*gormigrate.Migration, len(migrations))

	copy(list, migrations)
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID > list[j].ID
	})

	for _, migration := range list {
		err := migrate.RollbackMigration(migration)
		if err != nil {
			return err
		}

		cb()
	}

	return nil
}

// RollbackLast rollback last migration to database
func RollbackLast(db *gorm.DB) error {
	migrate := gormigrate.New(db, option, migrations)
	return migrate.RollbackLast()
}
