

--git a/pom.xml b/pom.xml
index 6d63a1bf8e..fcfe301744 100644
--- a/pom.xml
+++ b/pom.xml
@@ -229,7 +229,6 @@
         <!-- Frontend -->
         <node.version>v18.18.2</node.version>
         <pnpm.version>8.10.2</pnpm.version>
-        <pnpm.args.install>install --prefer-offline --frozen-lockfile --ignore-scripts</pnpm.args.install>
     </properties>

     <url>http://keycloak.org</url>
@@ -1905,23 +1904,6 @@
                     <groupId>com.github.eirslett</groupId>
                     <artifactId>frontend-maven-plugin</artifactId>
                     <version>${frontend.plugin.version}</version>
-                    <executions>
-                        <execution>
-                            <goals>
-                                <goal>install-node-and-pnpm</goal>
-                            </goals>
-                        </execution>
-                        <execution>
-                            <id>pnpm-install</id>
-                            <goals>
-                                <goal>pnpm</goal>
-                            </goals>
-                            <configuration>
-                                <arguments>${pnpm.args.install}</arguments>
-                                <workingDirectory>${maven.multiModuleProjectDirectory}</workingDirectory>
-                            </configuration>
-                        </execution>
-                    </executions>
                     <configuration>
                         <nodeVersion>${node.version}</nodeVersion>
                         <pnpmVersion>${pnpm.version}</pnpmVersion>
@@ -1950,6 +1932,34 @@
                 </plugin>
             </plugins>
         </pluginManagement>
+
+        <plugins>
+            <plugin>
+                <groupId>com.github.eirslett</groupId>
+                <artifactId>frontend-maven-plugin</artifactId>
+                <version>${frontend.plugin.version}</version>
+
+                <executions>
+                    <execution>
+                        <goals>
+                            <goal>install-node-and-pnpm</goal>
+                        </goals>
+                    </execution>
+                    <execution>
+                        <id>pnpm-install</id>
+                        <inherited>false</inherited>
+                        <goals>
+                            <goal>pnpm</goal>
+                        </goals>
+                        <configuration>
+                            <arguments>--recursive install --prefer-offline --frozen-lockfile --ignore-scripts</arguments>
+                            <workingDirectory>${maven.multiModuleProjectDirectory}</workingDirectory>
+                        </configuration>
+                    </execution>
+                </executions>
+            </plugin>
+        </plugins>
+
     </build>

     <profiles>
diff --git a/themes/pom.xml b/themes/pom.xml
index 410cf9938e..9f48f9130c 100755
--- a/themes/pom.xml
+++ b/themes/pom.xml
@@ -4,6 +4,7 @@
         <artifactId>keycloak-parent</artifactId>
         <groupId>org.keycloak</groupId>
         <version>999.0.0-SNAPSHOT</version>
+        <relativePath>../pom.xml</relativePath>
     </parent>
     <modelVersion>4.0.0</modelVersion>

@@ -12,6 +13,19 @@
     <description />
     <packaging>jar</packaging>

+    <dependencies>
+        <dependency>
+            <groupId>org.keycloak</groupId>
+            <artifactId>keycloak-admin-ui</artifactId>
+            <version>${project.version}</version>
+        </dependency>
+        <dependency>
+            <groupId>org.keycloak</groupId>
+            <artifactId>keycloak-account-ui</artifactId>
+            <version>${project.version}</version>
+        </dependency>
+    </dependencies>
+
     <properties>
         <dir.common>src/main/resources/theme/keycloak/common/resources</dir.common>
         <dir.account2>src/main/resources/theme/keycloak.v2/account/src</dir.account2>
@@ -132,10 +146,12 @@
                                 </goals>
                                 <configuration>
                                     <arguments>run build</arguments>
-                                    <workingDirectory>${dir.account2}</workingDirectory>
                                 </configuration>
                             </execution>
                         </executions>
+                        <configuration>
+                            <workingDirectory>${dir.account2}</workingDirectory>
+                        </configuration>
                     </plugin>
                 </plugins>
             </build>



pom.xml
https://eclipse.dev/m2e/documentation/m2e-execution-not-covered.html

        <pluginManagement>

            <plugins>

               <!--This plugin's configuration is used to store Eclipse m2e settings only. It has no influence on the Maven build itself.-->
                <plugin>
                    <groupId>org.eclipse.m2e</groupId>
                    <artifactId>lifecycle-mapping</artifactId>
                    <version>1.0.0</version>
                    <configuration>
                        <lifecycleMappingMetadata>
                            <pluginExecutions>
                                <pluginExecution>
                                    <pluginExecutionFilter>
                                        <groupId>com.google.code.maven-replacer-plugin</groupId>
                                        <artifactId>replacer</artifactId>
                                        <versionRange>[${replacer.plugin.version},)</versionRange>
                                        <goals>
                                            <goal>replace</goal>
                                        </goals>
                                    </pluginExecutionFilter>
                                    <action>
                                       <execute>
                                          <runOnIncremental>false</runOnIncremental>
                                       </execute>
                                    </action>
                                </pluginExecution>
                            </pluginExecutions>
                        </lifecycleMappingMetadata>
                    </configuration>
                </plugin>
            </plugins>
        </pluginManagement>






or directly where the plugin gets executed
https://eclipse.dev/m2e/documentation/release-notes-17.html#new-syntax-for-specifying-lifecycle-mapping-metadata

            <plugin>
                <groupId>com.google.code.maven-replacer-plugin</groupId>
                <artifactId>maven-replacer-plugin</artifactId>
                <executions>
                    <execution>
                        <?m2e execute?>
                        <phase>process-resources</phase>
                        <goals>
                            <goal>replace</goal>
                        </goals>
                    </execution>


diff --git a/js/apps/account-ui/pom.xml b/js/apps/account-ui/pom.xml
index 2eb40d0362..7682af839a 100644
--- a/js/apps/account-ui/pom.xml
+++ b/js/apps/account-ui/pom.xml
@@ -65,6 +65,7 @@
                 <artifactId>maven-replacer-plugin</artifactId>
                 <executions>
                     <execution>
+                        <?m2e execute?>
                         <phase>process-resources</phase>
                         <goals>
                             <goal>replace</goal>
@@ -139,4 +140,4 @@
             </plugin>
         </plugins>
     </build>
-</project>
\ No newline at end of file
+</project>
diff --git a/js/apps/admin-ui/pom.xml b/js/apps/admin-ui/pom.xml
index bb8c059f4c..d19f91d64e 100644
--- a/js/apps/admin-ui/pom.xml
+++ b/js/apps/admin-ui/pom.xml
@@ -107,6 +107,7 @@
                 <artifactId>maven-replacer-plugin</artifactId>
                 <executions>
                     <execution>
+                        <?m2e execute?>
                         <phase>process-resources</phase>
                         <goals>
                             <goal>replace</goal>
@@ -174,4 +175,4 @@
             </plugin>
         </plugins>
     </build>
-</project>
\ No newline at end of file
+</project>



diff --git a/themes/pom.xml b/themes/pom.xml
index 4bd9997da4..93686d5a4e 100755
--- a/themes/pom.xml
+++ b/themes/pom.xml
@@ -163,6 +163,10 @@
                                 </configuration>
                             </execution>
                         </executions>
+                        <configuration>
+                            <workingDirectory>${dir.account2}</workingDirectory>
+                            <installDirectory>${maven.multiModuleProjectDirectory}</installDirectory>
+                        </configuration>
                     </plugin>
                 </plugins>
             </build>
@@ -208,6 +212,9 @@
                                 </configuration>
                             </execution>
                         </executions>
+                        <configuration>
+                            <workingDirectory>${dir.common}</workingDirectory>
+                        </configuration>
                     </plugin>
                 </plugins>
             </build>


diff --git a/pom.xml b/pom.xml
index 93060d115f..0d76aeecd0 100644
--- a/pom.xml
+++ b/pom.xml
@@ -1893,7 +1893,7 @@
                         <nodeVersion>${node.version}</nodeVersion>
                         <pnpmVersion>${pnpm.version}</pnpmVersion>
                         <!-- Warning, this is an undocumented property https://issues.apache.org/jira/browse/MNG-6589. But there is nothing better. -->
-                        <installDirectory>${maven.multiModuleProjectDirectory}</installDirectory>
+                        <installDirectory>/home/tsaarni/work/keycloak-worktree/ldap-referral</installDirectory>
                     </configuration>
                 </plugin>
                 <plugin>



$ find . -name index.ftl
./js/apps/admin-ui/target/classes/theme/keycloak.v2/admin/index.ftl
./js/apps/account-ui/target/classes/theme/keycloak.v3/account/index.ftl
./themes/src/main/resources/theme/keycloak.v2/account/index.ftl
./themes/src/main/resources/theme/keycloak/welcome/index.ftl
./themes/target/classes/theme/keycloak.v2/account/index.ftl
./themes/target/classes/theme/keycloak/welcome/index.ftl
./examples/themes/src/main/resources/theme/logo-example/welcome/index.ftl

$ find . -name patternfly.css
./node_modules/.pnpm/@patternfly+patternfly@4.224.5/node_modules/@patternfly/patternfly/patternfly.css
./node_modules/.pnpm/patternfly@3.59.5/node_modules/patternfly/dist/css/patternfly.css
./themes/target/classes/theme/keycloak/common/resources/node_modules/patternfly/dist/css/patternfly.css

vs m2eclipse build

$ find . -name index.ftl
./themes/src/main/resources/theme/keycloak.v2/account/index.ftl
./themes/src/main/resources/theme/keycloak/welcome/index.ftl
./themes/target/classes/theme/keycloak.v2/account/index.ftl
./themes/target/classes/theme/keycloak/welcome/index.ftl
./examples/themes/src/main/resources/theme/logo-example/welcome/index.ftl
./examples/themes/target/classes/theme/logo-example/welcome/index.ftl

$ find . -name patternfly.css
./node_modules/.pnpm/@patternfly+patternfly@4.224.5/node_modules/@patternfly/patternfly/patternfly.css
./node_modules/.pnpm/patternfly@3.59.5/node_modules/patternfly/dist/css/patternfly.css



replacer plugin does not run in 
js/apps/admin-ui/pom.xml
js/apps/account-ui/pom.xml



git stash && git clean -fdx && git reset --hard && git stash pop
mvn -T4C clean install -DskipTestsuite -DskipExamples -DskipTests
mkdir -p .vscode && cp -a ~/work/devenvs/keycloak/configs/launch.json .vscode




mkdir -p .vscode
cat <<EOF | envsubst > .vscode/settings.json
{
    "java.configuration.maven.defaultMojoExecutionAction": "execute",
    //"java.jdt.ls.vmargs": "-Xmx16G -Xms100m -Dlog.level=ALL -Djdt.ls.debug=true -Dmaven.multiModuleProjectDirectory=$PWD",
    "java.jdt.ls.vmargs": "-Xmx16G -Xms100m -Dmaven.multiModuleProjectDirectory=$PWD",
}
EOF






maven.multiModuleProjectDirectory

https://issues.apache.org/jira/browse/MNG-6589
https://youtrack.jetbrains.com/issue/IDEA-190202

