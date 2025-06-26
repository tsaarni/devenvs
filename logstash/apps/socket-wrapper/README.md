# Example: Socket Wrapper Library for Setting DSCP

## Description

This is a simple example of an `LD_PRELOAD`-based socket wrapper library that allows setting the DSCP (Differentiated Services Code Point) value for outgoing packets.

## Usage

### Compilation

Compile the library using the following command:

```bash
make
```

This will generate the shared library `libsocket_wrapper.so`.

### Using the Library

To use the library, set the `LD_PRELOAD` environment variable to point to the compiled library. This will intercept socket-related system calls to set the DSCP value.

#### Example with Python HTTP Server and Java Client

Run an HTTP server:

```bash
python3 -m http.server
```

Run an HTTP client test application with the library:

```bash
LD_PRELOAD=./libsocket_wrapper.so java TestApp
```

#### Example with C-Based Applications

Compile your C application:

```bash
gcc -o testapp testapp.c
```

Run the application with the library:

```bash
LD_PRELOAD=./libsocket_wrapper.so ./testapp
```

#### Example with JRuby Client

Run an HTTP client test application with the library and JRuby:

```bash
LD_PRELOAD=./libsocket_wrapper.so jruby testapp.rb
```


### Verifying DSCP Value

To verify the DSCP value, use Wireshark to capture packets.
Start Wireshark with the following command:

```bash
wireshark -i lo -k -f "port 8000"
```

Then, run your application with and without the `LD_PRELOAD` environment variable to observe the difference in DSCP values.

### Test with non-glibc environment

```
docker run -it --rm --volume $(pwd):/data --network host alpine:latest ash
```

then run the following commands inside the container:

```
cd /data
apk add --no-cache gcc musl-dev make
make
LD_PRELOAD=./libsocket_wrapper.so ./testapp
```

### Notes

- The example code wraps both `bind()` and `connect()` functions so DSCP will be set on both server and client sockets.
  - To limit the impact to only client sockets, you can modify the code to only wrap `connect()`.
  - To limit the impact further, in theory it could be possible to e.g. check the target address range and set DSCP only for certain ranges but it would require prior knowledge on the target address range - or some other criteria.
- The DSCP value is hardcoded to `46` in the example code.
- This method is compatible only with application runtimes that utilize the standard socket API calls provided by the C library.
For instance, it will not work with Go, as Go does not depend on the C library.
- In case of ld problems, use LD_DEBUG=all
