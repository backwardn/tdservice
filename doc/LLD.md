# Threat Detection Service Low Level Documentation

## Acronyms

|     | Description                 |
|-----|-----------------------------|
| TDT | Threat Detection Technology |
| TDA | Threat Detection Agent      |
| TDS | Threat Detection Service    |
| TDL | Threat Detection Library    |
| TDD | Threat Detection Driver     |
|     |                             |


# Overview 

The `Threat Detection Service` is a web service whose purpose is to manage many deployed instances of the `Threat Detection Agent`. 

The `Threat Detection Service` has two core functionalities:

1. Aggregate threat reports from `Threat Detection Agent` (Phase 1)
2. Push updated heuristics models to `Threat Detection Agent` (Phase 2)

# API Endpoints

## Node Management

### POST `/tds/hosts`
Register an Agent to the Service. 

`Content-Type: application/json`

`Authorization: Basic ...`

Example body:
```json
{
  "hostname": "10.105.168.1",
  "signingcert": "base64",
  "os": "linux",
  // below is an embedding of the `discover` function output
  "version": "1.2.1",
  "build": "201910012012"
}
```

Example Response:
```json
{
  "id": "123e4567-e89b-12d3-a456-426655440000",
  "hostname": "10.105.168.1",
  // below is an embedding of the `discover` function output
  "version": "1.2.1",
  "build": "201910012012",
  "os": "linux",
  "status": "online"
}
```

### GET `/tds/hosts`

Query all registered hosts

`GET /tds/hosts`

Example Response:
```json
[
    {
        "id": "123e4567-e89b-12d3-a456-426655440000",
        "hostname": "10.105.168.1",
        // below is an embedding of the `discover` function output
        "version": "1.2.1",
        "build": "201910012012",
        "os": "linux",
        "status": "online",
    }, 
    {
        "id": "223e4567-e89b-12d3-a456-426655440000",
        "hostname": "10.105.168.2",
        // below is an embedding of the `discover` function output
        "version": "1.2.1",
        "build": "201910012012",
        "os": "linux",
        "status": "offline",
    }, 
]
```

Available Query parameters:

- hostname=(hostname)
- version=(version)
- build=(build)
- os=(os)
- status=(status)

Query parameters can be conjoined in any combination, so for example: `GET /tds/hosts?os=linux&status=offline`

### GET `/tds/hosts/{id-uuidv4}` or GET `/tds/hosts/{hostname}`

Get a single host by ID or its hostname

`GET /tds/hosts/123e4567-e89b-12d3-a456-426655440000`

`GET /tds/hosts/10.105.168.1`

Example Response:

```json
{
    "id": "123e4567-e89b-12d3-a456-426655440000",
    "hostname": "10.105.168.1",
    // below is an embedding of the `discover` function output
    "version": "1.2.1",
    "build": "201910012012",
    "os": "linux",
    "status": "online",
}
```

### DELETE `/tds/hosts/{id-uuidv4}` or DELETE `/tds/hosts/{hostname}`

Unregister node from `TDS`

## Reports

### POST `/tds/reports`
Create a new threat detection report event.

Example body:

```json
{
    "host_id": "123e4567-e89b-12d3-a456-426655440000",
    "detection": {
        "description": "Crypto mining suspected",
        "pid": 1234,
        "tid": 3, // thread id
        "process_name": "malicious.exe",
        "process_image_path": "C:\temp\malicious.exe",
        "process_cmd_line": "C:\temp\malicious.exe -h exfil.onion",
        "timestamp": 1234758758, // time since unix epoch
        "severity": 10,
        "profile_name": "rfc_ml_sca",
        "cve_ids": "CVE-...",
        "threat_class": "spectre variant 1",
    },
    "error": { 
        "description": "error message",
    }
}
```

TDS will create the ID, and log the event date. 

Example response:
```json
{
    "id": "123e4567-e89b-12d3-a456-426655440000",
    "date": "2019-02-04T20:56:31Z",
    "hostname": "10.1.1.1",
    "detection": {
        "description": "Crypto mining suspected",
        "pid": 1234,
        "tid": 3, // thread id
        "process_name": "malicious.exe",
        "process_image_path": "C:\temp\malicious.exe",
        "process_cmd_line": "C:\temp\malicious.exe -h exfil.onion",
        "timestamp": 1234758758, // time since unix epoch
        "severity": 10,
        "profile_name": "rfc_ml_sca",
        "cve_ids": "CVE-...",
        "threat_class": "spectre variant 1",
    },
    "error": {
        "description": "error message",
    }
}
```

### GET `/tds/reports`
Query reports by filter criteria.

With no query parameters, it returns ALL reports:
```json
[
    <report 1>,
    <report 2>,
    ...
]
```

With query parameter `?hostname=10.1.1.1`, returns all reports from the specified host

`GET /tds/reports?hostname=10.1.1.1`
```json
[
    {
        "id": "123e4567-e89b-12d3-a456-426655440000",
        "date": "2019-02-04T20:56:31Z",
        "hostname": "10.1.1.1",
        "detection": {
            "description": "Crypto mining suspected",
            "pid": 1234,
            "tid": 3, // thread id
            "process_name": "malicious.exe",
            "process_image_path": "C:\temp\malicious.exe",
            "process_cmd_line": "C:\temp\malicious.exe -h exfil.onion",
        },
        "error": {
            "description": "error message",
        }
    },
    {
        "id": "223e4567-e89b-12d3-a456-426655440000",
        "date": "2019-02-04T20:56:31Z",
        "hostname": "10.1.1.1",
        "detection": {
            "description": "Side channel detected",
            "pid": 1235,
            "tid": 3, // thread id
            "process_name": "chrome.exe",
            "process_image_path": "C:\Users\admin\AppData\Roaming\chrome.exe",
            "process_cmd_line": "C:\Users\admin\AppData\Roaming\chrome.exe",
        },
        "error": {
            "description": "error message",
        }
    }
]
```

With query parameter `?from=<RFC3339Date>`, returns all reports with date later than or equal to the specified date.

With query parameter `?to=<RFC3339Date>`, returns all reports with date before or equal to the specified date

`GET /tds/reports?from=2018-02-04T20:56:31Z&to=2020-02-04T20:56:31Z`

```json
[
    {
        "id": "123e4567-e89b-12d3-a456-426655440000",
        "date": "2019-02-04T20:56:31Z",
        "hostname": "10.1.1.1",
        "detection": {
            "description": "Crypto mining suspected",
            "pid": 1234,
            "tid": 3, // thread id
            "process_name": "malicious.exe",
            "process_image_path": "C:\temp\malicious.exe",
            "process_cmd_line": "C:\temp\malicious.exe -h exfil.onion",
        },
        "error": {
            "description": "error message",
        }
    },
]
```

Available Query parameters:

- hostname=(hostname)
- hostid=(host uuid)
- from=(from_date)
- to=(to_date)

### GET `/tds/reports/{id}`
Get a single report by its unique identifier

`GET /tds/reports/123e4567-e89b-12d3-a456-426655440000`

```json
{
    "id": "123e4567-e89b-12d3-a456-426655440000",
    "date": "2019-02-04T20:56:31Z",
    "hostname": "10.1.1.1",
    "detection": {
        "description": "Crypto mining suspected",
        "pid": 1234,
        "tid": 3, // thread id
        "process_name": "malicious.exe",
        "process_image_path": "C:\temp\malicious.exe",
        "process_cmd_line": "C:\temp\malicious.exe -h exfil.onion",
    },
    "error": {
        "description": "error message",
    }
}
```

## Configuration and Heuristics

API's for pushing configuration and heuristics stubbed out until Phase 2

# Threat Detection Service Installation

There are two modes of installation:

1. Bare Metal
2. Container

## Bare Metal Installation

The daemon will create and use the following files on the OS:

1. /var/run/tdservice/tdservice.pid (PID file to track daemon)
2. /var/log/tdservice/tdservice.log
3. /var/log/tdservice/http.log
4. /var/lib/tdservice/* (misc files)
5. /etc//tdservice/config.yaml (Configuration)
6. /usr/*/bin/tdservice (executable binary)
7. /etc/tdservice/key.pem (TLS key)
8. /etc/tdservice/cert.pem (TLS cert)

## Container Installation

Since `TDS` is a standalone web service, container deployment is trivial. 

All necessary setup options should be readable from environment variables, so the container can be spun up by only passing environment variables

# TLS Configuration

By default, `TDS` will use the system's cert pool for trusting the `TDA` TLS certificate. If a certificate is to be manually added to the trust pool, place the respective certificate file in `/etc/tdservice/certs/<hostname.pem>

For example, if the `TDA` node is at https://tda.node-2.intel.com, place its pem cert in /etc/tdservice/certs/tda.node-2.intel.com

`TDS` will then use both the systems cert pool as well as any certificates found in this directory for TLS validation. Ensure any missing intermediates are concatenated with the supplied .pem.

# Command Line Operations

## Setup

```bash
> tdservice setup 
  Available setup tasks:
    - database
    - admin
    - server
    - tls
    ---------------------
    - [all]
```

### Setup - Database

```bash
> tdservice setup database [-force] --db-host=postgres.com --db-port=5432 --db-username=admin --db-password=password --db-name=tds_db
```
Environment variables `TDS_DB_HOSTNAME`, `TDS_DB_PORT`, `TDS_DB_USERNAME`, `TDS_DB_PASSWORD`, `TDS_DB_NAME` can be used instead of command line flags

### Setup - HTTP Server

```bash
> tdservice setup server --port=8443
```
Environment variable `TDS_PORT` can be used instead of command line flags

### Setup - TLS

```bash
> tdservice setup tls [--force] [--hosts=intel.com,10.1.168.2]
```

Creates a Self Signed TLS Keypair in /etc/tdservice/ for quality of life. It is expected that consumers of this product will provide their own key and certificate in /etc/threat-detection before or after running setup, to make `TDA` use those instead. 

Environment variable `TDS_TLS_HOSTS` can be used instead of command line flags

`--force` overwrites any existing files, and will always generate a self signed pair.


### Setup - Admin

```bash
> tdservice setup admin --admin-user=admin --admin-pass=password
```

Environment variable `TDS_ADMIN_USERNAME` and `TDS_ADMIN_PASSWORD` can be used instead

This task can be used to create multiple admin-users, but if a duplicate username is specified it will error out.

## Start/Stop

```bash
> tdservice start
  Threat Detection Service started
> tdservice stop
  Threat Detection Service stopped
```

## Uninstall

```bash
> tdservice uninstall [--keep-config]
  Threat Detection Service uninstalled
```
Uninstalls Threat Detection Service, with optional flag to keep configuration

## Help

```bash
> tdservice (help|-h|-help)
  Usage: tdservice <command> <flags>
    Commands:
    - setup
    - help
    - start
    - stop
    - status
    - uninstall
    - version
```

## Version

```bash
> tdservice version
    Threat Detection Service v1.0.0 build 9cf83e2
```


# Container Operations

A container can be started using the command:

`docker run isecl/tdservice:latest -e TDS_PORT=8443 -e ... -p 8443:8443`

Volume mounts for specifying the TLS cert files must be provided

Preferably, a docker-compose.yml would be used instead

```yaml
version: "3.2"

services:
  database:
    image: postgres:latest
    ...
  tds:
    image: isecl/tdservice:latest
    environment:
      TDS_PORT: 8443
    secrets:
      - source: tls.cert
        target: /run/secrets/tls.pem
    ...
    # NOT A COMPLETE DOCKER-COMPOSE EXAMPLE
```