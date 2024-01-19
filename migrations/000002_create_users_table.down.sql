-- Drop the trigger if it exists
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_user_id_trigger') THEN
        DROP TRIGGER IF EXISTS set_user_id_trigger ON users;
    END IF;
END $$;

-- Drop the function
DROP FUNCTION IF EXISTS set_user_id();

-- Remove the default value for the ID column
ALTER TABLE users ALTER COLUMN id DROP DEFAULT;

-- Drop the sequence
DROP SEQUENCE IF EXISTS users_id_seq;

-- Drop the table
DROP TABLE IF EXISTS users;