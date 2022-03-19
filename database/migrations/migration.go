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

// TargetMigrations/TargetSeeds is CreateFile  for `migrations/seeds`
var (
	dir, _                                   = os.Getwd()
	TargetMigrations string                  = "migrations"
	TargetSeeds      string                  = "seeds"
	timestampFormat  string                  = "20060102150405"
	migrations       []*gormigrate.Migration = []*gormigrate.Migration{}
)

func getOption(target string) *gormigrate.Options {
	var option *gormigrate.Options

	if target == TargetSeeds {
		option = &gormigrate.Options{TableName: TargetSeeds}
	} else {
		option = &gormigrate.Options{TableName: TargetMigrations}
	}

	return option
}

func getTableName(target string) string {
	var name string

	if target == TargetSeeds {
		name = "seed_file"
	} else {
		name = "migration_file"
	}

	return name
}

// CreateFile generate `name.go` file and init to migration
func CreateFile(filename string, target string) (output string, err error) {
	templateBuf, err := ioutil.ReadFile(path.Join(dir, "/database/migrations/template.txt"))
	if err != nil {
		return "", err
	}

	id := time.Now().Format(timestampFormat)
	filename = fmt.Sprintf("%s_%s", id, filename)
	relativePath := fmt.Sprintf("/database/migrations/%s", getTableName(target))
	output = path.Join(dir, relativePath, fmt.Sprintf("%s.go", filename))
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
func Migrate(db *gorm.DB, target string) error {
	option := getOption(target)
	migrate := gormigrate.New(db, option, migrations)
	if err := migrate.Migrate(); err != nil {
		return err
	}

	return nil
}

// RollbackAll rollback all migration to database
func RollbackAll(db *gorm.DB, target string, cb func()) error {
	option := getOption(target)
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
func RollbackLast(db *gorm.DB, target string) error {
	option := getOption(target)
	migrate := gormigrate.New(db, option, migrations)
	return migrate.RollbackLast()
}
