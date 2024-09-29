

https://www.tenable.com/plugins/nessus/12085

https://wiki.owasp.org/index.php/Securing_tomcat
https://tomcat.apache.org/tomcat-9.0-doc/security-howto.html
https://tomcat.apache.org/tomcat-9.0-doc/security-howto.html#Valves


# change port if necessary (e.g. 8080 -> 18080)
conf/server.xml

# run server on foreground
bin/catalina.sh run


### Default applications

# check for example following pages for the default apps
# for more info, see
#  https://tomcat.apache.org/tomcat-9.0-doc/security-howto.html#Default_web_applications

http://localhost:18080/
http://localhost:18080/docs
http://localhost:18080/examples/


# delete web apps
rm -r webapps/ROOT webapps/docs webapps/examples
rm -r webapps/manager webapps/host-manager




### Hide server version info in default error page

# Request for non-existing page returns server information in error page
http http://localhost:18080/foo


# To change the returned server info, create following file
#
# See also
#   https://tomcat.apache.org/tomcat-9.0-doc/security-howto.html#Valves
# Note:
#   See next tip "error page and backtraces" which also removes server info


mkdir -p lib/org/apache/catalina/util

cat >lib/org/apache/catalina/util/ServerInfo.properties <<EOF
server.info=Apache Tomcat
EOF

# restart server




### error page and backtraces

CATALINA_HOME=/home/tsaarni/package/tomcat/apache-tomcat-9.0.84
(cd webapps/failing/WEB-INF/classes; javac -classpath $CATALINA_HOME/lib/servlet-api.jar FailingExample.java)
cp -r webapps/failing/ $CATALINA_HOME/webapps/
# restart server

http http://localhost:18080/failing/backtrace


# To disable backtraces, append following to conf/server.xml under <Host>
# see
#   https://tomcat.apache.org/tomcat-9.0-doc/security-howto.html#Valves
#   https://stackoverflow.com/questions/65987346/tomcat-8-5-50-the-errorreportvalve-is-not-working


        <Valve className="org.apache.catalina.valves.ErrorReportValve"
               showReport="false"
               showServerInfo="false"/>


# Note:
# showReport="false" removes the backtrace
# showServerInfo="false" removes the server info



### Script for tomcat default files

https://vulners.com/nessus/TOMCAT_SERVER_DEFAULT_FILES.NASL

Following line checks if version string is in error page

   if ( 'Apache Tomcat/' >< response &&

Which matches for example following error page string:

   Apache Tomcat/9.0.84


### Nessus Essentials download

https://www.tenable.com/downloads/nessus?loginAttempted=true

# free license from
https://www.tenable.com/products/nessus/nessus-essentials
