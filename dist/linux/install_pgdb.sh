#!/bin/bash

# download postgres repo
yes | yum install https://download.postgresql.org/pub/repos/yum/11/redhat/rhel-7-x86_64/pgdg-redhat11-11-2.noarch.rpm
yes | yum install postgresql11 postgresql11-server postgresql11-contrib postgresql11-libs

# setup postgres group
groupadd pg_wheel && getent group pg_wheel
gpasswd -a postgres pg_wheel

# cleanup and create folders for db
rm -Rf /usr/local/pgsql
mkdir /usr/local/pgsql
mkdir /usr/local/pgsql/data
chown -R postgres:pg_wheel /usr/local/pgsql

# generate setup script
db_setup_sh=/var/tmp/setup_db.sh
rm -f $db_setup_sh
echo "echo \"umask 077\" >> ~/.bash_rc" >> $db_setup_sh
echo "source ~/.bash_rc" >> $db_setup_sh
echo "cd /usr/local/pgsql" >> $db_setup_sh
echo "export PGHOST=${TDS_DB_HOSTNAME}" >> $db_setup_sh
echo "export PGPORT=${TDS_DB_PORT}" >> $db_setup_sh
echo "export PGDATA=/usr/local/pgsql/data" >> $db_setup_sh
echo "/usr/pgsql-11/bin/pg_ctl -D /usr/local/pgsql/data initdb" >> $db_setup_sh
echo "/usr/pgsql-11/bin/pg_ctl -D /usr/local/pgsql/data -l /usr/local/pgsql/logfile start" >> $db_setup_sh

echo "echo \"local all postgres peer\" >> /usr/local/pgsql/data/pg_hba.conf" >> $db_setup_sh
echo "echo \"local all all peer\" >> /usr/local/pgsql/data/pg_hba.conf" >> $db_setup_sh
echo "echo \"listen_addresses = '*'\" >> /usr/local/pgsql/data/pg_hba.conf" >> $db_setup_sh
echo "echo \"host all postgres 127.0.0.1/32 md5\" >> /usr/local/pgsql/data/pg_hba.conf" >> $db_setup_sh

echo "psql -c \"alter system set log_connections = 'on';\"" >> $db_setup_sh
echo "psql -c \"alter system set log_disconnections = 'on';\"" >> $db_setup_sh
echo "psql -c \"select pg_reload_conf();\"" >> $db_setup_sh
echo "psql -c \"CREATE EXTENSION \\\"uuid-ossp\\\";\"" >> $db_setup_sh

echo "psql -c \"CREATE USER ${TDS_DB_USERNAME} WITH SUPERUSER PASSWORD '${TDS_DB_PASSWORD}';\"" >> $db_setup_sh
echo "psql -c \"CREATE DATABASE ${TDS_DB_NAME}\"" >> $db_setup_sh

echo "psql -c \"GRANT ALL PRIVILEGES ON DATABASE ${TDS_DB_NAME} TO ${TDS_DB_USERNAME};\"" >> $db_setup_sh
echo "psql -c \"ALTER ROLE ${TDS_DB_USERNAME} NOCREATEROLE;\"" >> $db_setup_sh
echo "psql -c \"ALTER ROLE ${TDS_DB_USERNAME} NOCREATEDB;\"" >> $db_setup_sh
echo "psql -c \"ALTER ROLE ${TDS_DB_USERNAME} NOREPLICATION;\"" >> $db_setup_sh
echo "psql -c \"ALTER ROLE ${TDS_DB_USERNAME} NOBYPASSRLS;\"" >> $db_setup_sh
echo "psql -c \"ALTER ROLE ${TDS_DB_USERNAME} NOINHERIT;\"" >> $db_setup_sh

# run setup script as postgres user and cleanup
chown -R postgres:pg_wheel $db_setup_sh
chmod +x $db_setup_sh
su postgres -c "source $db_setup_sh"
rm -f $db_setup_sh
