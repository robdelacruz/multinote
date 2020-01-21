package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/russross/blackfriday.v2"
)

const ADMIN_ID = 1

type User struct {
	Userid   int64
	Username string
}

func main() {
	port := "8000"

	db, err := sql.Open("sqlite3", "notes.db")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", notesHandler(db))
	http.HandleFunc("/note/", noteHandler(db))
	http.HandleFunc("/createnote/", createNoteHandler(db))
	http.HandleFunc("/editnote/", editNoteHandler(db))
	http.HandleFunc("/delnote/", delNoteHandler(db))
	http.HandleFunc("/file/", fileHandler(db))
	http.HandleFunc("/uploadfile/", uploadFileHandler(db))
	//http.HandleFunc("/editfile/", editFileHandler(db))
	//http.HandleFunc("/delfile/", delFileHandler(db))
	http.HandleFunc("/newreply/", newReplyHandler(db))
	http.HandleFunc("/editreply/", editReplyHandler(db))
	http.HandleFunc("/delreply/", delReplyHandler(db))
	http.HandleFunc("/login/", loginHandler(db))
	http.HandleFunc("/logout/", logoutHandler(db))
	http.HandleFunc("/adminsetup/", adminsetupHandler(db))
	http.HandleFunc("/usersettings/", usersettingsHandler(db))
	http.HandleFunc("/newuser/", newUserHandler(db))
	http.HandleFunc("/edituser/", editUserHandler(db))
	http.HandleFunc("/sitesettings/", sitesettingsHandler(db))
	http.HandleFunc("/userssetup/", userssetupHandler(db))
	fmt.Printf("Listening on %s...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}

func parseNoteid(url string) int64 {
	sre := `^/note/(\d+)$`
	re := regexp.MustCompile(sre)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return -1
	}
	return idtoi(matches[1])
}

func notesHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getLoginUser(r, db)

		w.Header().Set("Content-Type", "text/html")

		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<ul class=\"vertical-list\">\n")
		s := `SELECT note_id, title, body, createdt, user.user_id, username, 
(SELECT COUNT(*) FROM notereply WHERE note.note_id = notereply.note_id) AS numreplies, 
(SELECT COALESCE(MAX(createdt), '') FROM notereply where note.note_id = notereply.note_id) AS maxreplydt 
FROM note 
LEFT OUTER JOIN user ON note.user_id = user.user_id 
ORDER BY MAX(createdt, maxreplydt) DESC;`
		rows, err := db.Query(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		for rows.Next() {
			var noteid int64
			var title, body, createdt, maxreplydt string
			var noteUser User
			var numreplies int
			rows.Scan(&noteid, &title, &body, &createdt, &noteUser.Userid, &noteUser.Username, &numreplies, &maxreplydt)
			tcreatedt, _ := time.Parse(time.RFC3339, createdt)

			fmt.Fprintf(w, "<li>\n")
			fmt.Fprintf(w, "<p class=\"doc-title\"><a href=\"/note/%d\">%s</a></p>\n", noteid, title)

			printByline(w, login, noteid, noteUser, tcreatedt, numreplies)
			fmt.Fprintf(w, "</li>\n")
		}
		fmt.Fprintf(w, "</ul>\n")

		printPageFoot(w)
	}
}

func noteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		noteid := parseNoteid(r.URL.Path)
		if noteid == -1 {
			log.Printf("display note: no noteid\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		login := getLoginUser(r, db)

		s := `SELECT title, body, createdt, user.user_id, username, (SELECT COUNT(*) FROM notereply WHERE note.note_id = notereply.note_id) AS numreplies 
FROM note
LEFT OUTER JOIN user ON note.user_id = user.user_id
WHERE note_id = ?
ORDER BY createdt DESC;`
		row := db.QueryRow(s, noteid)

		var title, body, createdt string
		var noteUser User
		var numreplies int
		err := row.Scan(&title, &body, &createdt, &noteUser.Userid, &noteUser.Username, &numreplies)
		if err == sql.ErrNoRows {
			// note doesn't exist so redirect to notes list page.
			log.Printf("display note: noteid %d doesn't exist\n", noteid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<article class=\"content\">\n")
		fmt.Fprintf(w, "<h1 class=\"heading doc-title\"><a href=\"/note/%d\">%s</a></h1>", noteid, title)
		tcreatedt, err := time.Parse(time.RFC3339, createdt)
		if err != nil {
			tcreatedt = time.Now()
		}
		printByline(w, login, noteid, noteUser, tcreatedt, numreplies)

		bodyMarkup := parseMarkdown(body)
		fmt.Fprintf(w, bodyMarkup)

		fmt.Fprintf(w, "<hr class=\"dotted\">\n")
		fmt.Fprintf(w, "<p>Replies:</p>\n")

		s = "SELECT notereply_id, replybody, createdt, user.user_id, username FROM notereply LEFT OUTER JOIN user ON notereply.user_id = user.user_id WHERE note_id = ? ORDER BY notereply_id"
		rows, err := db.Query(s, noteid)
		if err != nil {
			fmt.Fprintf(w, "<p class=\"error\">Error loading replies</p>\n")
			fmt.Fprintf(w, "</article>\n")
			printPageFoot(w)
			return
		}
		i := 1
		for rows.Next() {
			var replyid int64
			var replybody, createdt string
			var replyUser User
			rows.Scan(&replyid, &replybody, &createdt, &replyUser.Userid, &replyUser.Username)
			tcreatedt, _ := time.Parse(time.RFC3339, createdt)
			createdt = tcreatedt.Format("2 Jan 2006")

			fmt.Fprintf(w, "<div class=\"reply compact\">\n")
			fmt.Fprintf(w, "<ul class=\"line-menu finetext\">\n")
			fmt.Fprintf(w, "<li>%d. %s</li>\n", i, replyUser.Username)
			fmt.Fprintf(w, "<li>%s</li>\n", createdt)
			if replyUser.Userid == login.Userid || login.Userid == ADMIN_ID {
				fmt.Fprintf(w, "<li><a href=\"/editreply/?replyid=%d\">Edit</a></li>\n", replyid)
				fmt.Fprintf(w, "<li><a href=\"/delreply/?replyid=%d\">Delete</a></li>\n", replyid)
			}
			fmt.Fprintf(w, "</ul>\n")
			replybodyMarkup := parseMarkdown(replybody)
			fmt.Fprintf(w, replybodyMarkup)
			fmt.Fprintf(w, "</div>\n")

			i++
		}
		fmt.Fprintf(w, "</article>\n")

		// New Reply form
		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/newreply/?noteid=%d\" method=\"post\">\n", noteid)
		if login.Userid == -1 {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<label><a href=\"/login/\">Log in</a> to post a reply.</label>")
			fmt.Fprintf(w, "</div>\n")
		} else {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<label>reply as %s:</label>\n", login.Username)
			fmt.Fprintf(w, "<textarea name=\"replybody\" rows=\"10\" cols=\"80\"></textarea>\n")
			fmt.Fprintf(w, "</div>\n")

			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<button class=\"submit\">add reply</button>\n")
			fmt.Fprintf(w, "</div>\n")
			fmt.Fprintf(w, "</form>\n")
		}

		printPageFoot(w)
	}
}

func createNoteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string
		var title, body string

		login := getLoginUser(r, db)
		if login.Userid == -1 {
			log.Printf("create note: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			title = r.FormValue("title")
			body = r.FormValue("body")
			body = strings.ReplaceAll(body, "\r", "") // CRLF => CR
			createdt := time.Now().Format(time.RFC3339)

			s := "INSERT INTO note (title, body, createdt, user_id) VALUES (?, ?, ?, ?);"
			_, err := sqlstmt(db, s).Exec(title, body, createdt, login.Userid)
			if err != nil {
				log.Printf("DB error creating note: %s\n", err)
				errmsg = "A problem occured. Please try again."
			}

			if errmsg == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/createnote/\" method=\"post\">\n")
		fmt.Fprintf(w, "<h1 class=\"heading\">Create Note</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>title</label>\n")
		fmt.Fprintf(w, "<input name=\"title\" type=\"text\" size=\"50\" value=\"%s\">\n", title)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>note</label>\n")
		fmt.Fprintf(w, "<textarea name=\"body\" rows=\"25\" cols=\"80\">%s</textarea>\n", body)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">create note</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func editNoteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string

		noteid := idtoi(r.FormValue("noteid"))
		if noteid == -1 {
			log.Printf("edit note: no noteid\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		login := getLoginUser(r, db)
		if login.Userid == -1 {
			log.Printf("edit note: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		s := "SELECT title, body, user_id FROM note WHERE note_id = ?"
		row := db.QueryRow(s, noteid)

		var title, body string
		var noteUserid int64
		err := row.Scan(&title, &body, &noteUserid)
		if err == sql.ErrNoRows {
			// note doesn't exist so redirect to notes list page.
			log.Printf("noteid %d doesn't exist\n", noteid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Allow only creators or admin to edit the note.
		if noteUserid != login.Userid && login.Userid != ADMIN_ID {
			log.Printf("User '%s' doesn't have access to note %d\n", login.Username, noteid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			title = r.FormValue("title")
			body = r.FormValue("body")
			createdt := time.Now().Format(time.RFC3339)

			// Strip out linefeed chars so that CRLF becomes just CR.
			// CRLF causes problems in markdown parsing.
			body = strings.ReplaceAll(body, "\r", "")

			s := "UPDATE note SET title = ?, body = ?, createdt = ? WHERE note_id = ?"
			_, err = sqlstmt(db, s).Exec(title, body, createdt, noteid)
			if err != nil {
				log.Printf("DB error updating noteid %d: %s\n", noteid, err)
				errmsg = "A problem occured. Please try again."
			}

			if errmsg == "" {
				http.Redirect(w, r, fmt.Sprintf("/note/%d", noteid), http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/editnote/?noteid=%d\" method=\"post\">\n", noteid)
		fmt.Fprintf(w, "<h1 class=\"heading\">Edit Note</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>title</label>\n")
		fmt.Fprintf(w, "<input name=\"title\" type=\"text\" size=\"50\" value=\"%s\">\n", title)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>note</label>\n")
		fmt.Fprintf(w, "<textarea name=\"body\" rows=\"25\" cols=\"80\">%s</textarea>\n", body)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">update note</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func delNoteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string

		noteid := idtoi(r.FormValue("noteid"))
		if noteid == -1 {
			log.Printf("del note: no noteid\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		login := getLoginUser(r, db)
		if login.Userid == -1 {
			log.Printf("del note: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		s := "SELECT title, body, user_id FROM note WHERE note_id = ?"
		row := db.QueryRow(s, noteid)

		var title, body string
		var noteUserid int64
		err := row.Scan(&title, &body, &noteUserid)
		if err == sql.ErrNoRows {
			// note doesn't exist so redirect to notes list page.
			log.Printf("noteid %d doesn't exist\n", noteid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Allow only creators or admin to delete the note.
		if noteUserid != login.Userid && login.Userid != ADMIN_ID {
			log.Printf("User '%s' doesn't have access to note %d\n", login.Username, noteid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			// todo: use transaction or trigger instead?

			for {
				// Delete note replies
				s = "DELETE FROM notereply WHERE note_id = ?"
				_, err = sqlstmt(db, s).Exec(noteid)
				if err != nil {
					log.Printf("DB error deleting notereplies of noteid %d: %s\n", noteid, err)
					errmsg = "A problem occured. Please try again."
					break
				}

				// Delete note
				s := "DELETE FROM note WHERE note_id = ?"
				_, err = sqlstmt(db, s).Exec(noteid)
				if err != nil {
					log.Printf("DB error deleting noteid %d: %s\n", noteid, err)
					errmsg = "A problem occured. Please try again."
					break
				}

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform displayonly\" action=\"/delnote/?noteid=%d\" method=\"post\">\n", noteid)
		fmt.Fprintf(w, "<h1 class=\"heading warning\">Delete Note</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label class=\"byline\">title</label>\n")
		fmt.Fprintf(w, "<input name=\"title\" type=\"text\" size=\"50\" readonly value=\"%s\">\n", title)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label class=\"byline\">note</label>\n")
		fmt.Fprintf(w, "<textarea name=\"body\" rows=\"25\" cols=\"80\" readonly>%s</textarea>\n", body)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit warning\">delete note</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func parseFileid(url string) int64 {
	sre := `^/file/(\d+)$`
	re := regexp.MustCompile(sre)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return -1
	}
	return idtoi(matches[1])
}

func parseUrlFilepath(url string) (string, string) {
	sre := `^/file/(.+)$`
	re := regexp.MustCompile(sre)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return "", ""
	}

	ss := strings.Split(matches[1], "/")
	if len(ss) == 1 {
		return "", ss[0]
	}
	path := strings.Join(ss[:len(ss)-1], "/")
	if len(path) == 0 {
		path = ""
	}
	name := ss[len(ss)-1]
	return path, name
}

// Return filename extension. Ex. "image.png" returns "png", "file1" returns "".
func fileext(filename string) string {
	ss := strings.Split(filename, ".")
	if len(ss) < 2 {
		return ""
	}
	return ss[len(ss)-1]
}

func fileHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var qname, qpath string
		fileid := parseFileid(r.URL.Path)
		if fileid == -1 {
			qpath, qname = parseUrlFilepath(r.URL.Path)
			if qname == "" {
				log.Printf("display file: no fileid or filepath\n")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		//login := getLoginUser(r, db)
		// todo: authenticate user?

		var row *sql.Row
		var name, path string
		var fileUserid int64
		var bsContent []byte
		if qname != "" {
			s := "SELECT name, path, content, user_id FROM file WHERE name = ? AND path = ?"
			row = db.QueryRow(s, qname, qpath)
		} else {
			s := "SELECT name, path, content, user_id FROM file WHERE file_id = ?"
			row = db.QueryRow(s, fileid)
		}

		err := row.Scan(&name, &path, &bsContent, &fileUserid)
		if err == sql.ErrNoRows {
			// file doesn't exist
			log.Printf("display file: file doesn't exist\n")
			http.Error(w, fmt.Sprintf("file doesn't exist"), 404)
			return
		} else if err != nil {
			log.Printf("display file: server error (%s)\n", err)
			http.Error(w, fmt.Sprintf("server error (%s)", err), 500)
			return
		}

		ext := fileext(name)
		if ext == "" {
			w.Header().Set("Content-Type", "application")
		} else if ext == "png" || ext == "gif" || ext == "bmp" {
			w.Header().Set("Content-Type", fmt.Sprintf("image/%s", ext))
		} else if ext == "jpg" {
			w.Header().Set("Content-Type", fmt.Sprintf("image/jpeg"))
		} else {
			w.Header().Set("Content-Type", fmt.Sprintf("application/%s", ext))
		}

		_, err = w.Write(bsContent)
		if err != nil {
			log.Printf("display file: server error (%s)\n", err)
			http.Error(w, fmt.Sprintf("server error (%s)", err), 500)
			return
		}
	}
}

func uploadFileHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string
		var successmsg string
		var name, path string

		login := getLoginUser(r, db)
		if login.Userid == -1 {
			log.Printf("upload file: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			for {
				file, header, err := r.FormFile("file")
				if err != nil {
					log.Printf("uploadfile: IO error reading file: %s\n", err)
					errmsg = "A problem occured. Please try again."
					break
				}
				defer file.Close()

				createdt := time.Now().Format(time.RFC3339)
				name = r.FormValue("name")
				path = r.FormValue("path")
				if name == "" {
					name = header.Filename
				}
				// Strip out any leading or trailing "/" from path.
				// Ex. "/abc/dir/" becomes "abc/dir".
				path = strings.TrimPrefix(path, "/")
				path = strings.TrimSuffix(path, "/")

				bsContent, err := ioutil.ReadAll(file)
				if err != nil {
					log.Printf("uploadfile: IO error reading file: %s\n", err)
					errmsg = "A problem occured. Please try again."
					break
				}

				s := "INSERT INTO file (name, path, content, createdt, user_id) VALUES (?, ?, ?, ?, ?);"
				_, err = sqlstmt(db, s).Exec(name, path, bsContent, createdt, login.Userid)
				if err != nil {
					log.Printf("uploadfile: DB error inserting file: %s\n", err)
					errmsg = "A problem occured. Please try again."
					break
				}

				// Successfully added file.
				filepath := name
				if path != "" {
					filepath = fmt.Sprintf("%s/%s", path, name)
				}
				link := fmt.Sprintf("<a href=\"/file/%s\">%s</a>", filepath, filepath)
				successmsg = fmt.Sprintf("Successfully added file: %s", link)
				name = ""
				path = ""
				break
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/uploadfile/\" method=\"post\" enctype=\"multipart/form-data\">\n")
		fmt.Fprintf(w, "<h1 class=\"heading\">Upload File</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		if successmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"success\">%s</p>\n", successmsg)
			fmt.Fprintf(w, "</div>\n")
		}

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label for=\"file\">select file</label>\n")
		fmt.Fprintf(w, "<input id=\"file\" name=\"file\" type=\"file\">\n")
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label for=\"name\">name</label>\n")
		fmt.Fprintf(w, "<input id=\"name\" name=\"name\" type=\"text\" size=\"50\" value=\"%s\">\n", name)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label for=\"path\">path</label>\n")
		fmt.Fprintf(w, "<input id=\"path\" name=\"path\" type=\"text\" size=\"50\" value=\"%s\">\n", path)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">upload file</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func newReplyHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string
		var replybody string

		noteid := idtoi(r.FormValue("noteid"))
		if noteid == -1 {
			log.Printf("new reply: no noteid\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		login := getLoginUser(r, db)
		if login.Userid == -1 {
			log.Printf("new reply: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			replybody = r.FormValue("replybody")
			replybody = strings.ReplaceAll(replybody, "\r", "") // CRLF => CR
			createdt := time.Now().Format(time.RFC3339)

			s := "INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (?, ?, ?, ?)"
			_, err := sqlstmt(db, s).Exec(noteid, replybody, createdt, login.Userid)
			if err != nil {
				log.Printf("DB error creating reply: %s\n", err)
				errmsg = "A problem occured. Please try again."
			}

			if errmsg == "" {
				http.Redirect(w, r, fmt.Sprintf("/note/%d", noteid), http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		// Reply re-entry form
		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/newreply/?noteid=%d\" method=\"post\">\n", noteid)
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>enter reply:</label>\n")
		fmt.Fprintf(w, "<textarea name=\"replybody\" rows=\"10\" cols=\"80\">%s</textarea>\n", replybody)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">add reply</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)

	}
}

func editReplyHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string

		replyid := idtoi(r.FormValue("replyid"))
		if replyid == -1 {
			log.Printf("edit reply: no noteid\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		login := getLoginUser(r, db)
		if login.Userid == -1 {
			log.Printf("edit reply: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var replybody string
		var replyUserid int64
		var noteid int64
		s := "SELECT replybody, user_id, note_id FROM notereply WHERE notereply_id = ?"
		row := db.QueryRow(s, replyid)
		err := row.Scan(&replybody, &replyUserid, &noteid)
		if err == sql.ErrNoRows {
			log.Printf("replyid %d doesn't exist\n", replyid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if err != nil {
			log.Fatal(err)
		}

		// Allow only creators or admin to edit the reply.
		if login.Userid != replyUserid && login.Userid != ADMIN_ID {
			log.Printf("User '%s' doesn't have access to replyid %d\n", login.Username, replyid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			replybody = r.FormValue("replybody")
			replybody = strings.ReplaceAll(replybody, "\r", "") // CRLF => CR

			s := "UPDATE notereply SET replybody = ? WHERE notereply_id = ?"
			_, err = sqlstmt(db, s).Exec(replybody, replyid)
			if err != nil {
				log.Printf("DB error updating replyid %d: %s\n", replyid, err)
				errmsg = "A problem occured. Please try again."
			}

			if errmsg == "" {
				http.Redirect(w, r, fmt.Sprintf("/note/%d", noteid), http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/editreply/?replyid=%d\" method=\"post\">\n", replyid)
		fmt.Fprintf(w, "<h1 class=\"heading\">Edit Reply</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<textarea name=\"replybody\" rows=\"10\" cols=\"80\">%s</textarea>\n", replybody)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">update reply</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		fmt.Fprintf(w, "</article>\n")
		printPageFoot(w)

	}
}

func delReplyHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string

		replyid := idtoi(r.FormValue("replyid"))
		if replyid == -1 {
			log.Printf("del reply: no noteid\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		login := getLoginUser(r, db)
		if login.Userid == -1 {
			log.Printf("del reply: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var replybody string
		var replyUserid int64
		var noteid int64
		s := "SELECT replybody, user_id, note_id FROM notereply WHERE notereply_id = ?"
		row := db.QueryRow(s, replyid)
		err := row.Scan(&replybody, &replyUserid, &noteid)
		if err == sql.ErrNoRows {
			log.Printf("replyid %d doesn't exist\n", replyid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if err != nil {
			log.Fatal(err)
		}

		// Allow only creators or admin to delete the reply.
		if login.Userid != replyUserid && login.Userid != ADMIN_ID {
			log.Printf("User '%s' doesn't have access to replyid %d\n", login.Username, replyid)
			http.Redirect(w, r, fmt.Sprintf("/note/%d", noteid), http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			s := "DELETE FROM notereply WHERE notereply_id = ?"
			_, err = sqlstmt(db, s).Exec(replyid)
			if err != nil {
				log.Printf("DB error deleting replyid %d: %s\n", replyid, err)
				errmsg = "A problem occured. Please try again."
			}
			if errmsg == "" {
				http.Redirect(w, r, fmt.Sprintf("/note/%d", noteid), http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform displayonly\" action=\"/delreply/?replyid=%d\" method=\"post\">\n", replyid)
		fmt.Fprintf(w, "<h1 class=\"heading warning\">Delete Reply</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<textarea name=\"replybody\" rows=\"10\" cols=\"80\" readonly>%s</textarea>\n", replybody)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit warning\">delete reply</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func logoutHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := http.Cookie{
			Name:   "userid",
			Value:  "",
			Path:   "/",
			MaxAge: 0,
		}
		http.SetCookie(w, &c)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func hashPassword(pwd string) string {
	hashedpwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedpwd)
}

func isCorrectPassword(inputPassword, hashedpwd string) bool {
	if hashedpwd == "" && inputPassword == "" {
		return true
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedpwd), []byte(inputPassword))
	if err != nil {
		return false
	}
	return true
}

func loginHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string

		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")

			s := "SELECT user_id, password FROM user WHERE username = ?"
			row := db.QueryRow(s, username, password)

			var userid int
			var hashedpwd string
			err := row.Scan(&userid, &hashedpwd)

			for {
				if err == sql.ErrNoRows {
					errmsg = "Incorrect username or password"
					break
				}
				if err != nil {
					errmsg = "A problem occured. Please try again."
					break
				}
				if !isCorrectPassword(password, hashedpwd) {
					errmsg = "Incorrect username or password"
					break
				}

				suserid := strconv.Itoa(userid)
				c := http.Cookie{
					Name:  "userid",
					Value: suserid,
					Path:  "/",
					// Expires: time.Now().Add(24 * time.Hour),
				}
				http.SetCookie(w, &c)

				// Display notes list page.
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/login/\" method=\"post\">\n")
		fmt.Fprintf(w, "<h1 class=\"heading\">Log In</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>username</label>\n")
		fmt.Fprintf(w, "<input name=\"username\" type=\"text\" size=\"20\">\n")
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>password</label>\n")
		fmt.Fprintf(w, "<input name=\"password\" type=\"password\" size=\"20\">\n")
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">login</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func adminsetupHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getLoginUser(r, db)
		if login.Userid != ADMIN_ID {
			log.Printf("adminsetup: admin not logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<h1 class=\"heading\">Admin Setup</h1>\n")

		fmt.Fprintf(w, "<ul class=\"vertical-list\">\n")
		fmt.Fprintf(w, "  <li>\n")
		fmt.Fprintf(w, "    <p><a href=\"/sitesettings/\">Site Settings</a></p>\n")
		fmt.Fprintf(w, "    <p class=\"finetext\">Set site title and description.</p>\n")
		fmt.Fprintf(w, "  </li>\n")
		fmt.Fprintf(w, "  <li>\n")
		fmt.Fprintf(w, "    <p><a href=\"/userssetup/\">Users</a></p>\n")
		fmt.Fprintf(w, "    <p class=\"finetext\">Set usernames and passwords.</p>\n")
		fmt.Fprintf(w, "  </li>\n")
		fmt.Fprintf(w, "</ul>\n")

		printPageFoot(w)
	}
}

func usersettingsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getLoginUser(r, db)
		if login.Userid == -1 {
			log.Printf("usersettings: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<h2 class=\"heading doc-title\">User</h2>\n")
		fmt.Fprintf(w, "<ul class=\"vertical-list\">\n")

		fmt.Fprintf(w, "<li>\n")
		fmt.Fprintf(w, "<p>%s</p>\n", login.Username)

		fmt.Fprintf(w, "<ul class=\"line-menu finetext\">\n")
		fmt.Fprintf(w, "  <li><a href=\"/edituser?userid=%d\">rename</a>\n", login.Userid)
		fmt.Fprintf(w, "  <li><a href=\"/edituser?userid=%d&setpwd=1\">set password</a>\n", login.Userid)
		fmt.Fprintf(w, "  <li><a href=\"/deactivateuser?userid=%d&setpwd=1\">deactivate</a>\n", login.Userid)
		fmt.Fprintf(w, "</ul>\n")
		fmt.Fprintf(w, "</li>\n")

		fmt.Fprintf(w, "</ul>\n")
		printPageFoot(w)
	}
}

func newUserHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string
		var username, password, password2 string

		login := getLoginUser(r, db)
		if login.Userid != ADMIN_ID {
			log.Printf("new user: admin not logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			username = r.FormValue("username")
			password = r.FormValue("password")
			password2 = r.FormValue("password2")

			for {
				if password != password2 {
					errmsg = "re-entered password doesn't match"
					password = ""
					password2 = ""
					break
				}
				if isUsernameExists(db, username) {
					errmsg = fmt.Sprintf("username '%s' already exists", username)
					break
				}

				hashedPassword := hashPassword(password)
				s := "INSERT INTO user (username, password) VALUES (?, ?);"
				_, err := sqlstmt(db, s).Exec(username, hashedPassword)
				if err != nil {
					log.Printf("DB error creating user: %s\n", err)
					errmsg = "A problem occured. Please try again."
					break
				}

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/newuser/\" method=\"post\">\n")
		fmt.Fprintf(w, "<h1 class=\"heading\">Create User</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>username</label>\n")
		fmt.Fprintf(w, "<input name=\"username\" type=\"text\" size=\"20\" value=\"%s\">\n", username)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>password</label>\n")
		fmt.Fprintf(w, "<input name=\"password\" type=\"password\" size=\"30\" value=\"%s\">\n", password)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>re-enter password</label>\n")
		fmt.Fprintf(w, "<input name=\"password2\" type=\"password\" size=\"30\" value=\"%s\">\n", password2)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">add user</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func editUserHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string
		var username, password, password2 string

		setpwd := r.FormValue("setpwd") // ?setpwd=1 to prompt for new password
		userid := idtoi(r.FormValue("userid"))
		if userid == -1 {
			log.Printf("edit user: no userid\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		login := getLoginUser(r, db)
		if login.Userid != ADMIN_ID && login.Userid != userid {
			log.Printf("edit user: admin or self user not logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		s := "SELECT username FROM user WHERE user_id = ?"
		row := db.QueryRow(s, userid)
		err := row.Scan(&username)
		if err == sql.ErrNoRows {
			// user doesn't exist
			log.Printf("userid %d doesn't exist\n", userid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			oldUsername := username
			username = r.FormValue("username")

			for {
				// If username was changed,
				// make sure the new username hasn't been taken yet.
				if username != oldUsername && isUsernameExists(db, username) {
					errmsg = fmt.Sprintf("username '%s' already exists", username)
					break
				}

				var err error
				if setpwd == "" {
					s := "UPDATE user SET username = ? WHERE user_id = ?"
					_, err = sqlstmt(db, s).Exec(username, userid)
				} else {
					// ?setpwd=1 to set new password
					password = r.FormValue("password")
					password2 = r.FormValue("password2")
					if password != password2 {
						errmsg = "re-entered password doesn't match"
						password = ""
						password2 = ""
						break
					}
					hashedPassword := hashPassword(password)
					s := "UPDATE user SET username = ?, password = ? WHERE user_id = ?"
					_, err = sqlstmt(db, s).Exec(username, hashedPassword, userid)
				}
				if err != nil {
					log.Printf("DB error updating user: %s\n", err)
					errmsg = "A problem occured. Please try again."
					break
				}

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/edituser/?userid=%d&setpwd=%s\" method=\"post\">\n", userid, setpwd)
		fmt.Fprintf(w, "<h1 class=\"heading\">Edit User</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>username</label>\n")
		fmt.Fprintf(w, "<input name=\"username\" type=\"text\" size=\"20\" value=\"%s\">\n", username)
		fmt.Fprintf(w, "</div>\n")

		// ?setpwd=1 to set new password
		if setpwd != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<label>password</label>\n")
			fmt.Fprintf(w, "<input name=\"password\" type=\"password\" size=\"30\" value=\"%s\">\n", password)
			fmt.Fprintf(w, "</div>\n")

			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<label>re-enter password</label>\n")
			fmt.Fprintf(w, "<input name=\"password2\" type=\"password\" size=\"30\" value=\"%s\">\n", password2)
			fmt.Fprintf(w, "</div>\n")
		}

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">update user</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func sitesettingsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string
		var title, desc string

		login := getLoginUser(r, db)
		if login.Userid != ADMIN_ID {
			log.Printf("sitesettings: admin not logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		s := "SELECT title, desc FROM site WHERE site_id = 1"
		row := db.QueryRow(s)
		err := row.Scan(&title, &desc)
		if err == sql.ErrNoRows {
			title = "Group Notes"
			desc = "Central repository for notes"
		} else if err != nil {
			// site settings doesn't exist
			log.Printf("error reading site settings for siteid %d\n", 1)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			title = r.FormValue("title")
			desc = r.FormValue("desc")

			for {
				s := "INSERT OR REPLACE INTO site (site_id, title, desc) VALUES (1, ?, ?)"
				_, err = sqlstmt(db, s).Exec(title, desc)
				if err != nil {
					log.Printf("DB error updating site settings: %s\n", err)
					errmsg = "A problem occured. Please try again."
					break
				}

				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<form class=\"simpleform\" action=\"/sitesettings/\" method=\"post\">\n")
		fmt.Fprintf(w, "<h1 class=\"heading\">Site Settings</h1>")
		if errmsg != "" {
			fmt.Fprintf(w, "<div class=\"control\">\n")
			fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
			fmt.Fprintf(w, "</div>\n")
		}
		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>site title</label>\n")
		fmt.Fprintf(w, "<input name=\"title\" type=\"text\" size=\"50\" value=\"%s\">\n", title)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<label>site description</label>\n")
		fmt.Fprintf(w, "<input name=\"desc\" type=\"text\" size=\"50\" value=\"%s\">\n", desc)
		fmt.Fprintf(w, "</div>\n")

		fmt.Fprintf(w, "<div class=\"control\">\n")
		fmt.Fprintf(w, "<button class=\"submit\">update settings</button>\n")
		fmt.Fprintf(w, "</div>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func userssetupHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string

		login := getLoginUser(r, db)
		if login.Userid != ADMIN_ID {
			log.Printf("userssetup: admin not logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<h1 class=\"heading\">Users Setup</h1>\n")

		fmt.Fprintf(w, "<ul class=\"vertical-list\">\n")

		fmt.Fprintf(w, "<li>\n")
		fmt.Fprintf(w, "  <ul class=\"line-menu finetext\">\n")
		fmt.Fprintf(w, "    <li><p><a href=\"/newuser/\">Create new user</a></p></li>\n")
		fmt.Fprintf(w, "  </ul>\n")
		//		fmt.Fprintf(w, "<p><a href=\"/newuser/\">Create new user</a></p>\n")
		fmt.Fprintf(w, "</li>\n")
		s := "SELECT user_id, username FROM user ORDER BY username"
		rows, err := db.Query(s)
		for {
			if err != nil {
				errmsg = "A problem occured while loading users. Please try again."
				fmt.Fprintf(w, "<li>\n")
				fmt.Fprintf(w, "<p class=\"error\">%s</p>\n", errmsg)
				fmt.Fprintf(w, "</li>\n")
				break
			}

			for rows.Next() {
				var u User
				rows.Scan(&u.Userid, &u.Username)
				fmt.Fprintf(w, "<li>\n")
				fmt.Fprintf(w, "<p>%s</p>\n", u.Username)

				fmt.Fprintf(w, "<ul class=\"line-menu finetext\">\n")
				fmt.Fprintf(w, "  <li><a href=\"/edituser?userid=%d\">rename</a>\n", u.Userid)
				fmt.Fprintf(w, "  <li><a href=\"/edituser?userid=%d&setpwd=1\">set password</a>\n", u.Userid)
				fmt.Fprintf(w, "  <li><a href=\"/deactivateuser?userid=%d&setpwd=1\">deactivate</a>\n", u.Userid)
				fmt.Fprintf(w, "</ul>\n")

				fmt.Fprintf(w, "</li>\n")
			}
			break
		}

		fmt.Fprintf(w, "</ul>\n")
		printPageFoot(w)
	}
}

func printByline(w io.Writer, login User, noteid int64, noteUser User, tcreatedt time.Time, nreplies int) {
	createdt := tcreatedt.Format("2 Jan 2006")
	fmt.Fprintf(w, "<ul class=\"line-menu finetext\">\n")
	fmt.Fprintf(w, "<li>%s</li>\n", createdt)
	fmt.Fprintf(w, "<li>%s</li>\n", noteUser.Username)
	if nreplies > 0 {
		if nreplies == 1 {
			fmt.Fprintf(w, "<li>(%d reply)</li>\n", nreplies)
		} else {
			fmt.Fprintf(w, "<li>(%d replies)</li>\n", nreplies)
		}
	}
	if noteUser.Userid == login.Userid || login.Userid == ADMIN_ID {
		fmt.Fprintf(w, "<li><a href=\"/editnote/?noteid=%d\">Edit</a></li>\n", noteid)
		fmt.Fprintf(w, "<li><a href=\"/delnote/?noteid=%d\">Delete</a></li>\n", noteid)
	}
	fmt.Fprintf(w, "</ul>\n")
}

func printPageHead(w io.Writer) {
	fmt.Fprintf(w, "<!DOCTYPE html>\n")
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "<meta charset=\"utf-8\">\n")
	fmt.Fprintf(w, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	fmt.Fprintf(w, "<title>Website</title>\n")
	fmt.Fprintf(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"/static/style.css\">\n")
	fmt.Fprintf(w, "</head>\n")
	fmt.Fprintf(w, "<body>\n")
}

func printPageFoot(w io.Writer) {
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
}

func printPageNav(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	login := getLoginUser(r, db)

	s := "SELECT title, desc FROM site WHERE site_id = 1"
	row := db.QueryRow(s)
	var title, desc string
	err := row.Scan(&title, &desc)
	if err != nil {
		title = "Group Notes"
		desc = "Central repository for notes"
	}

	fmt.Fprintf(w, "<header class=\"masthead\">\n")
	fmt.Fprintf(w, "<nav class=\"navbar\">\n")
	fmt.Fprintf(w, "<div>\n")

	// Menu row 1
	fmt.Fprintf(w, "<div>\n")
	fmt.Fprintf(w, "<h1><a href=\"/\">%s</a></h1>\n", title)
	if login.Userid != -1 {
		fmt.Fprintf(w, "<a class=\"actiontext\" href=\"/createnote/\">create note</a>\n")
	}
	fmt.Fprintf(w, "<a href=\"/\">browse notes</a>\n")

	fmt.Fprintf(w, "<span class=\"finetext\">\n")
	if login.Userid != -1 {
		fmt.Fprintf(w, "<a href=\"/uploadfile/\">upload file</a>\n")
	}
	fmt.Fprintf(w, "<a href=\"/\">browse files</a>\n")
	if login.Userid == ADMIN_ID {
		fmt.Fprintf(w, "<a href=\"/adminsetup\">setup</a>\n")
	} else if login.Userid != -1 {
		fmt.Fprintf(w, "<a href=\"/usersettings\">settings</a>\n")
	}
	fmt.Fprintf(w, "</span>\n")
	fmt.Fprintf(w, "</div>\n")

	fmt.Fprintf(w, "</div>\n")

	// User section
	fmt.Fprintf(w, "<div>\n")
	fmt.Fprintf(w, "<span>%s</span>\n", login.Username)
	if login.Userid != -1 {
		fmt.Fprintf(w, "<a href=\"/logout\">logout</a>\n")
	} else {
		fmt.Fprintf(w, "<a href=\"/login\">login</a>\n")
	}
	fmt.Fprintf(w, "</div>\n")

	fmt.Fprintf(w, "</nav>\n")
	fmt.Fprintf(w, "<p class=\"finetext\">%s</p>\n", desc)
	fmt.Fprintf(w, "</header>\n")
}

func getLoginUser(r *http.Request, db *sql.DB) User {
	c, err := r.Cookie("userid")
	if err != nil {
		return User{-1, ""}
	}

	userid := idtoi(c.Value)
	if userid == -1 {
		return User{-1, ""}
	}

	var username string
	s := "SELECT username FROM user WHERE user_id = ?"
	row := db.QueryRow(s, userid)
	err = row.Scan(&username)
	if err == sql.ErrNoRows {
		return User{-1, ""}
	}
	return User{Userid: userid, Username: username}
}

func isUsernameExists(db *sql.DB, username string) bool {
	s := "SELECT user_id FROM user WHERE username = ?"
	row := db.QueryRow(s, username)
	var userid int64
	err := row.Scan(&userid)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		return false
	}
	return true
}

func parseMarkdown(s string) string {
	return string(blackfriday.Run([]byte(s), blackfriday.WithExtensions(blackfriday.HardLineBreak)))
}

func idtoi(sid string) int64 {
	if sid == "" {
		return -1
	}
	n, err := strconv.Atoi(sid)
	if err != nil {
		return -1
	}
	return int64(n)
}

func sqlstmt(db *sql.DB, s string) *sql.Stmt {
	stmt, err := db.Prepare(s)
	if err != nil {
		log.Fatalf("db.Prepare() sql: '%s'\nerror: '%s'", s, err)
	}
	return stmt
}
