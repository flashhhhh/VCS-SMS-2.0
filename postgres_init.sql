CREATE DATABASE user_service_db;

\c user_service_db;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(255) NOT NULL
);

INSERT INTO users (id, username, password, name, email, role) VALUES
    (gen_random_uuid(), 'admin', '$2a$10$6smDl/of0VnSPLw.1qFUgurMGLBaEg.FTLvuXtCTlv8fMQT1dVC2C', 'Admin', 'admin@gmail.com', 'admin');

CREATE DATABASE server_administration_db;

\c server_administration_db;

CREATE TABLE IF NOT EXISTS servers (
    id SERIAL,
    server_id VARCHAR(255) PRIMARY KEY,
    server_name VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ipv4 VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL
);
