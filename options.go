package pgx

// Options is connection parameters
// Additional info: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
type Options struct {
	DBName         string
	User           string
	Password       string
	Host           string
	Port           uint
	SSLMode        string
	ConnectTimeout int
	SSLCert        string
	SSLKey         string
	SSLRootCert    string

	// Additional parameters for database setup

	// The maximum amount of time a connection may be reused.
	// Additional info: https://golang.org/pkg/database/sql/#DB.SetConnMaxLifetime
	ConnMaxLifetime int
	// The maximum number of connections in the idle connection pool.
	// Additional info: https://golang.org/pkg/database/sql/#DB.SetMaxIdleConns
	MaxOpenConns int
	// The maximum number of open connections to the database.
	// Additional info: https://golang.org/pkg/database/sql/#DB.SetMaxOpenConns
	MaxIdleConns int
}
