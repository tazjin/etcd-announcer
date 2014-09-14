# Generic etcd announcer

This project is a tiny, scratch-based Docker container that simply announces a service to `etcd` based on some arguments that were passed in.

There is support for announcing the IP address of a specific interface, as well as passing in a custom string to announce.

Announcements are done with a short-lived TTL (by default 30 seconds), with a new announcement happening every TTL/2 seconds.

## Usage

There is a build of this on [Docker Hub](https://registry.hub.docker.com/u/tazjin/etcd-announcer/), so you can run the examples directly. Skip to the `Building` section if you don't trust my build.

By default, `etcd-announcer` expects an `etcd` instance at `127.0.0.1:4001`. In order to have access to the local etcd instance as well as to be able to fetch interface IP addresses, you can run it with `--net="host"`.

```sh
# Announce a service called "foo" with the content "bar" and a TTL of 1 minute
docker run --net="host" tazjin/etcd-announcer -service "foo" -msg "bar" -ttl 60

# Announce a service called "web" using the IP of interface 'eth0'
docker run --net="host" tazjin/etcd-announcer -service "web" -type "net" -if "eth0"

# Announce 'eth0'-IP for a service called "riak-1" under the "/service/riak" directory
docker run --net="host" tazjin/etcd-announcer -service "riak-1" -path "/service/riak" -type "net" -if "eth0"

# Announce the host's FQDN for a service called "foo", speaking to etcd through the Docker0 interface
DOCKER_HOST_IP=$(ip -o -4 addr show docker0 | awk '{print $4}' | cut -d/ -f1)
ETCD_ADDR="http://${DOCKER_HOST_IP}:4001"
docker run tazjin/etcd-announcer -service "foo" -msg $(hostname --fqdn) -etcd "${ETCD_ADDR}"
```

## Building

There is a simple Makefile for this project that performs all build steps in Docker. The only dependency for building this is having Docker running locally as well as some disk space for the [google/golang](https://registry.hub.docker.com/u/google/golang/) Docker image.

```sh
# Just build the binary
make gobuild

# Build Docker image (builds with tag tazjin/etcd-announcer by default)
make docker

# Build both
make

# Clean up
make clean
```

## TODO

* Add support for interfaces with multiple ID addresses (currently this just takes the first one)
* Clean this up (I'm not a Go programmer and just hacked this together quickly)
