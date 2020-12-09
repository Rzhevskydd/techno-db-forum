create extension if not exists citext;

CREATE TABLE IF NOT EXISTS users
(
    id       BIGSERIAL    NOT NULL
        CONSTRAINT users_pk PRIMARY KEY,
    nickname citext  NOT NULL,
    email    citext NOT NULL,
    fullname TEXT         NOT NULL,
    about    TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS users_nickname_uindex
    ON users (nickname);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_uindex
    ON users (email);


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

CREATE TABLE IF NOT EXISTS votes
(
    nickname CITEXT        REFERENCES users(nickname) NOT NULL,
    thread   INT           REFERENCES threads(id) NOT NULL,
    voice    SMALLINT      NOT NULL CHECK (voice = 1 OR voice = -1),
    PRIMARY KEY (nickname, thread)
);

CREATE TABLE IF NOT EXISTS forum_users (
                                           forum CITEXT NOT NULL,
                                           nickname CITEXT NOT NULL,
                                           FOREIGN KEY (forum) REFERENCES forums (slug),
                                           FOREIGN KEY (nickname) REFERENCES users (nickname),
                                           UNIQUE (forum, nickname)
);

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
