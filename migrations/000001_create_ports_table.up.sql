CREATE TABLE IF NOT EXISTS ports(
    id serial PRIMARY KEY,
    locode CHAR(5) UNIQUE NOT NULL,
    name text NOT NULL,
    city text NOT NULL,
    country text NOT NULL,
    alias text[],
    regions text[],
    lat text NOT NULL,
    lon text NOT NULL,
    province    text NOT NULL,
    timezone    text NOT NULL,
    unlocs text[] NOT NULL,
    foreign_code integer NOT NULL
);