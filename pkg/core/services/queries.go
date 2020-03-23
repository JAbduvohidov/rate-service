package services

const ratesDDL = `CREATE TABLE IF NOT EXISTS rates
(
    id          BIGSERIAL PRIMARY KEY,
    movie_id    INTEGER NOT NULL,
    title       TEXT    NOT NULL,
    description TEXT    NOT NULL,
    image       TEXT    NOT NULL,
    year        TEXT    NOT NULL,
    country     TEXT    NOT NULL,
    actors      TEXT[]  NOT NULL,
    genres      TEXT[]  NOT NULL,
    creators    TEXT[]  NOT NULL,
    studio      TEXT    NOT NULL,
    extLink     TEXT    NOT NULL,
    user_id     INTEGER NOT NULL,
    user_name   TEXT    NOT NULL,
    rate        INTEGER NOT NULL,
    removed     BOOLEAN DEFAULT FALSE
);`
