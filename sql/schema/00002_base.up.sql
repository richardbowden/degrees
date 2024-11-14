CREATE TYPE account_type AS ENUM ('org', 'user');

create table if not exists accounts
(
    account_id                      bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    account_number                  varchar not null,

    first_name                      varchar not null,
    middle_name                     varchar,
    surname                         varchar,

    email                           varchar not null,

    acc_type			            account_type not null,

    sign_up_stage                   int default 0,

    password_hash		            varchar not null,

    enabled                         bool        not null         default true,

    created_on                      timestamptz not null         default now(),
    updated_at                      timestamptz not null         default now()
);

-- ALTER table accounts add constraint email_unique UNIQUE (email);
SELECT add_updated_at_trigger('accounts');
CREATE INDEX idx_email ON accounts (email);

create table if not exists email_addresses
(
      email             varchar not null PRIMARY KEY,
      account_id        bigint  not null,
      verified		    bool	not null default false,
      verified_on	    timestamptz,
      updated_at        timestamptz,
      FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

select add_updated_at_trigger('email_addresses');

create table if not exists email_verification_code (

    email_verification_code_id      bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    code                            varchar not null,
    account_id  bigint not null,
    expires_at  timestamptz,

    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

create table if not exists addresses (

    address_id      bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    line1           varchar not null,
    line2           varchar,
    postcode        varchar not null,
    state           varchar,
    country         varchar not null
);

create table if not exists coffee_shops 
(

    coffee_shop_id  bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name            varchar not null,
    description     varchar,
    owner_id        bigint not null, 
    enabled         bool not null default false,

    FOREIGN KEY (owner_id) REFERENCES accounts(account_id)

);

create table if not exists coffee_shop_locations (

    coffee_shop_location_id              bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    
    coffee_shop_id  bigint not null,
    address_id      bigint not null,
    
    FOREIGN KEY (coffee_shop_id) REFERENCES coffee_shops(coffee_shop_id),
    FOREIGN KEY (address_id) REFERENCES addresses(address_id)

);

create table if not exists coffee_clubs
(
    coffee_club_id  bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name            varchar not null,
    owner_id        bigint  not null,
    public          bool not null default false,

    FOREIGN KEY (owner_id) REFERENCES accounts(account_id)

);

create table if not exists coffee_club_members
(
    coffee_club_member_id   bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    coffee_club_id          bigint not null,
    account_id              bigint not null,
    

    FOREIGN KEY (coffee_club_id) REFERENCES coffee_clubs(coffee_club_id),
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);
