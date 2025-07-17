
create table if not exists users
(
    id                              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name                      varchar not null,
    middle_name                     varchar,
    surname                         varchar,
    username                        varchar NOT NULL,
    login_email                     varchar not null,
    primary_email_id                BIGINT not null,
    sign_up_stage                   int default 0,
    password_hash                   varchar not null,
    enabled                         bool        not null         default true,
    sysop                           bool not null default false,
    created_on                      timestamptz not null         default now(),
    updated_at                      timestamptz not null         default now()
);

ALTER table users add constraint users_login_email_unique UNIQUE (login_email);
CREATE INDEX if not exists users_idx_login_email ON users (login_email);
CREATE INDEX if not exists users_idx_username ON users (username);
SELECT add_updated_at_trigger('users');

CREATE TABLE IF NOT EXISTS user_emails (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR NOT NULL,
    is_verified BOOL NOT NULL DEFAULT false,
    enabled BOOL NOT NULL DEFAULT true,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT unique_email UNIQUE (email),
    CONSTRAINT unique_user_email UNIQUE (user_id, email)
);

CREATE INDEX if not exists user_emails_idx_email ON user_emails (email);
SELECT add_updated_at_trigger('user_emails');

CREATE TABLE if not exists profiles (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name VARCHAR(255) NULL,
    avatar VARCHAR(500) NULL, -- URL to avatar image
    bio TEXT NULL,
    website VARCHAR(500) NULL,
    location VARCHAR(255) NULL,
    company VARCHAR(255) NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

SELECT add_updated_at_trigger('profiles');
