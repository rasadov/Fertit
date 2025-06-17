-- Subscriber table
CREATE TABLE subscribers (
    email VARCHAR(120) UNIQUE NOT NULL,
    uuid UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    policy_updates BOOLEAN NOT NULL DEFAULT TRUE,
    incidents BOOLEAN NOT NULL DEFAULT TRUE,
    new_features BOOLEAN NOT NULL DEFAULT TRUE,
    news BOOLEAN NOT NULL DEFAULT TRUE,
    others BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- user table
CREATE TABLE users (
    username VARCHAR(80) NOT NULL PRIMARY KEY,
    password VARCHAR(120) NOT NULL
);
