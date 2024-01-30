# Spot

## Parking space tracking service
Spot is a web service with exposed REST API for tracking parking spots
in public garages. It was developed as the final BSc
studies assignment before graduation.

## Parking monitor on Raspberry Pi
Monitor is implemented using Go programming language and a RaspberryPi, with a couple of LEDs and
an ultrasonic distance sensor. For more information about the device,
check the [rpi](./rpi/monitor) directory.

## Prerequisites
* Linux distribution
* Go programming language v1.10 or newer
* `dep` - Go dependency management tool
* MongoDB database server

## Programmer's deployment
Assuming `GOPATH` is set correctly and MongoDB server is running, execute following commands:
```bash
$ cd $GOPATH/src
$ go get -v github.com/cicovic-andrija/spot
$ cd github.com/cicovic-andrija/spot
$ dep ensure
$ make deploy
```

##  REST API Overview
| Operation  | Request |
| :--- | :--- |
| Create a garage | `POST /v1/garages {"name": "Union Sq. Garage", "city": "San Francisco", "address": "333 Post Street", "geolocation": {"longitude": -122.40754, "latitude": 37.788062}}` |
| Get all garages' properties  | `GET /v1/garages` |
| Get garage properties | `GET /v1/garages/{id}` |
| Change garage properties | `PUT /v1/garages/{id} {"name": "Union Square Garage", "city": "San Francisco"}` |
| Delete a garage | `DELETE /v1/garages/{id}` |
| Create a garage section | `POST /v1/garages/{id}/sections {"name": "A", "level": "Ground", "description": "Regular parking space", "total_spots": 42}` |
| Get all sections' properties | `GET /v1/garages/{id}/sections` |
| Get section properties | `GET /v1/garages/{id}/sections/{name}` |
| Change section properties | `PUT /v1/garages/{id}/sections/{name} {"name": "A1", "total_spots": 10}` |
| Delete a section | `DELETE /v1/garages/{id}/sections/{name}` |
| Update parking spot status (connect device) | `POST /v1/garages/{id}/sections/{name}/actions {"action": "update", "params": [{"number": 1, "label": "A1-1", "taken": false}]}` |
| Parking spot status: bulk update | `POST /v1/garages/{id}/sections/{name}/actions {"action": "update", "params": [{"number": 2, "label": "A1-2", "taken": false}, {"number": 3, "label": "A1-3", "taken": true}, {"number": 4, "label": "A1-4", "taken": false}]}` |
| Disconnect device | `POST /v1/garages/{id}/sections/{name}/actions {"action": "disconnect", "params": [{"number": 1}]}` |
| Disconnect device - bulk | `POST /v1/garages/{id}/sections/{name}/actions {"action": "disconnect", "params": [{"number": 2}, {"number": 3}, {"number": 4}]}` |
| Shutdown the server | `POST /v1/control {"action": "shutdown"}` |

## Running tests

In order to run tests, assuming the current working directory is the project's root, run:

```bash
$ cd test
$ go test -v ./...
```

