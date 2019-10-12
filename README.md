# OnError

## Description
This application runs a command and captures its stdout and stderr. If the command exits successfully, the command output is sent to syslog. If the command fails, then the output is sent to stderr and syslog.

## Building

```
git clone https://github.com/rsubr/onerror.git
cd onerror
go build onerror.go
strip onerror
```

## Usage

Use this application in crontab where you want to be emailed only if the cron job fails.

```
# Sample /etc/crontab
# Cron sends email every hour irrespective of borg backup success
# or failure
@hourly root /usr/local/sbin/borg-backup-foo.sh


# Cron sends email only if the backup script failed 
@hourly root onerror /usr/local/sbin/borg-backup-bar.sh
```

When calling bash scripts from `onerror`, you must use `set -e pipefail` to ensure the script fails and exits immediately whenever any command fails.

```
#!/bin/bash
# Sample bash script

# do not continue script on errors
set -euo pipefail

pg_dump ...
borg create $REPO::foo-{now} $MY_FILES
```
