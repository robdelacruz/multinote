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


-- fulltext search (fts) table and triggers to update fts
DROP TABLE IF EXISTS fts;
CREATE VIRTUAL TABLE fts USING FTS5(title, body, user_id, thing, thing_id);

-- enum 'thing'
-- note = 0
-- notereply = 1
-- file = 2

-- note triggers
DROP TRIGGER IF EXISTS note_after_insert;
CREATE TRIGGER note_after_insert
AFTER INSERT ON note
BEGIN
    INSERT INTO fts (title, body, user_id, thing, thing_id)
    VALUES (new.title, new.body, new.user_id, 0, new.note_id);
END;

DROP TRIGGER IF EXISTS note_after_update;
CREATE TRIGGER note_after_update
AFTER UPDATE ON note
WHEN old.title <> new.title OR old.body <> new.body
BEGIN
    UPDATE fts SET title = new.title, body = new.body
    WHERE thing = 0 AND thing_id = new.note_id;
END;

DROP TRIGGER IF EXISTS note_after_delete;
CREATE TRIGGER note_after_delete
AFTER DELETE ON note
BEGIN
    DELETE FROM fts WHERE thing = 0 AND thing_id = old.note_id;
END;

-- notereply triggers
DROP TRIGGER IF EXISTS notereply_after_insert;
CREATE TRIGGER notereply_after_insert
AFTER INSERT ON notereply
BEGIN
    INSERT INTO fts (body, user_id, thing, thing_id)
    VALUES (new.replybody, new.user_id, 1, new.notereply_id);
END;

DROP TRIGGER IF EXISTS notereply_after_update;
CREATE TRIGGER notereply_after_update
AFTER UPDATE ON notereply
WHEN old.replybody <> new.replybody
BEGIN
    UPDATE fts SET body = new.replybody
    WHERE thing = 1 AND thing_id = new.notereply_id;
END;

DROP TRIGGER IF EXISTS notereply_after_delete;
CREATE TRIGGER notereply_after_delete
AFTER DELETE ON notereply
BEGIN
    DELETE FROM fts WHERE thing = 1 AND thing_id = old.notereply_id;
END;

-- file triggers
DROP TRIGGER IF EXISTS file_after_insert;
CREATE TRIGGER file_after_insert
AFTER INSERT ON file
BEGIN
    INSERT INTO fts (title, body, user_id, thing, thing_id)
    VALUES (new.filename, new.desc, new.user_id, 2, new.file_id);
END;

DROP TRIGGER IF EXISTS file_after_update;
CREATE TRIGGER file_after_update
AFTER UPDATE ON file
WHEN old.filename <> new.filename OR old.desc <> new.desc
BEGIN
    UPDATE fts SET title = new.filename, body = new.desc
    WHERE thing = 2 AND thing_id = new.file_id;
END;

DROP TRIGGER IF EXISTS file_after_delete;
CREATE TRIGGER file_after_delete
AFTER DELETE ON file
BEGIN
    DELETE FROM fts WHERE thing = 2 AND thing_id = old.file_id;
END;

COMMIT;

