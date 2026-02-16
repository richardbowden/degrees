create table if not exists settings (
    subsystem text not null,
    key text not null,
    value text not null,
    PRIMARY  KEY (subsystem, key)
);
