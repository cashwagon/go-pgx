package pgx

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Libarary demands such import
	"github.com/pkg/errors"
)

// Create creates database by options
func Create(opts *Options) error {
	return rootExec(*opts, fmt.Sprintf("CREATE DATABASE %s", opts.DBName))
}

// Drop drops database by options
func Drop(opts *Options) error {
	return rootExec(*opts, fmt.Sprintf("DROP DATABASE %s", opts.DBName))
}

// Connect creates database connection and returns sqlx.DB pool
func Connect(opts *Options) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", BuildURL(opts))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(opts.ConnMaxLifetime) * time.Second)
	db.SetMaxIdleConns(opts.MaxIdleConns)
	db.SetMaxOpenConns(opts.MaxOpenConns)

	return db, nil
}

// BuildURL build database connection URL
func BuildURL(opts *Options) string {
	dbURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(opts.User, opts.Password),
		Host:   fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Path:   opts.DBName,
	}

	query := dbURL.Query()
	query.Set("connect_timeout", strconv.Itoa(opts.ConnectTimeout))
	query.Set("sslmode", opts.SSLMode)
	query.Set("sslcert", opts.SSLCert)
	query.Set("sslkey", opts.SSLKey)
	query.Set("sslrootcert", opts.SSLRootCert)
	dbURL.RawQuery = query.Encode()

	return dbURL.String()
}

// rootExec opens connection without database and executes one query
func rootExec(opts Options, query string) (err error) { // nolint:gocritic
	opts.DBName = ""
	opts.MaxIdleConns = 0
	opts.MaxOpenConns = 1

	db, err := Connect(&opts)
	if err != nil {
		return errors.Wrap(err, "could not connect to database")
	}
	defer func() {
		e := db.Close()
		if e != nil {
			err = errors.Wrap(e, "could not close database")
		}
	}()

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "could not execute query - '%s'", query)
	}

	return nil
}
