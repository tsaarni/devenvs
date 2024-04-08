CC = gcc
TARGET = libgetlogin.so
INSTALL_DIR = /usr/lib/
CFLAGS = -Wall -Werror
LDFLAGS = -fPIC -shared

all: $(TARGET)

clean:
	rm -f $(TARGET)

$(TARGET): getlogin.c
	$(CC) $(CFLAGS) $(LDFLAGS) -o $(TARGET) $<

install: $(TARGET)
	install -D -m 755 $(TARGET) $(DESTDIR)$(INSTALL_DIR)/$(TARGET)
