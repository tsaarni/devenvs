package envoy.authz

import input.attributes.request.http as http_request


# https://www.openpolicyagent.org/docs/latest/oauth-oidc/
jwks_request(url) = http.send({
    "url": url,
    "method": "GET",
    "force_cache": true,
    "force_cache_duration_seconds": 60,
})

jwks = json.marshal(jwks_request("http://localhost:8081/auth/realms/master/protocol/openid-connect/certs").body)

token = {"valid": valid, "payload": payload} {
    [_, encoded] := split(http_request.headers.authorization, " ")
    [valid, _, payload] := io.jwt.decode_verify(encoded, {"cert": jwks})
}

#default allow = false

# default allow = {
#     "allowed": false,
#     "body": "{ \"errors\": [\"Unauthorized request\"] }",
#     "http_status": 403,
#     "headers": {
#         "content-type": "application/json",
#     }
# }

allow = response {
    count(errors) > 0
    response := {
        "allowed": false,
        "body": json.marshal({ "errors": errors}),
        "http_status": 403,
        "headers": {
            "content-type": "application/json",
        }
    }
}


errors[err] {
    not token.valid
    err := "token is invalid"
}

errors[err] {
    token.valid
    now := time.now_ns() / 1000000000
    token.payload.exp < now
    err := "token is expired"
}

allow {
    is_token_valid
    action_allowed
}

is_token_valid {
    token.valid
    #now := time.now_ns() / 1000000000
    #token.payload.nbf <= now
    #now < token.payload.exp
}

action_allowed {
  http_request.method == "GET"
  #token.payload.role == "guest"
  trace(sprintf("TOKEN %v", [token.payload]))
  #token.payload.scope == "email profile"
  glob.match("/allow*", [], http_request.path)
}
