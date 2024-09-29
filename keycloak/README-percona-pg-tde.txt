

# https://percona-lab.github.io/pg_tde/main/index.html
# https://www.percona.com/blog/trying-out-the-postgresql-pg_tde-tech-preview-release/
# https://hub.docker.com/r/perconalab/pg_tde


# start new cluster
kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak

kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl apply -f manifests/pg-tde.yaml



CREATE EXTENSION pg_tde;

SELECT pg_tde_add_key_provider_file('my-key-provider','/var/lib/postgresql/data/pg.key');
SELECT pg_tde_set_principal_key('principal-key', 'my-key-provider');





SELECT pg_tde_is_encrypted('component_config');
SELECT to_regclass('component_config')::oid;   # prints e.g. 25770

kubectl exec -it postgres-0 find /var/lib/postgresql/data/base/ | grep 25770
/var/lib/postgresql/data/base/16384/25770_fsm
/var/lib/postgresql/data/base/16384/25770

kubectl exec -it postgres-0 strings /var/lib/postgresql/data/base/16384/25770










The files belonging to this database system will be owned by user "postgres".
This user must also own the server process.

The database cluster will be initialized with locale "en_US.utf8".
The default database encoding has accordingly been set to "UTF8".
The default text search configuration will be set to "english".

Data page checksums are disabled.

fixing permissions on existing directory /var/lib/postgresql/data ... ok
creating subdirectories ... ok
selecting dynamic shared memory implementation ... posix
selecting default max_connections ... 100
selecting default shared_buffers ... 128MB
selecting default time zone ... Etc/UTC
creating configuration files ... ok
running bootstrap script ... ok
performing post-bootstrap initialization ... ok
syncing data to disk ... ok


Success. You can now start the database server using:

    pg_ctl -D /var/lib/postgresql/data -l logfile start

initdb: warning: enabling "trust" authentication for local connections
initdb: hint: You can change this by editing pg_hba.conf or using the option -A, or --auth-local and --auth-host, the next time you run initdb.
waiting for server to start....2024-08-13 15:16:31.151 UTC [55] LOG:  starting PostgreSQL 16.4 (Debian 16.4-1.pgdg120+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit
2024-08-13 15:16:31.158 UTC [55] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
2024-08-13 15:16:31.178 UTC [58] LOG:  database system was shut down at 2024-08-13 15:16:30 UTC
2024-08-13 15:16:31.186 UTC [55] LOG:  database system is ready to accept connections
 done
server started
2024-08-13 15:16:31.257 UTC [64] LOG:  statement: SELECT 1 FROM pg_database WHERE datname = 'keycloak' ;
2024-08-13 15:16:31.281 UTC [66] LOG:  statement: CREATE DATABASE "keycloak" ;
CREATE DATABASE


/usr/local/bin/docker-entrypoint.sh: sourcing /docker-entrypoint-initdb.d/pg-tde-create-ext.sh
2024-08-13 15:16:31.344 UTC [68] LOG:  statement: CREATE EXTENSION pg_tde;
2024-08-13 15:16:31.346 UTC [68] WARNING:  pg_tde can only be loaded at server startup. Restart required.
2024-08-13 15:16:31.346 UTC [68] LOG:  Initializing TDE principal key info
2024-08-13 15:16:31.346 UTC [68] STATEMENT:  CREATE EXTENSION pg_tde;
2024-08-13 15:16:31.346 UTC [68] LOG:  initializing TDE key provider info
2024-08-13 15:16:31.346 UTC [68] STATEMENT:  CREATE EXTENSION pg_tde;
2024-08-13 15:16:31.346 UTC [68] ERROR:  failed to register custom resource manager "test_tdeheap_custom_rmgr" with ID 128
2024-08-13 15:16:31.346 UTC [68] DETAIL:  Custom resource manager must be registered while initializing modules in shared_preload_libraries.
2024-08-13 15:16:31.346 UTC [68] STATEMENT:  CREATE EXTENSION pg_tde;
WARNING:  pg_tde can only be loaded at server startup. Restart required.
ERROR:  failed to register custom resource manager "test_tdeheap_custom_rmgr" with ID 128
DETAIL:  Custom resource manager must be registered while initializing modules in shared_preload_libraries.





kubectl exec -it postgres-0 -- psql -U keycloak -c 'CREATE EXTENSION pg_tde;'
kubectl exec -it postgres-0 -- psql -U keycloak -d template1 -c 'CREATE EXTENSION pg_tde;'






docker run --name pg_tde --rm perconalab/pg_tde
