#!/usr/bin/env python3

from ldap3 import Server, Connection, ASYNC_STREAM
server = Server('localhost', 9011)

#conn = Connection(server, user='cn=config', password='secret', auto_bind=True)
#conn.search('cn=Bind,cn=Operations,cn=Monitor', '(objectclass=*)', attributes=['monitorOpInitiated'])
#print(conn.entries)
def change_detected(result_entry):
    print(result_entry['dn'])
    print(result_entry['attributes'])

conn = Connection(server, user='cn=config', password='secret', client_strategy=ASYNC_STREAM, auto_bind=True)
#p = conn.extend.standard.persistent_search('cn=Bind,cn=Operations,cn=Monitor', '(objectclass=*)', callback=change_detected)


#p = conn.extend.standard.persistent_search('cn=Bind,cn=Operations,cn=Monitor', '(objectclass=*)', streaming=False)
p = conn.extend.standard.persistent_search('cn=accesslog', '(objectclass=*)', streaming=False)

p.start()
while True:
    print(p.next(block=True))
