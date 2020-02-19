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
CREATE TABLE user (user_id INTEGER PRIMARY KEY NOT NULL, username TEXT, password TEXT, active INTEGER NOT NULL, mdeditor INTEGER NOT NULL, CONSTRAINT unique_username UNIQUE (username));

DROP TABLE IF EXISTS site;
CREATE TABLE site (site_id INTEGER PRIMARY KEY NOT NULL, title TEXT NOT NULL, desc TEXT NOT NULL, requireloginforpageview INTEGER NOT NULL, allowanonreplies INTEGER NOT NULL, loginmsg TEXT, sidebar1 TEXT NOT NULL, sidebar2 TEXT NOT NULL);

INSERT INTO user (user_id, username, password, active, mdeditor) VALUES (1, 'admin', '', 1, 0);
INSERT INTO user (user_id, username, password, active, mdeditor) VALUES (2, 'guest', '', 1, 0);

INSERT INTO site (site_id, title, desc, requireloginforpageview, allowanonreplies, loginmsg, sidebar1, sidebar2) VALUES (1, 'Group Notes', 'Repository for Notes', 0, 0,  '', 
'## About GroupNotes
GroupNotes is a multi-user web based note sharing system. Inspired by PLATO Notes. 

- Add and edit notes.
- Upload files, images.
- Reply to notes.
- Multiple users.
- Use to keep track of records, as a weblog, or CMS.
- MIT License
',
'<div style="color: #d6deeb; background-color: #2d2c5d; padding: 0.5em; margin: 0; border: 1px solid; text-align: center;">
Donate to the Developer
<form action="https://www.paypal.com/cgi-bin/webscr" method="post" target="_top"><input type="hidden" name="cmd" value="_s-xclick" /><input type="hidden" name="hosted_button_id" value="N5GLMFNSW9UVN" /><input type="image" src="https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif" border="0" name="submit" title="PayPal - The safer, easier way to pay online!" alt="Donate with PayPal button" /><img alt="" border="0" src="https://www.paypal.com/en_PH/i/scr/pixel.gif" width="1" height="1" /></form>
</div>');

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

