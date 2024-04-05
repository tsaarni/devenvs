SRC = fake.c
OUT = libnss_fake.so.2
SONAME = libnss_fake.so.2
INSTALL_DIR = /usr/lib/x86_64-linux-gnu/
LDFLAGS = -fPIC -shared -Wl,-soname,$(SONAME)

all: $(OUT)

$(OUT): fake.c
	gcc $(LDFLAGS) -o $(OUT) $<

install: $(OUT)
	install -D -m 755 $(OUT) $(DESTDIR)$(INSTALL_DIR)/$(OUT)
