create extension if not exists citext;

-- DROP INDEX IF EXISTS users_nickname_uindex;
-- DROP INDEX IF EXISTS users_email_uindex;
-- DROP INDEX IF EXISTS users_coverage_index;
-- DROP INDEX IF EXISTS forums_slug_uindex;;
-- DROP INDEX IF EXISTS forums_usernickname_idx;
-- DROP INDEX IF EXISTS threads_slug_idx;
-- DROP INDEX IF EXISTS threads_forum_idx;
-- DROP INDEX IF EXISTS threads_created_idx;
-- DROP INDEX IF EXISTS threads_created_forum_idx;
-- DROP INDEX IF EXISTS threads_coverage_idx;
-- DROP INDEX IF EXISTS posts_created_thread_idx;
-- DROP INDEX IF EXISTS posts_thread_idx;
-- DROP INDEX IF EXISTS posts_forum_idx;
-- DROP INDEX IF EXISTS posts_thread_id_idx;
-- DROP INDEX IF EXISTS posts_thread_path_idx;
-- DROP INDEX IF EXISTS forum_users_forum_idx;
-- DROP INDEX IF EXISTS forums_users_nickname_idx;
-- DROP INDEX IF EXISTS forums_users_forum_nickname_idx;
--
-- DROP INDEX IF EXISTS idx_posts_path;
-- DROP INDEX IF EXISTS idx_posts_parent;
-- DROP INDEX IF EXISTS idx_posts_thread_id;
-- DROP INDEX IF EXISTS idx_posts_pok;
-- DROP INDEX IF EXISTS idx_posts_created;
-- DROP INDEX IF EXISTS idx_votes_nickname_thread_unique2;

DROP TABLE IF EXISTS forum_users;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS users;


CREATE TABLE IF NOT EXISTS users
(
    id       BIGSERIAL    NOT NULL
        CONSTRAINT users_pk PRIMARY KEY,
    nickname citext COLLATE "POSIX" NOT NULL UNIQUE ,
    email    citext NOT NULL,
    fullname TEXT         NOT NULL,
    about    TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS users_nickname_uindex
    ON users (nickname);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_uindex
    ON users (email);

CREATE INDEX IF NOT EXISTS users_coverage_index
    ON users (nickname, email, fullname, about);

DROP TABLE IF EXISTS forums;
CREATE TABLE IF NOT EXISTS forums
(
    id      BIGSERIAL    NOT NULL PRIMARY KEY,
    slug    citext NOT NULL,
    title   VARCHAR(128) NOT NULL,
    nickname citext NOT NULL,
    posts   BIGINT       NOT NULL DEFAULT 0,
    threads INT          NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS forums_slug_uindex
    ON forums (slug);

CREATE INDEX IF NOT EXISTS forums_usernickname_idx
    ON forums (nickname);

CREATE TABLE IF NOT EXISTS threads
(
    id       BIGSERIAL    NOT NULL PRIMARY KEY,
    forum    citext NOT NULL,
    author   VARCHAR(128) NOT NULL,
    created  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    message  TEXT         NOT NULL,
    title    VARCHAR(128) NOT NULL,
    votes    INT          NOT NULL DEFAULT 0,
    slug     citext NOT NULL
);

CREATE INDEX IF NOT EXISTS threads_slug_idx
    ON threads (slug);

CREATE INDEX IF NOT EXISTS threads_forum_idx
    ON threads (forum);

CREATE INDEX IF NOT EXISTS threads_created_idx
    ON threads (created);

CREATE INDEX IF NOT EXISTS threads_created_forum_idx
    ON threads (created, forum);

CREATE INDEX IF NOT EXISTS threads_coverage_idx
    ON threads (id, forum, author, slug, created, title, message, votes);

CREATE TABLE IF NOT EXISTS posts
(
    id        BIGSERIAL   NOT NULL PRIMARY KEY,
    parent    BIGINT      DEFAULT NULL,
    thread    INT REFERENCES threads(id) NOT NULL,
    forum     citext NOT NULL,
    author    citext NOT NULL,
    created   TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_edited BOOLEAN     NOT NULL DEFAULT FALSE,
    message   TEXT        NOT NULL,
    path      BIGINT[]  NOT NULL DEFAULT '{0}'
);

CREATE INDEX IF NOT EXISTS posts_created_thread_idx
    ON posts (created, thread);

CREATE INDEX IF NOT EXISTS posts_thread_idx
    ON posts (thread);

CREATE INDEX IF NOT EXISTS posts_thread_id_idx
    ON posts (thread, id);

CREATE INDEX IF NOT EXISTS posts_forum_idx
    ON posts (forum);

CREATE INDEX IF NOT EXISTS posts_thread_path_idx
    ON posts (thread, path);

CREATE INDEX IF NOT EXISTS idx_posts_path
    ON posts USING GIN (path);

CREATE INDEX IF NOT EXISTS idx_posts_parent
    ON posts (parent);

CREATE INDEX IF NOT EXISTS idx_posts_thread_id
    ON posts (thread, id);

CREATE INDEX IF NOT EXISTS idx_posts_pok
    ON posts (id, parent, thread, forum, author, created, message, is_edited, path);

CREATE INDEX IF NOT EXISTS idx_posts_created
    ON posts (created);

CREATE TABLE IF NOT EXISTS votes
(
    nickname CITEXT        REFERENCES users(nickname) NOT NULL,
    thread   INT           REFERENCES threads(id) NOT NULL,
    voice    SMALLINT      NOT NULL CHECK (voice = 1 OR voice = -1),
    PRIMARY KEY (nickname, thread)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_nickname_thread_unique2
    ON votes (nickname, thread);

CREATE TABLE IF NOT EXISTS forum_users (
           forum CITEXT NOT NULL,
           nickname CITEXT NOT NULL,
           FOREIGN KEY (forum) REFERENCES forums (slug),
           FOREIGN KEY (nickname) REFERENCES users (nickname),
           UNIQUE (forum, nickname)
);

CREATE INDEX IF NOT EXISTS forum_users_forum_idx ON forum_users (forum);
CREATE INDEX IF NOT EXISTS forums_users_nickname_idx ON forum_users (nickname);
CREATE INDEX IF NOT EXISTS forums_users_forum_nickname_idx ON forum_users (forum, nickname);

-- CREATE TABLE IF NOT EXISTS forum_users (
--    user_id BIGINT REFERENCES users(id),
--    forum_id BIGINT REFERENCES forums(id),
--    UNIQUE (user_id, forum_id)
-- );

-- CREATE FUNCTION  on_forum_users_update()
--     RETURNS TRIGGER AS '
--     BEGIN
--         INSERT INTO forum_users (user_id, forum_id) VALUES ((SELECT id FROM users WHERE NEW.author = nickname),
--                                                             (SELECT id FROM forums WHERE NEW.forum = slug));
--         RETURN NULL;
--     END;
-- ' LANGUAGE plpgsql;
--
-- CREATE TRIGGER on_new_thread_inserted
--     AFTER INSERT ON threads
--     FOR EACH ROW EXECUTE PROCEDURE on_forum_users_update();
