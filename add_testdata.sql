PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;
DELETE FROM notereply WHERE 1=1;
DELETE FROM note WHERE 1=1;
DELETE FROM file WHERE 1=1;
DELETE FROM user WHERE user_id > 2;

INSERT INTO user (username, password) VALUES ('robdelacruz', '$2a$10$QBKdo66QfkyqNczexwGFwul3731pQ970B96Bn1hgmvXLBu.LaJhFK'); -- password is '123'
INSERT INTO user (username, password) VALUES ('lky', '');

INSERT INTO note (title, body, createdt, user_id) VALUES ('Aimee Teagarden', 'All about Aimee Teagarden Hallmark show', '2019-12-01T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('Emma Fielding', 'All about Emma Fielding Hallmark show', '2019-12-02T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('Mystery 101', 'All about Mystery 101 Hallmark show', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('test note', 'test note 1', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('test note 2', 'test note 2', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('test note 3', 'test note 3', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO note (title, body, createdt, user_id) VALUES ('markdown test', '# Gettysburg Address

*Versions*

- Bliss Copy
- Nicolay Copy
- Hay Copy
- Everett Copy
- Bancroft Copy

### Related Links

[Robert Todd Lincoln''s "Gettysburg Story"](https://quod.lib.umich.edu/j/jala/2629860.0038.103/--robert-todd-lincolns-gettysburg-story?rgn=main;view=fulltext) (JALA)
[Who stole the Gettysburg Address?](https://quod.lib.umich.edu/j/jala/2629860.0024.203/--who-stole-the-gettysburg-address?rgn=main;view=fulltext) (JALA)

---

Four score and seven years ago our fathers brought forth on this continent, a new nation, conceived in Liberty, and dedicated to the proposition that all men are created equal.

Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting place for those who here gave their lives that that nation might live. It is altogether fitting and proper that we should do this.

But, in a larger sense, we can not dedicate -- we can not consecrate -- we can not hallow -- this ground. The brave men, living and dead, who struggled here, have consecrated it, far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us -- that from these honored dead we take increased devotion to that cause for which they gave the last full measure of devotion -- that we here highly resolve that these dead shall not have died in vain -- that this nation, under God, shall have a new birth of freedom -- and that government of the people, by the people, for the people, shall not perish from the earth.

Abraham Lincoln
November 19, 1863

![Soldiers National Cemetery](http://www.abrahamlincolnonline.org/lincoln/sites/gettycem.jpg)

[source](http://www.abrahamlincolnonline.org/lincoln/speeches/gettysburg.htm)
', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO note (title, body, createdt, user_id) VALUES ('code test', '## Hello, World

Code for Hello, World:

    #include <stdio.h>

    int main() {
        printf("Hello, World!\n");
    }
', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO note (title, body, createdt, user_id) VALUES ('Lee Kuan Yew Quotes', '### Lee Kuan Yew Quotes:

>"If there was one formula for our success,it was that we were constantly studying how to make things work,or how to make them work better."

>"I’m very determined. If I decide that something is worth doing, then I’ll put my heart and soul to it. The whole ground can be against me, but if I know it is right, I’ll do it."

>"I always tried to be correct, not politically correct."

', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (5, 'first comment!', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (5, 'second comment!', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (5, 'third comment!', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (6, 'a comment!', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (6, 'another comment!', '2019-12-05T14:00:00+08:00', 3);

COMMIT;

