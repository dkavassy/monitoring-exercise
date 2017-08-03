Simple monitoring service
=========================

## How to use

Run webserver instances:
`FLASK_APP=webserver.py python -m flask run --port 2000`
`FLASK_APP=webserver.py python -m flask run --port 3000`
`FLASK_APP=webserver.py python -m flask run --port 4000`
`FLASK_APP=webserver.py python -m flask run --port 5000`

Run the request generator

`./request.sh`

Run the monitoring service:

`go build monitoring.go`
`./monitoring --interval 5 --targets
localhost:5000,localhost:4000,localhost:3000,localhost:2000 --metrics
requests,errors`

## Comments

Main areas lacking:
  - error handling in monitoring service
  - tests

## Directory structure

server/ - instrumented web server
client/ - monitoring service
request.sh - request generator
