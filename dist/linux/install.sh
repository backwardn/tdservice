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

COMPONENT_NAME=tdservice
PRODUCT_HOME=/opt/$COMPONENT_NAME
BIN_PATH=$PRODUCT_HOME/bin
LOG_PATH=/var/log/$COMPONENT_NAME/
CONFIG_PATH=/etc/$COMPONENT_NAME/

mkdir -p $BIN_PATH && chown tds:tds $BIN_PATH/
cp $COMPONENT_NAME $BIN_PATH/
chmod 775 $BIN_PATH/*
ln -sfT $BIN_PATH/$COMPONENT_NAME /usr/bin/$COMPONENT_NAME

# Create configuration directory in /etc
mkdir -p $CONFIG_PATH && chown tds:tds $CONFIG_PATH
chmod 775 $CONFIG_PATH
chmod g+s $CONFIG_PATH

# Create logging dir in /var/log
mkdir -p $LOG_PATH && chown tds:tds $LOG_PATH
chmod 775 $LOG_PATH
chmod g+s $LOG_PATH

# Install systemd script
cp tdservice.service $PRODUCT_HOME && chown tds:tds $PRODUCT_HOME/tdservice.service && chown tds:tds $PRODUCT_HOME

# Enable systemd service
systemctl enable $PRODUCT_HOME/tdservice.service

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