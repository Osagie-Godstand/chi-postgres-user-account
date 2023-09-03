-- Create the users table
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
-- Create the "users" table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    encrypted_password VARCHAR(60) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT false,
    UNIQUE (email)
);
  END IF;
END $$;
