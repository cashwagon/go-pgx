BEGIN;
  CREATE TABLE IF NOT EXISTS samples (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NUL // Break migration for purpose
  );
COMMIT;