PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;
DROP TABLE IF EXISTS note;
CREATE TABLE note (note_id INTEGER PRIMARY KEY NOT NULL, title TEXT, body TEXT, createdt TEXT, user_id INTEGER, FOREIGN KEY(user_id) REFERENCES user);

DROP TABLE IF EXISTS notereply;
CREATE TABLE notereply (notereply_id INTEGER PRIMARY KEY NOT NULL, note_id INTEGER, replybody TEXT, createdt TEXT, user_id INTEGER, FOREIGN KEY(user_id) REFERENCES user, FOREIGN KEY(note_id) REFERENCES note);

DROP TABLE IF EXISTS user;
CREATE TABLE user (user_id INTEGER PRIMARY KEY NOT NULL, username TEXT, password TEXT);

INSERT INTO user (username, password) VALUES ('admin', 'password');
INSERT INTO user (username, password) VALUES ('robdelacruz', '123');

INSERT INTO note (title, body, createdt, user_id) VALUES ('Aimee Teagarden', 'All about Aimee Teagarden Hallmark show', '2019-12-01T14:00:00+08:00', 2);
INSERT INTO note (title, body, createdt, user_id) VALUES ('Emma Fielding', 'All about Emma Fielding Hallmark show', '2019-12-02T14:00:00+08:00', 2);
INSERT INTO note (title, body, createdt, user_id) VALUES ('Mystery 101', 'All about Mystery 101 Hallmark show', '2019-12-05T14:00:00+08:00', 2);

END TRANSACTION;
