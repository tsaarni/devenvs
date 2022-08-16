# import client from running Keycloak

patch -p1 <<EOF
diff --git a/scripts/import-client.mjs b/scripts/import-client.mjs
index 5cbc001a..fc414d33 100755
--- a/scripts/import-client.mjs
+++ b/scripts/import-client.mjs
@@ -12,7 +12,7 @@ await importClient();

 async function importClient() {
   const adminClient = new KcAdminClient({
-    baseUrl: "http://127.0.0.1:8180",
+    baseUrl: "http://127.0.0.1:8080",
     realmName: "master",
   });

diff --git a/src/environment.ts b/src/environment.ts
index 124ba21f..60474afa 100644
--- a/src/environment.ts
+++ b/src/environment.ts
@@ -25,8 +25,8 @@ const realm =
 // The default environment, used during development.
 const defaultEnvironment: Environment = {
   loginRealm: realm,
-  authServerUrl: "http://localhost:8180",
-  authUrl: "http://localhost:8180",
+  authServerUrl: "http://localhost:8080",
+  authUrl: "http://localhost:8080",
   consoleBaseUrl: "/admin/master/console/",
   resourceUrl: ".",
   masterRealm: "master",

EOF

npm run server:import-client




# run keycloak and modify client settings for "security-admin-console-v2"
#
# change all ports from 8080 to 8181



# run admin ui with node

npm run dev
