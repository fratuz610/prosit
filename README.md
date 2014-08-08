prosit
======

A simple process manager in Go, suitable for golang as well as node apps

It does what I need (restarts, notifications, queriable logs) whereas pm2 doesn't

## Features

- run a process with a specific user and make sure it's always running
- supports alerts via email (mailgun supported only for now - register it's free)
- keep the process stdout and stderr logs in memory and autorotates them every 10k lines
- offers a complete restful API 
- offers a command line API
- the main executable is both a service (if called with &) and a client
- restart a process if it dies and send an email if an alert is associated with the process
- it doesn't suck (btw)
- compiled with Go 1.3.1

## Install

TODO

## Restful APIs

TODO

## Command line Interface

The command line interface is designed to be as helpful and predictive as possible
No command line parameters only q/a approach. i.e.:

```
todo
```

It's implemented as a wrapper for the restful apis


## TODO

- Support gmail for alerts
- Support 
