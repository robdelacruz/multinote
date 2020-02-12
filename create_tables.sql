PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;

DROP TABLE IF EXISTS entry;
DROP TABLE IF EXISTS file;

-- thing enum values:
--   note = 0
--   reply = 1
--   file = 2

CREATE TABLE entry (entry_id INTEGER PRIMARY KEY NOT NULL, thing INTEGER NOT NULL, title TEXT, body TEXT, createdt TEXT, user_id INTEGER, parent_id INTEGER, FOREIGN KEY(user_id) REFERENCES user, FOREIGN KEY(parent_id) REFERENCES entry(entry_id));

CREATE TABLE file (entry_id INTEGER PRIMARY KEY NOT NULL, folder TEXT, content BLOB, FOREIGN KEY(entry_id) REFERENCES entry);

DROP TABLE IF EXISTS user;
CREATE TABLE user (user_id INTEGER PRIMARY KEY NOT NULL, username TEXT, password TEXT, CONSTRAINT unique_username UNIQUE (username));

DROP TABLE IF EXISTS site;
CREATE TABLE site (site_id INTEGER PRIMARY KEY NOT NULL, title TEXT, desc TEXT);

INSERT INTO user (user_id, username, password) VALUES (1, 'admin', '');
INSERT INTO user (user_id, username, password) VALUES (2, 'guest', '');


COMMIT;

