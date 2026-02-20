CC      ?= cc
CFLAGS  ?= -O2 -g -Wall
BINDIR  ?= /usr/local/libexec/snapstore

exiftrim: src/main.c
	$(CC) $(CFLAGS) -o $@ src/main.c

install: exiftrim
	install -d $(BINDIR)
	install -m 0755 exiftrim $(BINDIR)/exiftrim

clean:
	rm -f exiftrim

.PHONY: install clean
