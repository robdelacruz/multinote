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
	"gopkg.in/russross/blackfriday.v2"
)

func main() {
	port := "8000"

	db, err := sql.Open("sqlite3", "notes.db")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", notesHandler(db))
	http.HandleFunc("/new/", newHandler(db))
	http.HandleFunc("/login/", loginHandler(db))
	http.HandleFunc("/logout/", logoutHandler(db))
	http.HandleFunc("/note/", noteHandler(db))
	fmt.Printf("Listening on %s...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}

func notesHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, "<ul class=\"links\">\n")

		s := `SELECT note_id, title, body, createdt, username, 
(SELECT COUNT(*) FROM notereply WHERE note.note_id = notereply.note_id) AS numreplies, 
(SELECT MAX(createdt) FROM notereply where note.note_id = notereply.note_id) AS maxreplydt 
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
			var title, body, createdt, username, maxreplydt string
			var numreplies int
			rows.Scan(&noteid, &title, &body, &createdt, &username, &numreplies, &maxreplydt)
			tcreatedt, _ := time.Parse(time.RFC3339, createdt)

			fmt.Fprintf(w, "<li>\n")
			fmt.Fprintf(w, "<a class=\"note-title\" href=\"/note/%d\">%s</a>\n", noteid, title)

			printByline(w, username, tcreatedt, numreplies)
			fmt.Fprintf(w, "</li>\n")
		}

		fmt.Fprintf(w, `</ul>
<p class="pager-links">
`)
		fmt.Fprintf(w, `<a href="#">more</a>
</p>
`)

		printPageFoot(w)
	}
}

func parseNoteid(url string) int64 {
	sre := `^/note/(\d+)$`
	re := regexp.MustCompile(sre)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return -1
	}
	noteid, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1
	}
	return int64(noteid)
}

func noteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		noteid := parseNoteid(r.URL.Path)
		if noteid == -1 {
			log.Printf("no noteid parameter in '%s'\n", r.URL.Path)
			// no note id parameter, so redirect to notes list page.
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			replybody := r.FormValue("replybody")
			createdt := time.Now().Format(time.RFC3339)

			userid, _ := getLoginUser(r, db)
			if userid == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			// Strip out linefeed chars so that CRLF becomes just CR.
			// CRLF causes problems in markdown parsing.
			replybody = strings.ReplaceAll(replybody, "\r", "")

			s := "INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (?, ?, ?, ?)"
			stmt, err := db.Prepare(s)
			if err != nil {
				log.Fatal(err)
			}
			_, err = stmt.Exec(noteid, replybody, createdt, userid)
			if err != nil {
				log.Fatal(err)
			}

			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}

		s := `SELECT title, body, createdt, username , (SELECT COUNT(*) FROM notereply WHERE note.note_id = notereply.note_id) AS numreplies 
FROM note
LEFT OUTER JOIN user ON note.user_id = user.user_id
WHERE note_id = ?
ORDER BY createdt DESC;`
		row := db.QueryRow(s, noteid)

		var title, body, createdt, username string
		var numreplies int
		err := row.Scan(&title, &body, &createdt, &username, &numreplies)
		if err == sql.ErrNoRows {
			// note doesn't exist so redirect to notes list page.
			log.Printf("noteid %d doesn't exist\n", noteid)
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
		printByline(w, username, tcreatedt, numreplies)

		bodyMarkup := parseMarkdown(body)
		fmt.Fprintf(w, string(bodyMarkup))
		fmt.Fprintf(w, "</article>\n")

		fmt.Fprintf(w, `<article class="replies">
<hr>
<p>Replies:</p>
`)

		s = "SELECT replybody, createdt, username FROM notereply LEFT OUTER JOIN user ON notereply.user_id = user.user_id WHERE note_id = ? ORDER BY notereply_id"
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
			var replybody, createdt, username string
			rows.Scan(&replybody, &createdt, &username)
			tcreatedt, _ := time.Parse(time.RFC3339, createdt)
			createdt = tcreatedt.Format("2 Jan 2006")

			fmt.Fprintf(w, "<p class=\"byline\">%d. %s wrote on %s:</p>", i, username, createdt)
			replybodyMarkup := parseMarkdown(replybody)
			fmt.Fprintf(w, string(replybodyMarkup))
			i++
		}

		fmt.Fprintf(w, `<form method="post">
        <label class="byline">post comment:</label>
        <textarea name="replybody" rows="10" cols="80"></textarea><br>
        <button class="submit">add reply</button>
    </form>
`)
		fmt.Fprintf(w, "</article>\n")

		printPageFoot(w)
	}
}

func newHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			title := r.FormValue("title")
			body := r.FormValue("body")
			createdt := time.Now().Format(time.RFC3339)

			userid, _ := getLoginUser(r, db)
			if userid == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			// Strip out linefeed chars so that CRLF becomes just CR.
			// CRLF causes problems in markdown parsing.
			body = strings.ReplaceAll(body, "\r", "")

			s := "INSERT INTO note (title, body, createdt, user_id) VALUES (?, ?, ?, ?);"
			stmt, err := db.Prepare(s)
			if err != nil {
				log.Fatal(err)
			}
			_, err = stmt.Exec(title, body, createdt, userid)
			if err != nil {
				log.Fatal(err)
			}

			// Display notes list page.
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		printPageHead(w)
		printPageNav(w, r, db)

		fmt.Fprintf(w, `
<form action="/new/" method="post">
    <label class="byline">title</label>
    <input name="title" type="text" size="50"><br>
    <label class="byline">note</label>
    <textarea name="body" rows="25" cols="80"></textarea><br>
    <button class="submit">add note</button>
</form>
`)

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

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func loginHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errMarkup := ""

		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")

			s := "SELECT user_id FROM user WHERE username = ? AND password = ?"
			row := db.QueryRow(s, username, password)

			var userid int
			err := row.Scan(&userid)
			if err == sql.ErrNoRows {
				errMarkup = "<p class=\"byline\">Incorrect username or password</p>\n"
			} else if err != nil {
				errMarkup = "<p class=\"byline\">Server error during login</p>\n"
			} else {
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

		if errMarkup != "" {
			fmt.Fprintf(w, errMarkup)
		}

		fmt.Fprintf(w, `
<form action="/login/" method="post">
    <label class="byline">username</label>
    <input name="username" type="text" size="20"><br>
    <label class="byline">password</label>
    <input name="password" type="password" size="20"><br>
    <button class="submit">login</button>
</form>
`)

		printPageFoot(w)
	}
}

func printByline(w io.Writer, username string, tcreatedt time.Time, nreplies int) {
	createdt := tcreatedt.Format("2 Jan 2006")
	fmt.Fprintf(w, "<p class=\"byline\">posted by %s on <time>%s</time> (%d replies)</p>", username, createdt, nreplies)
}

func printPageHead(w io.Writer) {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Website</title>
<link rel="stylesheet" type="text/css" href="/static/style.css">
</head>
<body>
`)
}

func printPageFoot(w io.Writer) {
	fmt.Fprintf(w, `</body>
</html>
`)
}

func printPageNav(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, username := getLoginUser(r, db)

	fmt.Fprintf(w, "<nav>\n")

	fmt.Fprintf(w, "<div>\n")
	fmt.Fprintf(w, "<h1><a href=\"/\">Group Notes</a></h1>\n")
	fmt.Fprintf(w, "<a href=\"/\">latest</a>\n")
	if username != "" {
		fmt.Fprintf(w, "<a href=\"/new\">new note</a>\n")
	}
	fmt.Fprintf(w, "<p class=\"byline\">I reserve the right to be biased, it makes life more interesting.</p>\n")
	fmt.Fprintf(w, "</div>\n")

	fmt.Fprintf(w, "<div>\n")
	if username != "" {
		fmt.Fprintf(w, "<span>%s</span>\n", username)
		fmt.Fprintf(w, "<a href=\"/logout\">logout</a>\n")
	} else {
		fmt.Fprintf(w, "<span></span>\n")
		fmt.Fprintf(w, "<a href=\"/login\">login</a>\n")
	}
	fmt.Fprintf(w, "</div>\n")

	fmt.Fprintf(w, "</nav>\n")
}

func getLoginUser(r *http.Request, db *sql.DB) (string, string) {
	userid := ""
	username := ""
	c, err := r.Cookie("userid")
	if err == nil {
		userid = c.Value
		s := "SELECT username FROM user WHERE user_id = ?"
		row := db.QueryRow(s, userid)
		row.Scan(&username)
	}

	return userid, username
}

func parseMarkdown(s string) string {
	return string(blackfriday.Run([]byte(s), blackfriday.WithExtensions(blackfriday.HardLineBreak)))
}
