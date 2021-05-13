.PHONY: all
all: build

.PHONY: FORCE
FORCE:

bin/deproc: cmd/deproc FORCE
	go build -o $@ ./$<

.PHONY: build
build: bin/deproc

.PHONY: clean
clean:
	git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: install
install: bin/deproc
	cp $< ~/bin/