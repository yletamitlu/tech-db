DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE IF NOT EXISTS users
(
    nickname citext NOT NULL PRIMARY KEY,
    fullname varchar(64) NOT NULL CHECK (fullname <> ''),
    email    citext NOT NULL UNIQUE CHECK (email <> ''),
    about    text
);

DROP TABLE IF EXISTS forums CASCADE;
CREATE TABLE IF NOT EXISTS forums
(
    author_nickname citext NOT NULL REFERENCES users (nickname) ON DELETE CASCADE,
    title           varchar(64) NOT NULL CHECK (title <> ''),
    slug            citext NOT NULL CHECK (slug <> '') PRIMARY KEY,
    threads         integer DEFAULT 0,
    posts           integer DEFAULT 0
);

DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE IF NOT EXISTS threads
(
    id              serial      NOT NULL PRIMARY KEY,
    slug            citext NOT NULL CHECK (slug <> ''),
    author_nickname citext REFERENCES users (nickname) ON DELETE CASCADE,
    created_at      timestamptz DEFAULT now(),
    forum_slug      citext REFERENCES forums (slug) ON DELETE CASCADE,
    message         text,
    title           varchar(64) NOT NULL,
    votes           integer
);

DROP TABLE IF EXISTS votes CASCADE;
CREATE TABLE IF NOT EXISTS votes
(
    user_nickname citext REFERENCES users (nickname) ON DELETE CASCADE,
    voice         integer NOT NULL,
    thread_id     integer REFERENCES threads (id) ON DELETE CASCADE,
    PRIMARY KEY (user_nickname, thread_id)
);

DROP TABLE IF EXISTS posts CASCADE;
CREATE TABLE IF NOT EXISTS posts
(
    id              serial  NOT NULL PRIMARY KEY,
    created_at      timestamptz      DEFAULT now(),
    message         text,
    author_nickname citext REFERENCES users (nickname) ON DELETE CASCADE,
    forum_slug      citext REFERENCES forums (slug) ON DELETE CASCADE,
    thread_id       integer NOT NULL DEFAULT 0 REFERENCES threads (id) ON DELETE CASCADE,
    parent          integer REFERENCES posts (id) ON DELETE CASCADE,
    is_edited       boolean NOT NULL DEFAULT false
)
