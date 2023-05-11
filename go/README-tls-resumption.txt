

crypto/tls: large handshake records may cause panics (CVE-2022-41724)
https://github.com/golang/go/issues/58001

https://go-review.googlesource.com/c/go/+/468117

https://github.com/golang/go/issues/58358   # Backport: Go1.19.6 ->
https://github.com/golang/go/issues/58359   # Backport: Go1.20.1 ->




cd go/src
./all.bash      # compile (you can press ctrl+c when testing begins)


PATH=$PWD/../bin/:$PATH go test -v ./crypto/tls/ -run ^TestResumption



# run only tls tests
PATH=$PWD/../bin/:$PATH go tool dist test -rebuild -run crypto/tls



# run test app
go run reproduce-tls-resumption-failure.go






 ~/work/go/bin/go run server/server.go
Request from: 127.0.0.1:58238
2023/03/29 20:34:35 http: TLS handshake error from 127.0.0.1:50606: tls: invalid PSK binder


# Test that reproduces the failure 
#   src/crypto/tls/handshake_client_test.go

// Trigger resumption by forcing CurveP521
func TestResumptionInvalidPSKBinder(t *testing.T) {
	serverConfig := &Config{
		CurvePreferences: []CurveID{CurveP521, CurveP384, CurveP256},
		Certificates: []Certificate{
			{
				Certificate: [][]byte{testP256Certificate},
				PrivateKey:  testP256PrivateKey,
			},
		},
	}

	clientConfig := &Config{
		ClientSessionCache: NewLRUClientSessionCache(32),
		ServerName:         "example.golang",
		InsecureSkipVerify: true,
	}

	testResumeState := func(test string, didResume bool) {
		_, hs, err := testHandshake(t, clientConfig, serverConfig)
		if err != nil {
			t.Fatalf("%s: handshake failed: %s", test, err)
		}
		if hs.DidResume != didResume {
			t.Fatalf("%s resumed: %v, expected: %v", test, hs.DidResume, didResume)
		}
		if got, want := hs.ServerName, clientConfig.ServerName; got != want {
			t.Errorf("%s: server name %s, want %s", test, got, want)
		}
	}

	testResumeState("Handshake", false)
	testResumeState("Resume", true)
}



Processing of TLSv1.3 Hello Retry Request can fail


If client sends Client Hello with mismatching ECDHE groups, the server sends Hello Retry Request.

From RFC 8446

2.1.  Incorrect DHE Share


   If the client has not provided a sufficient "key_share" extension
   (e.g., it includes only DHE or ECDHE groups unacceptable to or
   unsupported by the server), the server corrects the mismatch with a
   HelloRetryRequest and the client needs to restart the handshake with
   an appropriate "key_share" extension

the logic in processHelloRetryRequest calculates wrong hash to PSK binder entry,
in the retried Client Hello message, leading to server aborting the handshake
with "invalid PSK binder" alert





https://www.rfc-editor.org/rfc/rfc8446.html
https://www.rfc-editor.org/rfc/rfc8446.html#section-4.2.11.2    # PSK binder

TLS slides
https://mycourses.aalto.fi/pluginfile.php/1142733/mod_resource/content/1/Network%20Security%20C3%20-%20TLS%20PSK.pdf


