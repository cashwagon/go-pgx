package main

import (
	"io/ioutil"
	"log"

	"github.com/cashwagon/go-pgx"
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
	SSLCert        string `yaml:"ssl_cert"`
	SSLKey         string `yaml:"ssl_key"`
	SSLRootCert    string `yaml:"ssl_root_cert"`
}

func main() {
	options, err := loadOptions("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating database...")
	err = pgx.Create(options)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migrating database...")
	err = pgx.MigrateUp(options, "db/migrations")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Seeding database...")
	err = pgx.SeedUp(options, "db/seeds")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to database...")
	db, err := pgx.Connect(options)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching users count...")
	var usersCount int
	err = db.Get(&usersCount, "SELECT count(*) FROM users")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Users count: %d\n", usersCount)

	log.Println("Closing database...")
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Rollback seeds...")
	err = pgx.SeedDown(options, "db/seeds")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Rollback migrations...")
	err = pgx.MigrateDown(options, "db/migrations")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Dropping database...")
	err = pgx.Drop(options)
	if err != nil {
		log.Fatal(err)
	}
}

func loadOptions(filePath string) (*pgx.Options, error) {
	cfgFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(cfgFile, &config)
	if err != nil {
		return nil, err
	}

	return &pgx.Options{
		DBName:          config.DBName,
		User:            config.User,
		Password:        config.Password,
		Host:            config.Host,
		Port:            config.Port,
		SSLMode:         config.SSLMode,
		ConnectTimeout:  config.ConnectTimeout,
		SSLCert:         config.SSLCert,
		SSLKey:          config.SSLKey,
		SSLRootCert:     config.SSLRootCert,
		ConnMaxLifetime: 0,
		MaxOpenConns:    1,
		MaxIdleConns:    0,
	}, nil
}
