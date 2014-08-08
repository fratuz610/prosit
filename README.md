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

Also each write operation to stderr trigger a notification email

The only mail gateway supported at the moment is Mailgun. Register, it's free.

### Process management

Prosit spawns a child process (tutor) with each process start request
This child process changes its uid and spawns the requested process (which inherits the uid)

i.e.

  prosit -> prosit (tutor) -> node server.js
  <root>    <userA>          <userA>

The tutor process 
- starts up
- spawns the target process
- redirects the target process stderr and stdout to its stderr and stdout streams
- wait for the target process to exit
- exits with the target process' exit code

## Install

TODO

## Run (as a service)

Prosit needs to run as root in this configuration. This is required to call setuid on the tutor process

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
Only one command line parameter is needed (the command name) everything else is interactive

If you need programmatic access use the RESTFul apis. Believe me: it's better that using bash and grep the output.

```
todo
```

It's implemented as a wrapper for the RESTFul apis


## TODO

- support a different http port
- support http authentication
- support https (trivial but really needed?)
- Support gmail for alerts
- Make the log rotation size modifiable
- Improve command line output formats

## License

MIT license. Go nuts.
