# legs-client

legs-client provides the cli command tool and client daemon for legs product.

## Installation
Please execute below command. install.sh downloads client tool to current directory.
```bash
$ curl --silent https://raw.githubusercontent.com/iij/legs-client/master/install.sh | sh
```
After installation, you can show version info by `version` command.
```bash
$ ./legsc version
```
And if you need, copy tool binary to `/usr/local/bin`
```bash
$ sudo cp ./legsc /usr/local/bin/legsc
```

## Basic Usage
```bash
# show help
$ ./legsc

# set secret key
$ ./legsc secret <your secret key>

# export config file with all current/default values
$ ./legsc export

# specify config file path
$ ./legsc -c path/to/config.yml secret <your secret key>

# start client in foreground
$ ./legsc start -f

# start client in background
$ ./legsc start

# stop client
$ ./legsc stop

# send data to server
$ ./legsc send routing/name '{"value": 1}'
```

## Development
```bash
# get go libraries which use in project.
$ make setup

# dep ensure
$ make dep

# start daemon with localconfig and tail log file.
$ make run

# stop daemon
$ make stop

# format by goimports
$ make fmt
```
