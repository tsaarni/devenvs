# Example: Socket Wrapper Library for Setting DSCP

## Description

This is a simple example of an `LD_PRELOAD`-based socket wrapper library that allows settubg the DSCP (Differentiated Services Code Point) value for outgoing packets.

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

### Verifying DSCP Value

To verify the DSCP value, use Wireshark to capture packets.
Start Wireshark with the following command:

```bash
wireshark -i lo -k -f "port 8000"
```

Then, run your application with and without the `LD_PRELOAD` environment variable to observe the difference in DSCP values.

### Notes

- This method is compatible only with application runtimes that utilize the standard socket API calls provided by the C library.
For instance, it will not work with Go, as Go does not depend on the C library.
