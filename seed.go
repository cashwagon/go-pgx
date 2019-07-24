package pgx

func SeedUp(opts *Options, seedsPath string) error {
	dbURL := BuildURL(opts) + "&x-migrations-table=schema_seeds"
	return migrateUp(seedsPath, dbURL)
}

func SeedDown(opts *Options, seedsPath string) error {
	dbURL := BuildURL(opts) + "&x-migrations-table=schema_seeds"
	return migrateDown(seedsPath, dbURL)
}

func SeedTo(opts *Options, seedsPath string, version uint) error {
	dbURL := BuildURL(opts) + "&x-migrations-table=schema_seeds"
	return migrateTo(seedsPath, dbURL, version)
}
