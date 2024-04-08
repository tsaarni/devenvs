CC = gcc
TARGET = pam_userenv.so
INSTALL_DIR = /usr/lib/x86_64-linux-gnu/security/
CFLAGS = -Wall -Werror -Wextra
LDFLAGS = -fPIC -shared

all: $(TARGET)

clean:
	rm -f $(TARGET)

$(TARGET): userenv.c
	$(CC) $(CFLAGS) $(LDFLAGS) -o $(TARGET) $<

install: $(TARGET)
	install -D -m 755 $(TARGET) $(DESTDIR)$(INSTALL_DIR)/$(TARGET)
