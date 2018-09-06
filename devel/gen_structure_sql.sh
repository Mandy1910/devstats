#!/bin/sh
if [ -z "${PG_PASS}" ]
then
  echo "You need to set PG_PASS environment variable to run this script"
  exit 1
fi
./devel/drop_psql_db.sh temp &&\
PG_DB=temp ./structure &&\
PG_DB=temp GHA2DB_SKIPTABLE=1 GHA2DB_INDEX=1 GHA2DB_MGETC=y ./structure &&\
sudo -u postgres psql temp -f util_sql/current_state_all.sql &&\
sudo -u postgres pg_dump -s temp > structure.sql &&\
./devel/drop_psql_db.sh temp && echo 'structure.sql generated'
