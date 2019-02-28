#!/bin/bash

# READ .env file 
echo PWD IS $(pwd)
if [ -f ~/tdservice.env ]; then 
    echo Reading Installation options from `realpath ~/tdservice.env`
    source ~/tdservice.env
elif [ -f ../tdservice.env ]; then
    echo Reading Installation options from `realpath ../tdservice.env`
    source ../tdservice.env
else
    echo No .env file found
fi

# Export all known variables
export TDS_DB_HOSTNAME
export TDS_DB_PORT
export TDS_DB_USERNAME
export TDS_DB_PASSWORD
export TDS_DB_NAME

export TDS_DB_PORT

export TDS_ADMIN_USERNAME
export TDS_ADMIN_PASSWORD

export TDS_TLS_HOSTS

if [[ $EUID -ne 0 ]]; then 
    echo "This installer must be run as root"
    exit 1
fi

echo Setting up Threat Detection Service Linux User...
id -u tds 2> /dev/null || useradd tds

echo Installing Threat Detection Service...

cp tdservice /usr/bin/tdservice
chmod +x /usr/bin/tdservice 
chmod 755 /usr/bin/tdservice

# Create configuration directory in /etc
mkdir -p /etc/tdservice && chown tds:tds /etc/tdservice
chmod 775 /etc/tdservice
chmod g+s /etc/tdservice
# Create data dir in /var/lib
mkdir -p /var/lib/tdservice && chown tds:tds /var/lib/tdservice
chmod 775 /var/lib/tdservice
chmod g+s /var/lib/tdservice
# Create logging dir in /var/log
mkdir -p /var/log/tdservice && chown tds:tds /var/log/tdservice
chmod 775 /var/log/tdservice
chmod g+s /var/log/tdservice

# Install systemd script
cp tdservice.service /etc/systemd/system/tdservice.service

# Enable systemd service
systemctl enable tdservice

# check if TDS_NOSETUP is defined
if [[ -z $TDS_NOSETUP ]]; then 
    tdservice setup all
    SETUPRESULT=$?
    if [ ${SETUPRESULT} == 0 ]; then 
        echo Installation completed successfully!
        systemctl start tdservice
    else 
        echo Installation completed with errors
    fi
else 
    echo flag TDS_NOSETUP is defined, skipping setup
    echo Installation completed successfully!
fi