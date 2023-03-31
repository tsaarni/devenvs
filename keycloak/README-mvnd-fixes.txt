



mvnd clean install -DskipTestsuite -DskipExamples -DskipTests
mvn -T4C clean install -DskipTestsuite -DskipExamples -DskipTests


[ERROR] Failed to execute goal com.github.eirslett:frontend-maven-plugin:1.12.1:npm (npm-ci) on project keycloak-js-adapter-jar: Failed to run task: 'npm ci --ignore-scripts' failed. org.apache.commons.exec.ExecuteException: Process exited with an error: 239 (Exit value: 239) -> [Help 1]
[ERROR] Failed to execute goal com.github.eirslett:frontend-maven-plugin:1.12.1:npm (npm-ci) on project keycloak-js-adapter: Failed to run task: 'npm ci --ignore-scripts' failed. org.apache.commons.exec.ExecuteException: Process exited with an error: 254 (Exit value: 254) -> [Help 1]
[ERROR] Failed to execute goal com.github.eirslett:frontend-maven-plugin:1.12.1:npm (npm-build) on project keycloak-admin-ui: Failed to run task: 'npm run build --workspace=admin-ui' failed. org.apache.commons.exec.ExecuteException: Process exited with an error: 1 (Exit value: 1) -> [Help 1]
[ERROR] Failed to execute goal com.github.eirslett:frontend-maven-plugin:1.12.1:npm (npm-ci) on project keycloak-js-admin-client: Failed to run task: 'npm ci --ignore-scripts' failed. org.apache.commons.exec.ExecuteException: Process exited with an error: 254 (Exit value: 254) -> [Help 1]




Support for hidden dependencies
https://github.com/apache/maven-mvnd/issues/12


Remove mvnd.builder.rule* and mvnd.builder.rules.provider.* features
https://github.com/apache/maven-mvnd/issues/264


# test compilation
CI=true npm ci --ignore-scripts
CI=true npm run build

# logs from frontend-maven-plugin npm execution
rm ~/.npm/_logs/*
ls ~/.npm/_logs/
grep error ~/.npm/_logs/*



mvn dependency:tree -f js/pom.xml
mvn dependency:tree -Dincludes=org.keycloak:keycloak-js-parent





[ERROR] Failed to execute goal org.apache.maven.plugins:maven-compiler-plugin:3.8.1-jboss-2:testCompile (default-testCompile) on project keycloak-model-map: Compilation failure: Compilation failure:
[ERROR] /home/tsaarni/work/keycloak-worktree/mvnd-documentation-build-fix/model/map/src/test/java/org/keycloak/models/map/storage/tree/TreeStorageNodePrescriptionTest.java:[22,38] cannot find symbol
[ERROR]   symbol:   class MapClientEntityFields
[ERROR]   location: package org.keycloak.models.map.client
[ERROR] /home/tsaarni/work/keycloak-worktree/mvnd-documentation-build-fix/model/map/src/test/java/org/keycloak/models/map/storage/tree/TreeStorageNodePrescriptionTest.java:[23,38] cannot find symbol
[ERROR]   symbol:   class MapClientEntityFields
[ERROR]   location: package org.keycloak.models.map.client



[ERROR] Failed to execute goal org.apache.maven.plugins:maven-compiler-plugin:3.8.1:compile (default-compile) on project keycloak-model-map-jpa: Compilation failure: Compilation failure:
[ERROR] /home/tsaarni/work/keycloak-worktree/mvnd-documentation-build-fix/model/map-jpa/src/main/java/org/keycloak/models/map/storage/jpa/authorization/policy/delegate/JpaPolicyDelegateProvider.java:[28,52] cannot find symbol
[ERROR]   symbol:   class MapPolicyEntityFields
[ERROR]   location: package org.keycloak.models.map.authorization.entity
[ERROR] /home/tsaarni/work/keycloak-worktree/mvnd-documentation-build-fix/model/map-jpa/src/main/java/org/keycloak/models/map/storage/jpa/authorization/policy/entity/JpaPolicyMetadata.java:[20,52] cannot find symbol



[ERROR] Failed to execute goal com.github.eirslett:frontend-maven-plugin:1.12.1:npm (npm-ci) on project keycloak-js-adapter: Failed to run task: 'npm ci --ignore-scripts' failed. org.apache.commons.exec.ExecuteException: Process exited with an error: 217 (Exit value: 217) -> [Help 1]






1061 verbose stack Error: EEXIST: file already exists, symlink '../libs/ui-shared' -> '/home/tsaarni/work/keycloak-worktree/mvnd-parallel-build-fix/js/node_modules/ui-shared'
1062 verbose cwd /home/tsaarni/work/keycloak-worktree/mvnd-parallel-build-fix/js
1063 verbose Linux 5.15.0-67-generic
1064 verbose node v18.14.2
1065 verbose npm  v9.5.0
1066 error code EEXIST
1067 error syscall symlink
1068 error path ../libs/ui-shared
1069 error dest /home/tsaarni/work/keycloak-worktree/mvnd-parallel-build-fix/js/node_modules/ui-shared
1070 error errno -17
1071 error EEXIST: file already exists, symlink '../libs/ui-shared' -> '/home/tsaarni/work/keycloak-worktree/mvnd-parallel-build-fix/js/node_modules/ui-shared'
1072 error File exists: /home/tsaarni/work/keycloak-worktree/mvnd-parallel-build-fix/js/node_modules/ui-shared
1073 error Remove the existing file and try again, or run npm
1074 error with --force to overwrite files recklessly.
1075 verbose exit -17
