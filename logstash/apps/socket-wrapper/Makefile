CC = gcc
TARGET = libsocket_wrapper.so
CFLAGS = -Wall -Werror -fPIC
LDFLAGS ?= -shared -fPIC -ldl

all: $(TARGET)

clean:
	rm -f *.o $(TARGET)

$(TARGET): socket_wrapper.o
	$(CC) $(CFLAGS) $(LDFLAGS) -o $(TARGET) $<
