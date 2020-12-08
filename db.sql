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
    ON users (LOWER(nickname));

CREATE UNIQUE INDEX IF NOT EXISTS users_email_uindex
    ON users (LOWER(email));


DROP TABLE IF EXISTS forums;
CREATE TABLE IF NOT EXISTS forums
(
    id      BIGSERIAL    NOT NULL PRIMARY KEY,
    slug    VARCHAR(128) NOT NULL,
    title   VARCHAR(128) NOT NULL,
    nickname VARCHAR(128) NOT NULL,
    posts   BIGINT       NOT NULL DEFAULT 0,
    threads INT          NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS forums_slug_uindex
    ON forums (LOWER(slug));

CREATE TABLE IF NOT EXISTS threads
(
    id       BIGSERIAL    NOT NULL PRIMARY KEY,
    forum    VARCHAR(128) NOT NULL,
    author   VARCHAR(128) NOT NULL,
    created  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    message  TEXT         NOT NULL,
    title    VARCHAR(128) NOT NULL,
    votes    INT          NOT NULL DEFAULT 0,
    slug     VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id        BIGSERIAL   NOT NULL PRIMARY KEY,
    parent    BIGINT      DEFAULT NULL,
    thread    INT REFERENCES threads(id) NOT NULL,
    forum     VARCHAR(128) NOT NULL,
    author    VARCHAR(128) NOT NULL,
    created   TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_edited BOOLEAN     NOT NULL DEFAULT FALSE,
    message   TEXT        NOT NULL,
    path      BIGINT[]  NOT NULL DEFAULT '{0}'
);

CREATE TABLE IF NOT EXISTS votes
(
    nickname VARCHAR(128)  REFERENCES users(nickname) NOT NULL,
    thread   INT           REFERENCES threads(id) NOT NULL,
    voice    SMALLINT      NOT NULL CHECK (voice = 1 OR voice = -1),
    PRIMARY KEY (nickname, thread)
);

CREATE TABLE IF NOT EXISTS forum_users (
                                           user_id BIGINT REFERENCES users(id),
                                           forum_id BIGINT REFERENCES forums(id)
);

CREATE FUNCTION  on_forum_users_update()
    RETURNS TRIGGER AS '
    BEGIN
        INSERT INTO forum_users (user_id, forum_id) VALUES ((SELECT id FROM users WHERE LOWER(NEW.author) = LOWER(nickname)),
                                                            (SELECT id FROM forums WHERE LOWER(NEW.forum) = LOWER(slug)));
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_new_thread_inserted
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE on_forum_users_update();
