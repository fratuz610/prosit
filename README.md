prosit
======

A simple process manager in Go, suitable for golang as well as node apps

It does what I need (restarts, notifications, queryable logs that autotate) whereas pm2 doesn't

## Features

- run a process with a specific user and make sure it's always running
- supports alerts via email (mailgun supported only for now - register it's free)
- keep the process stdout and stderr logs in memory and autorotates them every 10k lines
- offers a complete RESTFul API 
- offers a command line API
- the main executable is both a service (if called with &) and a client
- restart a process if it dies and send an email if an alert is associated with the process
- alerts
- it doesn't suck (btw)
- compiled with Go 1.3.1

### Alerts

Alerts are what every ops guy need. Node apps can crash at any time (like Chrome) for no reason.
Being notified when this happens is useful

Also each write to stderr trigger a notification email

The only mail gateway supported at the moment is Mailgun. Register, it's free.

### Process management

Prosit spawns a child process with each process start request

This child process changes its uid and spawns the requested process (which inherits the uid)

## Install

TODO

## Run (as a service)

Prosit needs to run as root in this configuration. This is required to call setuid

```
prosit &
```

TODO

## Run (as a CLI)

TODO

## RESTFul APIs

TODO

## Command line Interface

The command line interface is designed to be as helpful and predictive as possible
No command line parameters only q/a approach. i.e.:

```
todo
```

It's implemented as a wrapper for the RESTFul apis


## TODO

- Support gmail for alerts
- Make the log rotation size modifiable
- Improve command line output formats