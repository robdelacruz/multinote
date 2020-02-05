PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;

DROP TABLE IF EXISTS notereply;
DROP TABLE IF EXISTS note;
DROP TABLE IF EXISTS file;

CREATE TABLE note (note_id INTEGER PRIMARY KEY NOT NULL, title TEXT, body TEXT, createdt TEXT, user_id INTEGER, FOREIGN KEY(user_id) REFERENCES user);

CREATE TABLE notereply (notereply_id INTEGER PRIMARY KEY NOT NULL, note_id INTEGER, replybody TEXT, createdt TEXT, user_id INTEGER, FOREIGN KEY(user_id) REFERENCES user, FOREIGN KEY(note_id) REFERENCES note);

CREATE TABLE file (file_id INTEGER PRIMARY KEY NOT NULL, filename TEXT, folder TEXT, desc TEXT, content BLOB, createdt TEXT, user_id INTEGER, FOREIGN KEY(user_id) REFERENCES user);

DROP TABLE IF EXISTS user;
CREATE TABLE user (user_id INTEGER PRIMARY KEY NOT NULL, username TEXT, password TEXT, CONSTRAINT unique_username UNIQUE (username));

DROP TABLE IF EXISTS site;
CREATE TABLE site (site_id INTEGER PRIMARY KEY NOT NULL, title TEXT, desc TEXT);

INSERT INTO user (user_id, username, password) VALUES (1, 'admin', '');
INSERT INTO user (user_id, username, password) VALUES (2, 'guest', '');

COMMIT;

