DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'sessions') THEN
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    userid UUID REFERENCES users(id) NOT NULL,
    token TEXT NOT NULL,
    expiresat TIMESTAMPTZ NOT NULL
);
 END IF;
END $$;

