
##  Start Vault in dev mode in debugger

VAULT_TOKEN=root

# store secret in kv


# Invalid JSON body started to fail in Vault 1.20.3

http -v POST localhost:8200/v1/secret/foo  X-Vault-Token:$VAULT_TOKEN --raw '{"foo":"bar"}"'

HTTP/1.1 500 Internal Server Error
Content-Length: 56
Content-Type: application/json
Date: Thu, 25 Sep 2025 11:00:20 GMT

{
   "errors": [ "error reading JSON token: unexpected EOF" ]
}



Beginning from vault 1.20.3 there is a denial-of-service protection for the parser by checking json payload depth (and some other limits)

 https://developer.hashicorp.com/vault/docs/updates/important-changes#json-payload-limits



This was implemented by JSON parser's Token() function by tokenizing the JSON document before unmarshaling it

It is done in:

https://github.com/hashicorp/vault/blob/7665ff29d77e5cb3ea9ddbeaed49ee312e53c6b8/sdk/helper/jsonutil/json.go#L166

// VerifyMaxDepthStreaming scans the JSON stream to enforce nesting depth, counts,
// and other limits without decoding the full structure into memory.
func VerifyMaxDepthStreaming(jsonReader io.Reader, limits JSONLimits) (int, error) {


The code could in theory ignore invalid tokens (err is io.ErrUnexpectedEOF) at the end of the stream, but it does not.

The code was recently changed again to custom tokenizer

https://github.com/hashicorp/vault/commit/b19e74c29a33ed2a99fc01626104db1a49345df3


it does not use JSON decoder anymore for the check.
The error to invalid JSON is now:


HTTP/1.1 500 Internal Server Error
Content-Length: 60
Content-Type: application/json
Date: Thu, 25 Sep 2025 12:36:16 GMT
{
   "errors": [ "invalid character '\"' after top-level value" ]
}
