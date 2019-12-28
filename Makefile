SHAREDIR=/usr/local/share/spot
BINDIR=/usr/local/bin

all: spotd spot

dep:
	go get -u github.com/mattn/go-sqlite3

spotd: spotd.go db.go
	go build -o spotd spotd.go db.go

spot: spot.go db.go
	go build -o spot spot.go db.go

clean:
	rm -rf spotd spot

install: spotd
	mkdir -p $(SHAREDIR)
	touch $(SHAREDIR)/spot.db
	chmod a+w $(SHAREDIR)
	chmod a+w $(SHAREDIR)/spot.db
	cp spotd $(BINDIR)
	cp spot $(BINDIR)

uninstall:
	rm -rf $(SHAREDIR)
	rm -rf $(BINDIR)/spotd $(BINDIR)/spot

