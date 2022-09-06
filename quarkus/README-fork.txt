
Releasing a fork


Stepping version numbers of all submodules

mvn versions:set -DnewVersion=2.7.6.Final-nordix-1 -DoldVersion=2.7.6.Final -DgroupId=* -DartifactId=*
mvn clean install -DskipTests=true



