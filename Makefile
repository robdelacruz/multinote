SHAREDIR=/usr/local/share/groupnotes
BINDIR=/usr/local/bin

all: groupnotesd

dep:
	go get -u github.com/mattn/go-sqlite3
	go get -u gopkg.in/russross/blackfriday.v2

groupnotesd: groupnotesd.go
	go build -o groupnotesd groupnotesd.go

clean:
	rm -rf groupnotesd

install: groupnotesd
	mkdir -p $(SHAREDIR)
	touch $(SHAREDIR)/groupnotes.db
	chmod a+w $(SHAREDIR)
	chmod a+w $(SHAREDIR)/groupnotes.db
	cp groupnotesd $(BINDIR)

uninstall:
	rm -rf $(SHAREDIR)
	rm -rf $(BINDIR)/groupnotesd

