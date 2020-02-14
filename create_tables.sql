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
CREATE TABLE user (user_id INTEGER PRIMARY KEY NOT NULL, username TEXT, password TEXT, active INTEGER NOT NULL, CONSTRAINT unique_username UNIQUE (username));

DROP TABLE IF EXISTS site;
CREATE TABLE site (site_id INTEGER PRIMARY KEY NOT NULL, title TEXT, desc TEXT);

INSERT INTO user (user_id, username, password, active) VALUES (1, 'admin', '', 1);
INSERT INTO user (user_id, username, password, active) VALUES (2, 'guest', '', 1);

DROP TABLE IF EXISTS fts;
CREATE VIRTUAL TABLE fts USING FTS5(title, body, entry_id);

-- entry triggers
DROP TRIGGER IF EXISTS entry_after_insert;
CREATE TRIGGER entry_after_insert
AFTER INSERT ON entry
BEGIN
    INSERT INTO fts (title, body, entry_id)
    VALUES (new.title, new.body, new.entry_id);
END;

DROP TRIGGER IF EXISTS entry_after_update;
CREATE TRIGGER entry_after_update
AFTER UPDATE ON entry
WHEN old.title <> new.title OR old.body <> new.body
BEGIN
    UPDATE fts SET title = new.title, body = new.body
    WHERE fts.entry_id = new.entry_id;
END;

DROP TRIGGER IF EXISTS entry_after_delete;
CREATE TRIGGER entry_after_delete
AFTER DELETE ON entry
BEGIN
    DELETE FROM fts WHERE fts.entry_id = old.entry_id;
END;

COMMIT;

