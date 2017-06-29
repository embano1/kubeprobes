REPO = embano1
APPNAME = probes
VERSION = 2.0
LDFLAGS1 = '-extldflags "-static"'
LDFLAGS2 = "-s -w -X main.version=$(VERSION)"

build:
	GOOS=linux CGO=0 go build -a --ldflags $(LDFLAGS1) --ldflags $(LDFLAGS2) -tags netgo -installsuffix netgo -o $(APPNAME) .

image: build
	docker build --rm -t $(REPO)/$(APPNAME):$(VERSION) .
	rm $(APPNAME)

clean:
	rm $(APPNAME)

all: image
