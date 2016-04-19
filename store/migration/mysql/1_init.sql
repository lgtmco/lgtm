-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
 user_id      INTEGER PRIMARY KEY AUTO_INCREMENT
,user_login   VARCHAR(255)
,user_token   VARCHAR(255)
,user_email   VARCHAR(255)
,user_avatar  VARCHAR(1024)
,user_secret  VARCHAR(255)

,UNIQUE(user_login)
);

CREATE TABLE IF NOT EXISTS repos (
 repo_id       INTEGER PRIMARY KEY AUTO_INCREMENT
,repo_user_id  INTEGER
,repo_owner    VARCHAR(255)
,repo_name     VARCHAR(255)
,repo_slug     VARCHAR(255)
,repo_link     VARCHAR(1024)
,repo_private  BOOLEAN
,repo_secret   VARCHAR(255)

,UNIQUE(repo_slug)
);

CREATE INDEX ix_repo_owner   ON repos (repo_owner);
CREATE INDEX ix_repo_user_id ON repos (repo_user_id);

-- +migrate Down

DROP TABLE repos;
DROP TABLE users;
