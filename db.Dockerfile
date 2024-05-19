FROM postgres:14.2-alpine


COPY pg_hba.conf /var/lib/postgresql/data/pg_hba.conf
# COPY schema.sql /docker-entrypoint-initdb.d/
