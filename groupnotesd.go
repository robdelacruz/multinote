package main

import (
	"database/sql"
	"fmt"
	"io"
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

func main() {
	port := "8000"

	db, err := sql.Open("sqlite3", "notes.db")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", notesHandler(db))
	http.HandleFunc("/note/", noteHandler(db))
	http.HandleFunc("/newnote/", newNoteHandler(db))
	http.HandleFunc("/editnote/", editNoteHandler(db))
	http.HandleFunc("/delnote/", delNoteHandler(db))
	http.HandleFunc("/newreply/", newReplyHandler(db))
	http.HandleFunc("/editreply/", editReplyHandler(db))
	http.HandleFunc("/delreply/", delReplyHandler(db))
	http.HandleFunc("/login/", loginHandler(db))
	http.HandleFunc("/logout/", logoutHandler(db))
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
		_, loginUsername := getLoginUser(r, db)

		w.Header().Set("Content-Type", "text/html")

		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<ul class=\"links\">\n")
		s := `SELECT note_id, title, body, createdt, username, 
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
			var title, body, createdt, noteUsername, maxreplydt string
			var numreplies int
			rows.Scan(&noteid, &title, &body, &createdt, &noteUsername, &numreplies, &maxreplydt)
			tcreatedt, _ := time.Parse(time.RFC3339, createdt)

			fmt.Fprintf(w, "<li>\n")
			fmt.Fprintf(w, "<a class=\"note-title\" href=\"/note/%d\">%s</a>\n", noteid, title)

			printByline(w, loginUsername, noteid, noteUsername, tcreatedt, numreplies)
			fmt.Fprintf(w, "</li>\n")
		}
		fmt.Fprintf(w, "</ul>\n")

		fmt.Fprintf(w, "<p class=\"pager-links\">\n")
		fmt.Fprintf(w, "<a href=\"#\">more</a>\n")
		fmt.Fprintf(w, "</p>\n")

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

		loginUserid, loginUsername := getLoginUser(r, db)

		s := `SELECT title, body, createdt, username, (SELECT COUNT(*) FROM notereply WHERE note.note_id = notereply.note_id) AS numreplies 
FROM note
LEFT OUTER JOIN user ON note.user_id = user.user_id
WHERE note_id = ?
ORDER BY createdt DESC;`
		row := db.QueryRow(s, noteid)

		var title, body, createdt, noteUsername string
		var numreplies int
		err := row.Scan(&title, &body, &createdt, &noteUsername, &numreplies)
		if err == sql.ErrNoRows {
			// note doesn't exist so redirect to notes list page.
			log.Printf("display note: noteid %d doesn't exist\n", noteid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<article>\n")
		fmt.Fprintf(w, "<h1><a class=\"note-title\" href=\"/note/%d\">%s</a></h1>", noteid, title)
		tcreatedt, err := time.Parse(time.RFC3339, createdt)
		if err != nil {
			tcreatedt = time.Now()
		}
		printByline(w, loginUsername, noteid, noteUsername, tcreatedt, numreplies)

		bodyMarkup := parseMarkdown(body)
		fmt.Fprintf(w, bodyMarkup)

		fmt.Fprintf(w, "<div class=\"replies\">\n")
		fmt.Fprintf(w, "<hr>\n")
		fmt.Fprintf(w, "<p>Replies:</p>\n")

		s = "SELECT notereply_id, replybody, createdt, username FROM notereply LEFT OUTER JOIN user ON notereply.user_id = user.user_id WHERE note_id = ? ORDER BY notereply_id"
		rows, err := db.Query(s, noteid)
		if err != nil {
			fmt.Fprintf(w, "<p class=\"byline\">Error loading replies</p>\n")
			fmt.Fprintf(w, "</article>\n")
			printPageFoot(w)
			log.Fatal(err)
			return
		}
		i := 1
		for rows.Next() {
			var replyid int64
			var replybody, createdt, replyUsername string
			rows.Scan(&replyid, &replybody, &createdt, &replyUsername)
			tcreatedt, _ := time.Parse(time.RFC3339, createdt)
			createdt = tcreatedt.Format("2 Jan 2006")

			fmt.Fprintf(w, "<p class=\"byline\">\n")
			fmt.Fprintf(w, "%d. %s wrote on %s:", i, replyUsername, createdt)
			if replyUsername == loginUsername {
				fmt.Fprintf(w, "<span class=\"actions\">\n")
				fmt.Fprintf(w, "<a href=\"/editreply/?replyid=%d\">Edit</a>\n", replyid)
				fmt.Fprintf(w, "<a href=\"/delreply/?replyid=%d\">Delete</a>\n", replyid)
				fmt.Fprintf(w, "</span>\n")
			}
			fmt.Fprintf(w, "</p>\n")
			replybodyMarkup := parseMarkdown(replybody)
			fmt.Fprintf(w, replybodyMarkup)
			i++
		}
		fmt.Fprintf(w, "</div>\n")

		// New Reply form
		if loginUserid == -1 {
			fmt.Fprintf(w, "<label class=\"byline\"><a href=\"/login/\">Log in</a> to post a reply.</label>")
		} else {
			fmt.Fprintf(w, "<form action=\"/newreply/?noteid=%d\" method=\"post\">\n", noteid)
			fmt.Fprintf(w, "<label class=\"byline\">reply as %s:</label>\n", loginUsername)
			fmt.Fprintf(w, "<textarea name=\"replybody\" rows=\"10\" cols=\"80\"></textarea><br>\n")
			fmt.Fprintf(w, "<button class=\"submit\">add reply</button>\n")
			fmt.Fprintf(w, "</form>\n")
		}

		fmt.Fprintf(w, "</article>\n")
		printPageFoot(w)
	}
}

func newNoteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errmsg string
		var title, body string

		loginUserid, _ := getLoginUser(r, db)
		if loginUserid == -1 {
			log.Printf("new note: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			title = r.FormValue("title")
			body = r.FormValue("body")
			body = strings.ReplaceAll(body, "\r", "") // CRLF => CR
			createdt := time.Now().Format(time.RFC3339)

			s := "INSERT INTO note (title, body, createdt, user_id) VALUES (?, ?, ?, ?);"
			_, err := sqlstmt(db, s).Exec(title, body, createdt, loginUserid)
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

		fmt.Fprintf(w, "<form action=\"/newnote/\" method=\"post\">\n")
		if errmsg != "" {
			fmt.Fprintf(w, "<label class=\"error\">%s</label><br>\n", errmsg)
		}
		fmt.Fprintf(w, "<label class=\"byline\">title</label>\n")
		fmt.Fprintf(w, "<input name=\"title\" type=\"text\" size=\"50\" value=\"%s\"><br>\n", title)
		fmt.Fprintf(w, "<label class=\"byline\">note</label>\n")
		fmt.Fprintf(w, "<textarea name=\"body\" rows=\"25\" cols=\"80\">%s</textarea><br>\n", body)
		fmt.Fprintf(w, "<button class=\"submit\">add note</button>\n")
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

		loginUserid, loginUsername := getLoginUser(r, db)
		if loginUserid == -1 {
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

		// Allow only creators (todo: also admin) to edit the note.
		if noteUserid != loginUserid {
			log.Printf("User '%s' doesn't have access to note %d\n", loginUsername, noteid)
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

		fmt.Fprintf(w, "<form action=\"/editnote/?noteid=%d\" method=\"post\">\n", noteid)
		if errmsg != "" {
			fmt.Fprintf(w, "<label class=\"error\">%s</label><br>\n", errmsg)
		}
		fmt.Fprintf(w, "<label class=\"byline\">title</label>\n")
		fmt.Fprintf(w, "<input name=\"title\" type=\"text\" size=\"50\" value=\"%s\"><br>\n", title)
		fmt.Fprintf(w, "<label class=\"byline\">note</label>\n")
		fmt.Fprintf(w, "<textarea name=\"body\" rows=\"25\" cols=\"80\">%s</textarea><br>\n", body)
		fmt.Fprintf(w, "<button class=\"submit\">update note</button>\n")
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

		loginUserid, loginUsername := getLoginUser(r, db)
		if loginUserid == -1 {
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

		// Allow only creators (todo: also admin) to delete the note.
		if noteUserid != loginUserid {
			log.Printf("User '%s' doesn't have access to note %d\n", loginUsername, noteid)
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

		fmt.Fprintf(w, "<form class=\"delete\" action=\"/delnote/?noteid=%d\" method=\"post\">\n", noteid)
		fmt.Fprintf(w, "<h2>Delete Note</h2>")
		if errmsg != "" {
			fmt.Fprintf(w, "<label class=\"error\">%s</label><br>\n", errmsg)
		}
		fmt.Fprintf(w, "<label class=\"byline\">title</label>\n")
		fmt.Fprintf(w, "<input class=\"readonly\" name=\"title\" type=\"text\" size=\"50\" readonly value=\"%s\"><br>\n", title)
		fmt.Fprintf(w, "<label class=\"byline\">note</label>\n")
		fmt.Fprintf(w, "<textarea class=\"readonly\" name=\"body\" rows=\"25\" cols=\"80\" readonly>%s</textarea><br>\n", body)
		fmt.Fprintf(w, "<button class=\"submit\">delete note</button>\n")
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

		loginUserid, _ := getLoginUser(r, db)
		if loginUserid == -1 {
			log.Printf("new reply: no user logged in\n")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			replybody = r.FormValue("replybody")
			replybody = strings.ReplaceAll(replybody, "\r", "") // CRLF => CR
			createdt := time.Now().Format(time.RFC3339)

			s := "INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (?, ?, ?, ?)"
			_, err := sqlstmt(db, s).Exec(noteid, replybody, createdt, loginUserid)
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
		fmt.Fprintf(w, "<form action=\"/newreply/?noteid=%d\" method=\"post\">\n", noteid)
		if errmsg != "" {
			fmt.Fprintf(w, "<label class=\"error\">%s</label><br>\n", errmsg)
		}
		fmt.Fprintf(w, "<label class=\"byline\">enter reply:</label>\n")
		fmt.Fprintf(w, "<textarea name=\"replybody\" rows=\"10\" cols=\"80\">%s</textarea><br>\n", replybody)
		fmt.Fprintf(w, "<button class=\"submit\">add reply</button>\n")
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

		loginUserid, loginUsername := getLoginUser(r, db)
		if loginUserid == -1 {
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

		// Only reply creator (todo: also admin) can edit the reply
		if loginUserid != replyUserid {
			log.Printf("User '%s' doesn't have access to replyid %d\n", loginUsername, replyid)
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

		fmt.Fprintf(w, "<form action=\"/editreply/?replyid=%d\" method=\"post\">\n", replyid)
		if errmsg != "" {
			fmt.Fprintf(w, "<label class=\"error\">%s</label><br>\n", errmsg)
		}
		fmt.Fprintf(w, "<label class=\"byline\">edit reply:</label>\n")
		fmt.Fprintf(w, "<textarea name=\"replybody\" rows=\"10\" cols=\"80\">%s</textarea><br>\n", replybody)
		fmt.Fprintf(w, "<button class=\"submit\">update reply</button>\n")
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

		loginUserid, loginUsername := getLoginUser(r, db)
		if loginUserid == -1 {
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

		// Only reply creator (todo: also admin) can delete the reply
		if loginUserid != replyUserid {
			log.Printf("User '%s' doesn't have access to replyid %d\n", loginUsername, replyid)
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

		fmt.Fprintf(w, "<form class=\"delete\" action=\"/delreply/?replyid=%d\" method=\"post\">\n", replyid)
		fmt.Fprintf(w, "<h2>Delete Reply</h2>")
		if errmsg != "" {
			fmt.Fprintf(w, "<label class=\"error\">%s</label><br>\n", errmsg)
		}
		fmt.Fprintf(w, "<textarea class=\"readonly\" name=\"replybody\" rows=\"10\" cols=\"80\">%s</textarea><br>\n", replybody)
		fmt.Fprintf(w, "<button class=\"submit\">delete reply</button>\n")
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

		fmt.Fprintf(w, "<form action=\"/login/\" method=\"post\">\n")
		if errmsg != "" {
			fmt.Fprintf(w, "<label class=\"error\">%s</label><br>\n", errmsg)
		}
		fmt.Fprintf(w, "<label class=\"byline\">username</label>\n")
		fmt.Fprintf(w, "<input name=\"username\" type=\"text\" size=\"20\"><br>\n")
		fmt.Fprintf(w, "<label class=\"byline\">password</label>\n")
		fmt.Fprintf(w, "<input name=\"password\" type=\"password\" size=\"20\"><br>\n")
		fmt.Fprintf(w, "<button class=\"submit\">login</button>\n")
		fmt.Fprintf(w, "</form>\n")

		printPageFoot(w)
	}
}

func printByline(w io.Writer, loginUsername string, noteid int64, noteUsername string, tcreatedt time.Time, nreplies int) {
	createdt := tcreatedt.Format("2 Jan 2006")
	fmt.Fprintf(w, "<p class=\"byline\">\n")
	fmt.Fprintf(w, "posted by %s on <time>%s</time> (%d replies)", noteUsername, createdt, nreplies)
	if noteUsername == loginUsername {
		fmt.Fprintf(w, "<span class=\"actions\">\n")
		fmt.Fprintf(w, "<a href=\"/editnote/?noteid=%d\">Edit</a>\n", noteid)
		fmt.Fprintf(w, "<a href=\"/delnote/?noteid=%d\">Delete</a>\n", noteid)
		fmt.Fprintf(w, "</span>\n")
	}
	fmt.Fprintf(w, "</p>\n")
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
	loginUserid, loginUsername := getLoginUser(r, db)

	fmt.Fprintf(w, "<nav>\n")

	fmt.Fprintf(w, "<div>\n")
	fmt.Fprintf(w, "<h1><a href=\"/\">Group Notes</a></h1>\n")
	fmt.Fprintf(w, "<a href=\"/\">latest</a>\n")
	if loginUserid != -1 {
		fmt.Fprintf(w, "<a href=\"/newnote/\">new note</a>\n")
	}
	fmt.Fprintf(w, "<p class=\"byline\">I reserve the right to be biased, it makes life more interesting.</p>\n")
	fmt.Fprintf(w, "</div>\n")

	fmt.Fprintf(w, "<div>\n")
	fmt.Fprintf(w, "<span>%s</span>\n", loginUsername)
	if loginUserid != -1 {
		fmt.Fprintf(w, "<a href=\"/logout\">logout</a>\n")
	} else {
		fmt.Fprintf(w, "<a href=\"/login\">login</a>\n")
	}
	fmt.Fprintf(w, "</div>\n")

	fmt.Fprintf(w, "</nav>\n")
}

func getLoginUser(r *http.Request, db *sql.DB) (int64, string) {
	c, err := r.Cookie("userid")
	if err != nil {
		return -1, ""
	}

	userid := idtoi(c.Value)
	if userid == -1 {
		return -1, ""
	}

	var username string
	s := "SELECT username FROM user WHERE user_id = ?"
	row := db.QueryRow(s, userid)
	err = row.Scan(&username)
	if err == sql.ErrNoRows {
		return -1, ""
	}
	return userid, username
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
