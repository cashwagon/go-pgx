package pgx

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

// Positive suite
type SeedPositiveSuite struct {
	suite.Suite
	options        *Options
	migrationsPath string
	seedsPath      string
	db             *sqlx.DB
}

func (s *SeedPositiveSuite) setup() {
	s.options = buildTestOptions(s.T())
	s.migrationsPath = "testdata/migrations"
	s.seedsPath = "testdata/seeds"
}

func (s *SeedPositiveSuite) SetupSuite() {
	s.setup()
}

func (s *SeedPositiveSuite) SetupTest() {
	Drop(s.options) // nolint:errcheck

	s.Require().NoError(Create(s.options))
	s.Require().NoError(MigrateUp(s.options, s.migrationsPath))

	db, err := Connect(s.options)
	s.Require().NoError(err)
	s.db = db
}

func (s *SeedPositiveSuite) TearDownTest() {
	if s.db != nil {
		s.Require().NoError(s.db.Close())
	}
	s.Require().NoError(Drop(s.options))
}

func (s *SeedPositiveSuite) TestSeedUp() {
	s.NoError(SeedUp(s.options, s.seedsPath))

	assertSeedsVersion(s.T(), s.db, 2)
	assertRowsCount(s.T(), s.db, "samples", 1)
	assertRowsCount(s.T(), s.db, "users", 1)
}

func (s *SeedPositiveSuite) TestSeedDown() {
	s.Require().NoError(SeedUp(s.options, s.seedsPath))

	s.NoError(SeedDown(s.options, s.seedsPath))

	assertSeedsVersion(s.T(), s.db, 0)
	assertRowsCount(s.T(), s.db, "samples", 0)
	assertRowsCount(s.T(), s.db, "users", 0)
}

func (s *SeedPositiveSuite) TestSeedToUp() {
	s.NoError(SeedTo(s.options, s.seedsPath, 1))

	assertSeedsVersion(s.T(), s.db, 1)
	assertRowsCount(s.T(), s.db, "samples", 1)
	assertRowsCount(s.T(), s.db, "users", 0)
}

func (s *SeedPositiveSuite) TestSeedToDown() {
	s.Require().NoError(SeedUp(s.options, s.seedsPath))

	s.NoError(SeedTo(s.options, s.seedsPath, 1))

	assertSeedsVersion(s.T(), s.db, 1)
	assertRowsCount(s.T(), s.db, "samples", 1)
	assertRowsCount(s.T(), s.db, "users", 0)
}

// Negative suite with broken migrations
type SeedBrokenSeedsSuite struct {
	SeedPositiveSuite
	brokenSeedsPath string
}

func (s *SeedBrokenSeedsSuite) SetupSuite() {
	s.setup()
	s.brokenSeedsPath = "testdata/seeds_broken"
}

func (s *SeedBrokenSeedsSuite) TestSeedUp() {
	s.Error(SeedUp(s.options, s.brokenSeedsPath))

	assertSeedsVersion(s.T(), s.db, 0)
	assertRowsCount(s.T(), s.db, "samples", 0)
	assertRowsCount(s.T(), s.db, "users", 0)
}

func (s *SeedBrokenSeedsSuite) TestSeedDown() {
	s.Require().NoError(SeedUp(s.options, s.seedsPath))

	s.Error(SeedDown(s.options, s.brokenSeedsPath))

	assertSeedsVersion(s.T(), s.db, 2)
	assertRowsCount(s.T(), s.db, "samples", 1)
	assertRowsCount(s.T(), s.db, "users", 1)
}

func (s *SeedBrokenSeedsSuite) TestSeedToUp() {
	s.Error(SeedTo(s.options, s.brokenSeedsPath, 1))

	assertSeedsVersion(s.T(), s.db, 0)
	assertRowsCount(s.T(), s.db, "samples", 0)
	assertRowsCount(s.T(), s.db, "users", 0)
}

func (s *SeedBrokenSeedsSuite) TestSeedToDown() {
	s.Require().NoError(SeedUp(s.options, s.seedsPath))

	s.Error(SeedTo(s.options, s.brokenSeedsPath, 1))

	assertSeedsVersion(s.T(), s.db, 2)
	assertRowsCount(s.T(), s.db, "samples", 1)
	assertRowsCount(s.T(), s.db, "users", 1)
}

// Negative suite with no migrations path
type SeedNoSeedsSuite struct {
	SeedPositiveSuite
	brokenSeedsPath string
}

func (s *SeedNoSeedsSuite) SetupSuite() {
	s.setup()
	s.brokenSeedsPath = "testdata/seeds/brokenpath"
}

func (s *SeedNoSeedsSuite) TestSeedUp() {
	s.Error(SeedUp(s.options, s.brokenSeedsPath))

	assertTableNotExist(s.T(), s.db, "schema_seeds")
	assertRowsCount(s.T(), s.db, "samples", 0)
	assertRowsCount(s.T(), s.db, "users", 0)
}

func (s *SeedNoSeedsSuite) TestSeedDown() {
	s.Require().NoError(SeedUp(s.options, s.seedsPath))

	s.Error(SeedDown(s.options, s.brokenSeedsPath))

	assertSeedsVersion(s.T(), s.db, 2)
	assertRowsCount(s.T(), s.db, "samples", 1)
	assertRowsCount(s.T(), s.db, "users", 1)
}

func (s *SeedNoSeedsSuite) TestSeedToUp() {
	s.Error(SeedTo(s.options, s.brokenSeedsPath, 1))

	assertTableNotExist(s.T(), s.db, "schema_seeds")
	assertRowsCount(s.T(), s.db, "samples", 0)
	assertRowsCount(s.T(), s.db, "users", 0)
}

func (s *SeedNoSeedsSuite) TestSeedToDown() {
	s.Require().NoError(SeedUp(s.options, s.seedsPath))

	s.Error(SeedTo(s.options, s.brokenSeedsPath, 1))

	assertSeedsVersion(s.T(), s.db, 2)
	assertRowsCount(s.T(), s.db, "samples", 1)
	assertRowsCount(s.T(), s.db, "users", 1)
}

// Run tests
func TestSeedPositiveSuite(t *testing.T) {
	suite.Run(t, new(SeedPositiveSuite))
}

func TestSeedBrokenSeedsSuite(t *testing.T) {
	suite.Run(t, new(SeedBrokenSeedsSuite))
}

func TestSeedNoSeedsSuite(t *testing.T) {
	suite.Run(t, new(SeedNoSeedsSuite))
}
