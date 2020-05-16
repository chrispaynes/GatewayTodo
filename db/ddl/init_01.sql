SET TIME ZONE 'UTC';

DROP SCHEMA public CASCADE;

CREATE SCHEMA app;

-- CREATE ADMIN ROLE
CREATE ROLE admin WITH LOGIN SUPERUSER CREATEDB CREATEROLE REPLICATION NOINHERIT PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE todos to admin;
ALTER ROLE admin SET search_path = app;

-- ALTER READ_WRITE USER ROLE
GRANT ALL PRIVILEGES ON DATABASE todos to read_write_user;
GRANT CONNECT ON DATABASE todos TO read_write_user;
GRANT USAGE ON SCHEMA app TO read_write_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app TO read_write_user;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA app TO read_write_user;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA app TO read_write_user;
ALTER ROLE read_write_user SET search_path = app;
ALTER ROLE read_write_user WITH NOINHERIT;

CREATE TYPE app.todo_status AS ENUM ('new', 'in progress', 'complete');

CREATE TABLE IF NOT EXISTS app.todo (
    id                      SERIAL NOT NULL,
    title                   VARCHAR(80) DEFAULT NULL,
    description             VARCHAR(1024) DEFAULT NULL,
    status                  app.todo_status DEFAULT 'new',
    constraint todo_pkey    PRIMARY KEY (id)
);

-- START THE SERIAL SEQUENCE WITH A LARGER VALUE THAN 1
ALTER SEQUENCE app.todo_id_seq RESTART WITH 100000;
