

Parallel builds stopped working
https://github.com/keycloak/keycloak/issues/24571



My suggestion to run pnpm install recursively on top level
https://github.com/keycloak/keycloak/pull/24537#discussion_r1389486794



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




*** Plugin execution not covered by lifecycle configuration



                 <executions>
                     <execution>
+                        <?m2e ignore?>
                         <id>echo-output</id>
                     </execution>
                 </executions>




*** No mojo definitions were found



Failed to execute mojo org.apache.maven.plugins:maven-plugin-plugin:3.11.0:descriptor {execution: generate-descriptor} (org.apache.maven.plugins:maven-plugin-plugin:3.11.0:descriptor:generate-descriptor:process-classes)

org.eclipse.core.runtime.CoreException: Failed to execute mojo org.apache.maven.plugins:maven-plugin-plugin:3.11.0:descriptor {execution: generate-descriptor}
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.executeMojo(MavenExecutionContext.java:404)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.lambda$2(MavenExecutionContext.java:355)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.executeBare(MavenExecutionContext.java:458)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.execute(MavenExecutionContext.java:339)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.execute(MavenExecutionContext.java:354)
	at org.eclipse.m2e.core.project.configurator.MojoExecutionBuildParticipant.build(MojoExecutionBuildParticipant.java:57)
	at org.eclipse.m2e.core.internal.builder.MavenBuilderImpl.lambda$2(MavenBuilderImpl.java:159)
	at java.base/java.util.LinkedHashMap.forEach(Unknown Source)
	at org.eclipse.m2e.core.internal.builder.MavenBuilderImpl.build(MavenBuilderImpl.java:139)
	at org.eclipse.m2e.core.internal.builder.MavenBuilder$1.method(MavenBuilder.java:164)
	at org.eclipse.m2e.core.internal.builder.MavenBuilder$1.method(MavenBuilder.java:1)
	at org.eclipse.m2e.core.internal.builder.MavenBuilder$BuildMethod.lambda$1(MavenBuilder.java:109)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.executeBare(MavenExecutionContext.java:458)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.execute(MavenExecutionContext.java:292)
	at org.eclipse.m2e.core.internal.builder.MavenBuilder$BuildMethod.lambda$0(MavenBuilder.java:100)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.executeBare(MavenExecutionContext.java:458)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.execute(MavenExecutionContext.java:339)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.execute(MavenExecutionContext.java:278)
	at org.eclipse.m2e.core.internal.builder.MavenBuilder$BuildMethod.execute(MavenBuilder.java:83)
	at org.eclipse.m2e.core.internal.builder.MavenBuilder.build(MavenBuilder.java:192)
	at org.eclipse.core.internal.events.BuildManager$2.run(BuildManager.java:1077)
	at org.eclipse.core.runtime.SafeRunner.run(SafeRunner.java:47)
	at org.eclipse.core.internal.events.BuildManager.basicBuild(BuildManager.java:296)
	at org.eclipse.core.internal.events.BuildManager.basicBuild(BuildManager.java:352)
	at org.eclipse.core.internal.events.BuildManager$1.run(BuildManager.java:441)
	at org.eclipse.core.runtime.SafeRunner.run(SafeRunner.java:47)
	at org.eclipse.core.internal.events.BuildManager.basicBuild(BuildManager.java:444)
	at org.eclipse.core.internal.events.BuildManager.basicBuildLoop(BuildManager.java:555)
	at org.eclipse.core.internal.events.BuildManager.basicBuildLoop(BuildManager.java:503)
	at org.eclipse.core.internal.events.BuildManager.build(BuildManager.java:585)
	at org.eclipse.core.internal.resources.Workspace.buildInternal(Workspace.java:594)
	at org.eclipse.core.internal.resources.Workspace.build(Workspace.java:483)
	at org.eclipse.jdt.ls.core.internal.handlers.BuildWorkspaceHandler.buildWorkspace(BuildWorkspaceHandler.java:65)
	at org.eclipse.jdt.ls.core.internal.handlers.JDTLanguageServer.lambda$28(JDTLanguageServer.java:1001)
	at org.eclipse.jdt.ls.core.internal.handlers.JDTLanguageServer.lambda$61(JDTLanguageServer.java:1236)
	at java.base/java.util.concurrent.CompletableFuture$UniApply.tryFire(Unknown Source)
	at java.base/java.util.concurrent.CompletableFuture$Completion.exec(Unknown Source)
	at java.base/java.util.concurrent.ForkJoinTask.doExec(Unknown Source)
	at java.base/java.util.concurrent.ForkJoinPool$WorkQueue.topLevelExec(Unknown Source)
	at java.base/java.util.concurrent.ForkJoinPool.scan(Unknown Source)
	at java.base/java.util.concurrent.ForkJoinPool.runWorker(Unknown Source)
	at java.base/java.util.concurrent.ForkJoinWorkerThread.run(Unknown Source)
Caused by: org.apache.maven.plugin.MojoExecutionException: Error extracting plugin descriptor: 'No mojo definitions were found for plugin: org.keycloak:keycloak-distribution-licenses-maven-plugin.'
	at org.apache.maven.plugin.plugin.DescriptorGeneratorMojo.generate(DescriptorGeneratorMojo.java:368)
	at org.apache.maven.plugin.plugin.AbstractGeneratorMojo.execute(AbstractGeneratorMojo.java:90)
	at org.apache.maven.plugin.DefaultBuildPluginManager.executeMojo(DefaultBuildPluginManager.java:126)
	at org.eclipse.m2e.core.internal.embedder.MavenExecutionContext.executeMojo(MavenExecutionContext.java:402)
	... 41 more
Caused by: org.apache.maven.plugin.descriptor.InvalidPluginDescriptorException: No mojo definitions were found for plugin: org.keycloak:keycloak-distribution-licenses-maven-plugin.
	at org.apache.maven.tools.plugin.scanner.DefaultMojoScanner.populatePluginDescriptor(DefaultMojoScanner.java:136)
	at org.apache.maven.plugin.plugin.DescriptorGeneratorMojo.generate(DescriptorGeneratorMojo.java:355)
	... 44 more




diff --git a/distribution/maven-plugins/pom.xml b/distribution/maven-plugins/pom.xml
index fe9b7db2e3..f73c9f468e 100644
--- a/distribution/maven-plugins/pom.xml
+++ b/distribution/maven-plugins/pom.xml
@@ -83,6 +83,9 @@
                 <plugin>
                     <groupId>org.apache.maven.plugins</groupId>
                     <artifactId>maven-plugin-plugin</artifactId>
+                    <configuration>
+                        <skipErrorNoDescriptorsFound>true</skipErrorNoDescriptorsFound>
+                    </configuration>
                     <executions>
                         <execution>
                             <id>generate-descriptor</id>










diff --git a/js/apps/account-ui/pom.xml b/js/apps/account-ui/pom.xml
index fe932ccc64..57aeea5536 100644
--- a/js/apps/account-ui/pom.xml
+++ b/js/apps/account-ui/pom.xml
@@ -71,6 +71,7 @@
                 <artifactId>frontend-maven-plugin</artifactId>
                 <executions>
                     <execution>
+                        <?m2e ignore?>
                         <id>lib-build</id>
                         <goals>
                             <goal>pnpm</goal>
@@ -83,6 +84,7 @@
                         </configuration>
                     </execution>
                     <execution>
+                        <?m2e ignore?>
                         <id>pack</id>
                         <phase>package</phase>
                         <goals>
@@ -99,6 +101,7 @@
                 <artifactId>maven-replacer-plugin</artifactId>
                 <executions>
                     <execution>
+                        <?m2e ignore?>
                         <phase>process-resources</phase>
                         <goals>
                             <goal>replace</goal>
@@ -198,4 +201,4 @@
             </plugin>
         </plugins>
     </build>
-</project>
\ No newline at end of file
+</project>
diff --git a/js/apps/admin-ui/pom.xml b/js/apps/admin-ui/pom.xml
index 2ea545c7a6..442a318b7a 100644
--- a/js/apps/admin-ui/pom.xml
+++ b/js/apps/admin-ui/pom.xml
@@ -75,6 +75,7 @@
                 <artifactId>maven-replacer-plugin</artifactId>
                 <executions>
                     <execution>
+                        <?m2e ignore?>
                         <phase>process-resources</phase>
                         <goals>
                             <goal>replace</goal>











diff --git a/distribution/licenses-common/pom.xml b/distribution/licenses-common/pom.xml
index bf615a63ee..a395a39146 100644
--- a/distribution/licenses-common/pom.xml
+++ b/distribution/licenses-common/pom.xml
@@ -27,20 +27,5 @@
     <packaging>jar</packaging>
     <name>Keycloak Distribution Licenses Common</name>

-    <build>
-        <resources>
-            <resource>
-                <targetPath>keycloak-licenses-common</targetPath>
-                <directory>../../</directory>
-                <includes>
-                    <include>LICENSE.txt</include>
-                </includes>
-            </resource>
-            <resource>
-                <targetPath>keycloak-licenses-common</targetPath>
-                <directory>src/main/resources</directory>
-            </resource>
-        </resources>
-    </build>

 </project>
