
# Rekey cancel returns 400 EOF error unless nonce is provided
# https://github.com/hashicorp/vault/issues/31649

# http: allow empty request bodies in JSON parsing
# https://github.com/hashicorp/vault/pull/31650


bin/vault server -dev -dev-root-token-id=root

export VAULT_ADDR=http://127.0.0.1:8200

bin/vault login root
bin/vault operator rekey -init -key-shares=3 -key-threshold=2
# wait for 10 minutes
bin/vault operator rekey -cancel





Timeout can be shortened for tests in
http/sys_rekey.go  handleSysRekeyInitDelete()

-	if err := core.RekeyCancel(recovery, req.Nonce, 10*time.Minute); err != nil {
+	if err := core.RekeyCancel(recovery, req.Nonce, 10*time.Second); err != nil {









# Failure when sendining cancel request with empty body (before fix)
Error canceling rekey: Error making API request.

URL: DELETE http://127.0.0.1:8200/v1/sys/rekey/init
Code: 400. Errors:

* EOF



# Error when sending cancel without nonce before 10 minutes (after the fix)

Error canceling rekey: Error making API request.

URL: DELETE http://127.0.0.1:8200/v1/sys/rekey/init
Code: 400. Errors:

* invalid request



# Success when sending cancel with nonce before 10 minutes (after the fix)

Success! Canceled rekeying (if it was started)
