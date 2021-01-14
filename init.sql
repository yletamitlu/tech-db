DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE IF NOT EXISTS users
(
    nickname citext NOT NULL PRIMARY KEY,
    fullname text   NOT NULL,
    email    citext NOT NULL UNIQUE,
    about    text
);

DROP TABLE IF EXISTS forums CASCADE;
CREATE TABLE IF NOT EXISTS forums
(
    author_nickname citext NOT NULL REFERENCES users (nickname) ON DELETE CASCADE,
    title           text   NOT NULL,
    slug            citext NOT NULL PRIMARY KEY,
    threads         integer DEFAULT 0,
    posts           integer DEFAULT 0
);

DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE IF NOT EXISTS threads
(
    id              serial  NOT NULL PRIMARY KEY,
    slug            citext  NOT NULL DEFAULT '',
    author_nickname citext REFERENCES users (nickname) ON DELETE CASCADE,
    created_at      timestamptz      DEFAULT now(),
    forum_slug      citext REFERENCES forums (slug) ON DELETE CASCADE,
    message         text,
    title           text    NOT NULL,
    votes           integer NOT NULL DEFAULT 0
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
    parent          integer NOT NULL DEFAULT 0,
    is_edited       boolean NOT NULL DEFAULT false,
    path            text    NOT NULL DEFAULT ''
);

DROP TABLE IF EXISTS user_forum CASCADE;
CREATE TABLE user_forum
(
    user_nickname CITEXT NOT NULL,
    forum_slug    CITEXT NOT NULL,
    PRIMARY KEY (forum_slug, user_nickname)
);

drop function if exists addUserForum;
create function addUserForum() returns trigger as
$$
begin
    insert into user_forum (forum_slug, user_nickname)
    values (new.forum_slug, new.author_nickname)
    on conflict do nothing;
    return new;
end;
$$ language plpgsql;

drop function if exists threadsCounter;
create or replace function threadsCounter()
    returns trigger as
$$
begin
    update forums
    set threads = threads + 1
    where slug = new.forum_slug;

    return null;
end;
$$ language plpgsql;

drop trigger if exists threadsIncrementer on threads;
create trigger threadsIncrementer
    after insert
    on threads
    for each row
execute procedure threadsCounter();

drop trigger if exists threadsActivity on threads;
create trigger threadsActivity
    after insert
    on threads
    for each row
execute procedure addUserForum();

drop trigger if exists postsActivity on posts;
create trigger postsActivity
    after insert
    on posts
    for each row
execute procedure addUserForum();

CREATE INDEX ON threads (slug);
CREATE INDEX ON threads (created_at, forum_slug);
CREATE INDEX ON posts (thread_id);
CREATE INDEX ON posts (substring(path, 1, 8));
CREATE INDEX ON votes (thread_id, user_nickname);
