CC = gcc
TARGET = libnss_fake.so.2
INSTALL_DIR = /usr/lib/x86_64-linux-gnu/
CFLAGS = -Wall -Werror
LDFLAGS = -fPIC -shared -Wl,-soname,$(TARGET)

all: $(TARGET)

clean:
	rm -f $(TARGET)

$(TARGET): fake.c
	$(CC) $(CFLAGS) $(LDFLAGS) -o $(TARGET) $<

install: $(TARGET)
	install -D -m 755 $(TARGET) $(DESTDIR)$(INSTALL_DIR)/$(TARGET)
