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
	http.HandleFunc("/note/", noteHandler(db))
	fmt.Printf("Listening on %s...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}

func notesHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		printPageHead(w)
		printPageNav(w)

		fmt.Fprintf(w, "<ul class=\"links\">\n")

		s := "SELECT note_id, title, body, createdt, username FROM note LEFT OUTER JOIN user ON note.user_id = user.user_id ORDER BY createdt DESC;"
		rows, err := db.Query(s)
		if err != nil {
			return
		}
		for rows.Next() {
			var noteid int64
			var title, body, createdt, username string
			rows.Scan(&noteid, &title, &body, &createdt, &username)
			tcreatedt, err := time.Parse(time.RFC3339, createdt)
			if err != nil {
				tcreatedt = time.Now()
			}

			fmt.Fprintf(w, "<li>\n")
			fmt.Fprintf(w, "<a href=\"/note/%d\">%s</a>\n", noteid, title)
			printByline(w, username, tcreatedt)
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

func noteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sre := `^/note/(\d+)$`
		re := regexp.MustCompile(sre)
		matches := re.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			log.Printf("no noteid parameter in '%s'\n", r.URL.Path)
			// no note id parameter, so redirect to notes list page.
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		noteid, _ := strconv.Atoi(matches[1])

		s := "SELECT title, body, createdt, username FROM note LEFT OUTER JOIN user ON note.user_id = user.user_id WHERE note_id = ? ORDER BY createdt DESC;"
		row := db.QueryRow(s, noteid)

		var title, body, createdt, username string
		err := row.Scan(&title, &body, &createdt, &username)
		if err == sql.ErrNoRows {
			// note doesn't exist so redirect to notes list page.
			log.Printf("noteid %d doesn't exist\n", noteid)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		printPageHead(w)
		printPageNav(w)

		fmt.Fprintf(w, "<article>\n")
		fmt.Fprintf(w, "<h1><a href=\"/note/%d\">%s</a></h1>", noteid, title)
		tcreatedt, err := time.Parse(time.RFC3339, createdt)
		if err != nil {
			tcreatedt = time.Now()
		}
		printByline(w, username, tcreatedt)

		//body = strings.ReplaceAll(body, "\r", "")
		bodyMarkup := blackfriday.Run([]byte(body))
		fmt.Fprintf(w, string(bodyMarkup))
		fmt.Fprintf(w, "</article>\n")

		fmt.Fprintf(w, `<article class="replies">
<hr>
<p>Replies:</p>
`)

		fmt.Fprintf(w, "<p class=\"byline\">1. robdelacruz wrote on 2019-12-06:</p>")
		fmt.Fprintf(w, "<p>comment text here</p>")

		fmt.Fprintf(w, `<form>
        <p class="byline">post comment:</p>
        <textarea rows="10" cols="80"></textarea><br>
        <button>add reply</button>
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
			userid := "robdelacruz"

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
		printPageNav(w)

		fmt.Fprintf(w, `
<form action="/new/" method="post">
    <p class="byline">title</p>
    <input name="title" type="text" size="50"><br>
    <p class="byline">note</p>
    <textarea name="body" rows="25" cols="80"></textarea><br>
    <button>add note</button>
</form>
`)

		printPageFoot(w)
	}
}

func printByline(w io.Writer, username string, tcreatedt time.Time) {
	createdt := tcreatedt.Format("2 Jan 2006")
	fmt.Fprintf(w, "<p class=\"byline\">posted by %s on <time>%s</time> (2 replies)</p>", username, createdt)
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

func printPageNav(w io.Writer) {
	fmt.Fprintf(w, `<nav>
    <h1><a href="/">Group Notes</a></h1>
    <a href="/">latest</a>
    <a href="/new">new note</a>
    <p class="byline">I reserve the right to be biased, it makes life more interesting.</p>
</nav>
`)
}
