-- +goose Up
CREATE TABLE posts(
    id uuid PRIMARY KEY ,
    created_at timestamp NOT NULL ,
    updated_at timestamp NOT NULL ,
    title text ,
    url text UNIQUE NOT NULL ,
    description text ,
    published_at timestamp ,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);