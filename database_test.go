package pgx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Positive suite
type DatabasePositiveSuite struct {
	suite.Suite
	options *Options
}

func (s *DatabasePositiveSuite) SetupSuite() {
	s.options = buildTestOptions(s.T())
}

func (s *DatabasePositiveSuite) SetupTest() {
	Drop(s.options) // nolint:errcheck
}

func (s *DatabasePositiveSuite) TestCreate() {
	s.NoError(Create(s.options))
	defer func() {
		s.Require().NoError(Drop(s.options))
	}()

	db, err := Connect(s.options)
	s.NoError(err)
	s.NoError(db.Close())
}

func (s *DatabasePositiveSuite) TestDrop() {
	s.Require().NoError(Create(s.options))

	s.NoError(Drop(s.options))

	_, err := Connect(s.options)
	s.Error(err)
}

func (s *DatabasePositiveSuite) TestConnect() {
	s.Require().NoError(Create(s.options))
	defer func() {
		s.Require().NoError(Drop(s.options))
	}()

	db, err := Connect(s.options)
	s.NoError(err)
	s.NoError(db.Ping())
	s.NoError(db.Close())
}

// Negative suite
type DatabaseNegativeSuite struct {
	suite.Suite
	options *Options
}

func (s *DatabaseNegativeSuite) SetupSuite() {
	s.options = buildTestOptions(s.T())
	s.options.Port = 5435
}

func (s *DatabaseNegativeSuite) TestCreate() {
	s.Error(Create(s.options))
}

func (s *DatabaseNegativeSuite) TestDrop() {
	s.Error(Drop(s.options))
}

func (s *DatabaseNegativeSuite) TestConnect() {
	_, err := Connect(s.options)
	s.Error(err)
}

// Run tests
func TestBuildURL(t *testing.T) {
	opts := &Options{
		DBName:         "mydb",
		User:           "admin",
		Password:       "qwerty",
		Host:           "127.0.0.1",
		Port:           5435,
		SSLMode:        "prefer",
		ConnectTimeout: 5,
		SSLCert:        "./pgssl.cert",
		SSLKey:         "./pgssl.key",
		SSLRootCert:    "./pgsslroot.cert",
	}
	url := "postgres://admin:qwerty@127.0.0.1:5435/mydb?connect_timeout=5&" +
		"sslcert=.%2Fpgssl.cert&sslkey=.%2Fpgssl.key&sslmode=prefer&" +
		"sslrootcert=.%2Fpgsslroot.cert"
	assert.Equal(t, url, BuildURL(opts), "Must return valid url")
}

func TestDatabasePositiveSuite(t *testing.T) {
	suite.Run(t, new(DatabasePositiveSuite))
}

func TestDatabaseNegativeSuite(t *testing.T) {
	suite.Run(t, new(DatabaseNegativeSuite))
}
