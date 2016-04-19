-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
 user_id      INTEGER PRIMARY KEY AUTOINCREMENT
,user_login   TEXT
,user_token   TEXT
,user_email   TEXT
,user_avatar  TEXT
,user_secret  TEXT

,UNIQUE(user_login)
);

CREATE TABLE IF NOT EXISTS repos (
 repo_id       INTEGER PRIMARY KEY AUTOINCREMENT
,repo_user_id  INTEGER
,repo_owner    TEXT
,repo_name     TEXT
,repo_slug     TEXT
,repo_link     TEXT
,repo_private  BOOLEAN
,repo_secret   TEXT

,UNIQUE(repo_slug)
);

CREATE INDEX IF NOT EXISTS ix_repo_owner   ON repos (repo_owner);
CREATE INDEX IF NOT EXISTS ix_repo_user_id ON repos (repo_user_id);

-- +migrate Down

DROP TABLE repos;
DROP TABLE users;
