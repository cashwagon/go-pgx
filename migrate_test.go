package pgx

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

// Positive suite
type MigratePositiveSuite struct {
	suite.Suite
	options        *Options
	migrationsPath string
	db             *sqlx.DB
}

func (s *MigratePositiveSuite) setup() {
	s.options = buildTestOptions(s.T())
	s.migrationsPath = "testdata/migrations"
}

func (s *MigratePositiveSuite) SetupSuite() {
	s.setup()
}

func (s *MigratePositiveSuite) SetupTest() {
	Drop(s.options) // nolint:errcheck

	s.Require().NoError(Create(s.options))

	db, err := Connect(s.options)
	s.Require().NoError(err)
	s.db = db
}

func (s *MigratePositiveSuite) TearDownTest() {
	if s.db != nil {
		s.Require().NoError(s.db.Close())
	}
	s.Require().NoError(Drop(s.options))
}

func (s *MigratePositiveSuite) TestMigrateUp() {
	s.NoError(MigrateUp(s.options, s.migrationsPath))

	assertMigrationsVersion(s.T(), s.db, 2)
	assertTableExist(s.T(), s.db, "samples")
	assertTableExist(s.T(), s.db, "users")
}

func (s *MigratePositiveSuite) TestMigrateDown() {
	s.Require().NoError(MigrateUp(s.options, s.migrationsPath))

	s.NoError(MigrateDown(s.options, s.migrationsPath))

	assertMigrationsVersion(s.T(), s.db, 0)
	assertTableNotExist(s.T(), s.db, "samples")
	assertTableNotExist(s.T(), s.db, "users")
}

func (s *MigratePositiveSuite) TestMigrateToUp() {
	s.NoError(MigrateTo(s.options, s.migrationsPath, 1))

	assertMigrationsVersion(s.T(), s.db, 1)
	assertTableExist(s.T(), s.db, "samples")
	assertTableNotExist(s.T(), s.db, "users")
}

func (s *MigratePositiveSuite) TestMigrateToDown() {
	s.Require().NoError(MigrateUp(s.options, s.migrationsPath))

	s.NoError(MigrateTo(s.options, s.migrationsPath, 1))

	assertMigrationsVersion(s.T(), s.db, 1)
	assertTableExist(s.T(), s.db, "samples")
	assertTableNotExist(s.T(), s.db, "users")
}

// Negative suite with broken migrations
type MigrateBrokenMigrationsSuite struct {
	MigratePositiveSuite
	brokenMigrationsPath string
}

func (s *MigrateBrokenMigrationsSuite) SetupSuite() {
	s.setup()
	s.brokenMigrationsPath = "testdata/migrations_broken"
}

func (s *MigrateBrokenMigrationsSuite) TestMigrateUp() {
	s.Error(MigrateUp(s.options, s.brokenMigrationsPath))

	assertMigrationsVersion(s.T(), s.db, 0)
	assertTableNotExist(s.T(), s.db, "samples")
	assertTableNotExist(s.T(), s.db, "users")
}

func (s *MigrateBrokenMigrationsSuite) TestMigrateDown() {
	s.Require().NoError(MigrateUp(s.options, s.migrationsPath))

	s.Error(MigrateDown(s.options, s.brokenMigrationsPath))

	assertMigrationsVersion(s.T(), s.db, 2)
	assertTableExist(s.T(), s.db, "samples")
	assertTableExist(s.T(), s.db, "users")
}

func (s *MigrateBrokenMigrationsSuite) TestMigrateToUp() {
	s.Error(MigrateTo(s.options, s.brokenMigrationsPath, 1))

	assertMigrationsVersion(s.T(), s.db, 0)
	assertTableNotExist(s.T(), s.db, "samples")
	assertTableNotExist(s.T(), s.db, "users")
}

func (s *MigrateBrokenMigrationsSuite) TestMigrateToDown() {
	s.Require().NoError(MigrateUp(s.options, s.migrationsPath))

	s.Error(MigrateTo(s.options, s.brokenMigrationsPath, 1))

	assertMigrationsVersion(s.T(), s.db, 2)
	assertTableExist(s.T(), s.db, "samples")
	assertTableExist(s.T(), s.db, "users")
}

// Negative suite with no migrations path
type MigrateNoMigrationsSuite struct {
	MigratePositiveSuite
	brokenMigrationsPath string
}

func (s *MigrateNoMigrationsSuite) SetupSuite() {
	s.setup()
	s.brokenMigrationsPath = "testdata/migrations/brokenpath"
}

func (s *MigrateNoMigrationsSuite) TestMigrateUp() {
	s.Error(MigrateUp(s.options, s.brokenMigrationsPath))

	assertTableNotExist(s.T(), s.db, "schema_migrations")
	assertTableNotExist(s.T(), s.db, "samples")
	assertTableNotExist(s.T(), s.db, "users")
}

func (s *MigrateNoMigrationsSuite) TestMigrateDown() {
	s.Require().NoError(MigrateUp(s.options, s.migrationsPath))

	s.Error(MigrateDown(s.options, s.brokenMigrationsPath))

	assertMigrationsVersion(s.T(), s.db, 2)
	assertTableExist(s.T(), s.db, "samples")
	assertTableExist(s.T(), s.db, "users")
}

func (s *MigrateNoMigrationsSuite) TestMigrateToUp() {
	s.Error(MigrateTo(s.options, s.brokenMigrationsPath, 1))

	assertTableNotExist(s.T(), s.db, "schema_migrations")
	assertTableNotExist(s.T(), s.db, "samples")
	assertTableNotExist(s.T(), s.db, "users")
}

func (s *MigrateNoMigrationsSuite) TestMigrateToDown() {
	s.Require().NoError(MigrateUp(s.options, s.migrationsPath))

	s.Error(MigrateTo(s.options, s.brokenMigrationsPath, 1))

	assertMigrationsVersion(s.T(), s.db, 2)
	assertTableExist(s.T(), s.db, "samples")
	assertTableExist(s.T(), s.db, "users")
}

// Run tests
func TestMigratePositiveSuite(t *testing.T) {
	suite.Run(t, new(MigratePositiveSuite))
}

func TestMigrateBrokenMigrationsSuite(t *testing.T) {
	suite.Run(t, new(MigrateBrokenMigrationsSuite))
}

func TestMigrateNoMigrationsSuite(t *testing.T) {
	suite.Run(t, new(MigrateNoMigrationsSuite))
}
