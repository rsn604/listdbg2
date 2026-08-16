SHELL=/bin/bash
GOBUILD=go build
GORUN=go run
MAIN_PROG=listdbg2
LOAD_PROG=loaddb
MAIN_SRCS=main.go
LOAD_SRCS=loaddb.go

BOLTDB=BOLT
BOLTSTRING=./db/ListDB.boltdb

CSVDIR=./csv

# ---------------------------------------------------
clean:
	-@rm $(MAIN_PROG)

# ---------------------------------------------------
gofmt:
	@gofmt -l -s -w .

# ---------------------------------------------------
run:
	$(GORUN) $(MAIN_SRCS) $(BOLTDB) $(BOLTSTRING)

# ---------------------------------------------------
load:
	$(GORUN) $(LOAD_SRCS) $(BOLTDB) $(BOLTSTRING) $(CSVDIR) 

# ---------------------------------------------------
build:
	$(GOBUILD) -o $(MAIN_PROG) $(MAIN_SRCS)
	$(GOBUILD) -o $(LOAD_PROG) $(LOAD_SRCS)

# ---------------------------------------------------
build-win64:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(MAIN_PROG).exe $(MAIN_SRCS)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(LOAD_PROG).exe $(LOAD_SRCS)

# ---------------------------------------------------
build-arm6:
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(MAIN_PROG)_arm6 $(MAIN_SRCS)
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(LOAD_PROG)_arm6 $(LOAD_SRCS)

