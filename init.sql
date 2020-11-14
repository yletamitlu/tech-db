DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE IF NOT EXISTS users
(
    nickname varchar(64) NOT NULL PRIMARY KEY,
    fullname varchar(64) NOT NULL CHECK (fullname <> ''),
    email    varchar(64) NOT NULL UNIQUE CHECK (email <> ''),
    about    text
);

CREATE TABLE IF NOT EXISTS forums
(
    author_nickname varchar(64) NOT NULL REFERENCES users (nickname) ON DELETE CASCADE,
    title           varchar(64) NOT NULL CHECK (title <> ''),
    slug            varchar(64) NOT NULL CHECK (slug <> '') PRIMARY KEY,
    threads         integer DEFAULT 0,
    posts           integer DEFAULT 0
);

CREATE TABLE IF NOT EXISTS threads
(
    id              serial      NOT NULL PRIMARY KEY,
    slug            varchar(64) NOT NULL CHECK (slug <> ''),
    author_nickname varchar(64) REFERENCES users (nickname) ON DELETE CASCADE,
    created_at      timestamptz DEFAULT now(),
    forum_slug      varchar(64) REFERENCES forums (slug) ON DELETE CASCADE,
    message         text,
    title           varchar(64) NOT NULL,
    votes           integer
);

CREATE TABLE IF NOT EXISTS votes
(
    user_nickname varchar(64) REFERENCES users (nickname) ON DELETE CASCADE,
    voice         integer NOT NULL,
    thread_id     integer REFERENCES threads (id) ON DELETE CASCADE,
    PRIMARY KEY (user_nickname, thread_id)
);

CREATE TABLE IF NOT EXISTS posts
(
    id              serial  NOT NULL PRIMARY KEY,
    created_at      timestamptz      DEFAULT now(),
    message         text,
    author_nickname varchar(64) REFERENCES users (nickname) ON DELETE CASCADE,
    forum_slug      varchar(64) REFERENCES forums (slug) ON DELETE CASCADE,
    thread_id       integer NOT NULL DEFAULT 0 REFERENCES threads (id) ON DELETE CASCADE,
    parent          integer REFERENCES posts (id) ON DELETE CASCADE,
    is_edited       boolean NOT NULL DEFAULT false
)
