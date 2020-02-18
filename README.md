## About GroupNotes

GroupNotes is a multi-user web based note sharing system. Inspired by PLATO Notes. 

- Add and edit notes.
- Upload files, images.
- Reply to notes.
- Multiple users.
- Use to keep track of records, as a weblog, or CMS.
- MIT License

## Build and Install

    $ make dep
    $ sqlite3 notes.db < create_tables.sql
    $ make
    $ ./groupnotesd -i notes.db

    Run 'groupnotesd <notes_file>' to start the web service.

groupnotes uses a single sqlite3 database file to store all notes, uploaded files, users, and site settings.

## Screenshots

![notes list](screenshots/note_list.png)
![note text](screenshots/note_text.png)
![note with image](screenshots/note_with_image.png)
![note editor](screenshots/note_edit.png)
![files gallery](screenshots/files_gridview.png)

## Contact
    Twitter: @robdelacruz
    Source: http://github.com/robdelacruz/groupnotes

