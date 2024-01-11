DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') THEN
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

-- The sequence and trigger work together to generate and 
-- assign UUIDs automatically when needed, providing a convenient 
-- and consistent way to handle UUID generation for the id column.

-- Create a sequence for the UUID if it doesn't exist
CREATE SEQUENCE IF NOT EXISTS posts_id_seq;

-- Set the default value for the ID column using the sequence
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();

-- Create a trigger to set the UUID on insert if not provided
CREATE OR REPLACE FUNCTION set_post_id()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.id IS NULL THEN
        NEW.id := uuid_generate_v4();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_post_id_trigger
    BEFORE INSERT ON users
    FOR EACH ROW EXECUTE FUNCTION set_post_id();

