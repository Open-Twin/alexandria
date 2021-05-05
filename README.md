# alexandria
A simple distributed, redundant, error correcting DNS server

## Start the DNS

To start a single DNS node simply direct youself to **alexandria/deployments/** then open the terminal in this path, start the node with:
```shell
docker-compose up
```
In the .env it is possible to set certain settings like the addresses and ports an example would be:
```config
RAFT_ADDR=10.5.0.2
HTTP_ADDR=0.0.0.0
RAFT_PORT=7000
HTTP_PORT=8000
```
In the following section about the **Raft demo** it can be witnessed that a another docker-compose.yml will be started to test the raft functionality.

## Raft demo
Author: Sebastian Bruckner-Hrubesch,
Last updated: 09.02.2021

The demo consists of a leader and two followers.
The open ports are:
* leader: 8000
* follower1: 8001
* follower2: 8002
### Start container
To use it, just head to the ``/deployments/dev`` directory and start the container:
```shell
docker-compose -f docker-compose.dev.yml up
```
The leader and the followers should now be up and running.

### Interact with the cluster via cURL
To interact with the cluster, a JSON message has to be sent. It looks like this:
```json
{
  "service": "temperature",
  "ip": "1.2.3.4",
  "type": "store",
  "key": "toilet",
  "value": "100C"
}
```
The supported types are ``store, update, delete, get``.

In order to change data (store, update, delete), send a POST request with your JSON to a server in the cluster:
```shell
curl -XPOST -d @testreq.json -H "Content-Type: application/json" http://127.0.0.1:8000/key
```

In order to retrieve data from the cluster (get), send a GET request:
```shell
curl -XGET -d @testreq.json -H "Content-Type: application/json" http://127.0.0.1:8001/key
```
For retrieving data, the value does not have to be in the JSON.

## Testing

Executing the tests can be done in two different ways.
Firstly it is possible to choose a certain test.go file and see its outcome by using the following command in the cmd:
```shell
go test file_test.go
```
Important to mention is that you'd need to be in the right directory otherwise you'd need to enter the path in the command like following:
```shell
go test src/tests/file_test.go
```
Lastly it is also possible to run all the tests at once by using:
```shell
go test ./...
```

It should be mentioned that it is also possible to simply run the tests in the IDE, the IDE that we used was GoLand where its pretty simple by just right clicking the file_test.go and pressing "run file_test.go".
