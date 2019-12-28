PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;
DROP TABLE IF EXISTS note;
CREATE TABLE note (note_id INTEGER PRIMARY KEY NOT NULL, title TEXT, body TEXT, createdt TEXT, user_id INTEGER, FOREIGN KEY(user_id) REFERENCES user);

DROP TABLE IF EXISTS notereply;
CREATE TABLE notereply (notereply_id INTEGER PRIMARY KEY NOT NULL, note_id INTEGER, body TEXT, createdt TEXT, user_id INTEGER, FOREIGN KEY(user_id) REFERENCES user, FOREIGN KEY(note_id) REFERENCES note);

DROP TABLE IF EXISTS user;
CREATE TABLE user (user_id INTEGER PRIMARY KEY NOT NULL, alias TEXT);

INSERT INTO user (alias) VALUES ('admin');
INSERT INTO user (alias) VALUES ('robdelacruz');

INSERT INTO note (title, body, createdt, user_id) VALUES ('Aimee Teagarden', 'All about Aimee Teagarden Hallmark show', '2019-12-01', 2);
INSERT INTO note (title, body, createdt, user_id) VALUES ('Emma Fielding', 'All about Emma Fielding Hallmark show', '2019-12-02', 2);
INSERT INTO note (title, body, createdt, user_id) VALUES ('Mystery 101', 'All about Mystery 101 Hallmark show', '2019-12-05', 2);

END TRANSACTION;
