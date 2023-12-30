DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
-- Create the "users" table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    encryptedpassword VARCHAR(60) NOT NULL,
    isadmin BOOLEAN NOT NULL DEFAULT false,
    UNIQUE (email)
);
  END IF;
END $$;
