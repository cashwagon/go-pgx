package pgx

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // library demands such import
	_ "github.com/golang-migrate/migrate/v4/source/file"       // library demands such import
	"github.com/pkg/errors"
)

// MigrateUp runs migrations on the given database
func MigrateUp(opts *Options, migrationsPath string) error {
	return migrateUp(migrationsPath, BuildURL(opts))
}

// MigrateDown rollbacks migrations on the given database
func MigrateDown(opts *Options, migrationsPath string) error {
	return migrateDown(migrationsPath, BuildURL(opts))
}

// MigrateTo runs migrations up to the given version on the given database
func MigrateTo(opts *Options, migrationsPath string, version uint) error {
	return migrateTo(migrationsPath, BuildURL(opts), version)
}

func migrateUp(migrationsPath, dbURL string) error {
	return wrapMigration(migrationsPath, dbURL, func(m *migrate.Migrate) error {
		return m.Up()
	})
}

func migrateDown(migrationsPath, dbURL string) error {
	return wrapMigration(migrationsPath, dbURL, func(m *migrate.Migrate) error {
		return m.Down()
	})
}

func migrateTo(migrationsPath, dbURL string, version uint) error {
	return wrapMigration(migrationsPath, dbURL, func(m *migrate.Migrate) error {
		return m.Migrate(version)
	})
}

func wrapMigration(migrationsPath, dbURL string, fn func(*migrate.Migrate) error) error {
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		return errors.Wrap(err, "could not open migration")
	}
	defer m.Close()

	version, err := getSchemaVersion(m)
	if err != nil {
		return errors.Wrap(err, "could not get schema version")
	}

	if err := fn(m); err != nil {
		if e := recoverSchema(m, version); e != nil {
			return errors.Wrap(err, "could not reciver schema")
		}
		return errors.Wrap(err, "could not migrate")
	}
	return nil
}

func recoverSchema(m *migrate.Migrate, prevVersion uint) error {
	brokenVersion, err := getSchemaVersion(m)
	if err != nil {
		return errors.Wrap(err, "could not get schema version")
	}

	actualVersion := brokenVersion
	if brokenVersion < prevVersion {
		actualVersion = brokenVersion + 1
	}
	if brokenVersion > prevVersion {
		actualVersion = brokenVersion - 1
	}

	err = m.Force(int(actualVersion))
	if err != nil {
		return errors.Wrap(err, "could not force migration version")
	}

	return nil
}

func getSchemaVersion(m *migrate.Migrate) (version uint, err error) {
	version, _, err = m.Version()
	if err != nil {
		if err.Error() == "no migration" {
			err = nil
		} else {
			err = errors.Wrap(err, "could get schema version")
		}
	}

	return
}
