## notes - Barebones group notes website. Inspired by PLATO Notes.

Still in development.

License:
  MIT

## Manual

Build and Install
-----------------
  $ make dep
  $ sqlite3 notes.db < create_tables.sql
  $ make

  Run 'groupnotes' to start the web service.

groupnotes uses a single sqlite3 database file to store all notes, uploaded files, users, and site settings.

## Screenshots

![notes list](screenshots/note_list.png)
![note text](screenshots/note_text.png)
![note with image](screenshots/note_with_image.png)
![note editor](screenshots/note_edit.png)
![files gallery](screenshots/files_gridview.png)

## Contact
  Twitter: @robdelacruz
  Source: http://github.com/robdelacruz/notes

