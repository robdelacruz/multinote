package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	port := "8000"

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/new", rootHandler())
	fmt.Printf("Listening on %s...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatal(err)
}

func rootHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		printPageHead(w)
		printPageNav(w)

		fmt.Fprintf(w, `
<form>
    <p class="byline">title</p>
    <input type="text" size="50"><br>
    <p class="byline">note</p>
    <textarea rows="25" cols="80"></textarea><br>
    <button>add note</button>
</form>
`)

		printPageFoot(w)
	}
}

func printPageHead(w io.Writer) {
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Website</title>
<link rel="stylesheet" type="text/css" href="static/style.css">
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
    <h1><a href="#">Group Notes</a></h1>
    <a href="notes.html">latest</a>
    <a href="newnote.html">new note</a>
    <p class="byline">I reserve the right to be biased, it makes life more interesting.</p>
</nav>
`)
}
