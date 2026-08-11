SHELL=/bin/bash
GOBUILD=go build
GORUN=go run
MAIN_PROG=listdbg2
LOAD_PROG=loaddb
SRCS=main.go
LOAD_SRCS=loaddb.go

# ---------------------------------------------------
clean:
	-@rm $(MAIN_PROG)

# ---------------------------------------------------
gofmt:
	@gofmt -l -s -w .

# ---------------------------------------------------
run:
	$(GORUN) $(SRCS)

# ---------------------------------------------------
load:
	$(GORUN) $(LOAD_SRCS)

# ---------------------------------------------------
build:
	$(GOBUILD) -o $(MAIN_PROG) $(SRCS)

# ---------------------------------------------------
build-win64:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(MAIN_PROG).exe $(SRCS)

# ---------------------------------------------------
#build-win32:
#	GOOS=windows GOARCH=386 $(GOBUILD) -o $(MAIN_PROG)32.exe $(SRCS)

# ---------------------------------------------------
build-arm6:
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(MAIN_PROG)_arm6 $(SRCS)
	#GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(MAIN_PROG) $(SRCS)

# ---------------------------------------------------
#build-arm7:
#	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(MAIN_PROG)_arm7 $(SRCS)

# ---------------------------------------------------
#build-arm64:
#	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(MAIN_PROG)_arm64 $(SRCS)
