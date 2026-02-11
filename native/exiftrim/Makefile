CC      ?= cc
CFLAGS  ?= -O2 -g -Wall -Iinclude
BINDIR  ?= /usr/local/libexec/snapstore

SRC := src/exiftrim.c src/parse.c
OBJ := $(SRC:.c=.o)

exiftrim: $(OBJ)
	$(CC) $(CFLAGS) -o $@ $(OBJ)

%.o: %.c include/parse.h
	$(CC) $(CFLAGS) -c -o $@ $<

install: exiftrim
	install -d $(BINDIR)
	install -m 0755 exiftrim $(BINDIR)/exiftrim

clean:
	rm -f $(OBJ) exiftrim

.PHONY: install clean
