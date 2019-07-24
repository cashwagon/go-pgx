package pgx

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type Config = struct {
	DBName         string `yaml:"db_name"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Host           string `yaml:"host"`
	Port           uint   `yaml:"port"`
	SSLMode        string `yaml:"ssl_mode"`
	ConnectTimeout int    `yaml:"connect_timeout"`
}

func assertSeedsVersion(t *testing.T, db *sqlx.DB, expectedVersion uint) {
	t.Helper()

	assertSchemaVersion(t, db, "schema_seeds", expectedVersion)
}

func assertMigrationsVersion(t *testing.T, db *sqlx.DB, expectedVersion uint) {
	t.Helper()

	assertSchemaVersion(t, db, "schema_migrations", expectedVersion)
}

func assertSchemaVersion(t *testing.T, db *sqlx.DB, table string, expectedVersion uint) {
	t.Helper()

	var rowsCount int
	err := db.Get(&rowsCount, fmt.Sprintf(
		"SELECT count(*) FROM %s", pq.QuoteIdentifier(table),
	))
	require.Nil(t, err)

	if rowsCount == 0 {
		assert.Equal(t, expectedVersion, uint(0))
	} else {
		var version uint
		err = db.Get(&version, fmt.Sprintf(
			"SELECT version FROM %s", pq.QuoteIdentifier(table),
		))
		require.Nil(t, err)
		assert.Equal(t, expectedVersion, version)
	}
}

func assertRowsCount(t *testing.T, db *sqlx.DB, table string, expectedCount int) {
	t.Helper()

	var rowsCount int
	err := db.Get(&rowsCount, fmt.Sprintf(
		"SELECT count(*) FROM %s", pq.QuoteIdentifier(table),
	))
	require.Nil(t, err)

	assert.Equal(t, expectedCount, rowsCount)
}

func assertTableExist(t *testing.T, db *sqlx.DB, tableName string) {
	t.Helper()

	isExist := isTableExist(t, db, tableName)
	assert.Truef(t, isExist, "Table %s should exist", tableName)
}

func assertTableNotExist(t *testing.T, db *sqlx.DB, tableName string) {
	t.Helper()

	isExist := isTableExist(t, db, tableName)
	assert.Falsef(t, isExist, "Table %s should not exist", tableName)
}

func isTableExist(t *testing.T, db *sqlx.DB, tableName string) bool {
	t.Helper()

	var isExist bool
	err := db.Get(&isExist,
		`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = $1
			)
		`,
		tableName,
	)
	require.Nil(t, err)

	return isExist
}

func buildTestOptions(t *testing.T) *Options {
	t.Helper()

	config := loadTestConfig(t)
	return &Options{
		DBName:          config.DBName,
		User:            config.User,
		Password:        config.Password,
		Host:            config.Host,
		Port:            config.Port,
		SSLMode:         config.SSLMode,
		ConnectTimeout:  config.ConnectTimeout,
		ConnMaxLifetime: 0,
		MaxOpenConns:    1,
		MaxIdleConns:    0,
	}
}

func loadTestConfig(t *testing.T) Config {
	t.Helper()

	var config Config

	cfgFile, err := ioutil.ReadFile("testdata/config/config.yaml")
	require.NoError(t, err)

	err = yaml.Unmarshal(cfgFile, &config)
	require.NoError(t, err)

	return config
}
