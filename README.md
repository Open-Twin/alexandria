# alexandria
A simple distributed, redundant, error correcting DNS server

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
