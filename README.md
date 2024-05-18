# Faraway assignment, Dmitry Belov

## The task

Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge

## Running the app

### Server

```
docker build -t pow-server -f Dockerfile.server .
docker run --rm -it --network=host pow-server
```

### Client

```
docker build -t pow-client -f Dockerfile.client .
docker run --rm -it --network=host pow-client
```

## Running tests

```
go test ./...
```

## Choice of the PoW algorithm

Simple hash-based PoW algorithm based on SHA-256 was chosed for the following reasons:
* Resource bound. Clients must spend computational resources, which might make flood of requests economically inefficient for attackers.
* Adjustable difficulty. It makes it possible to increase the difficulty based on current load.
* It is a widely used algorithm with proved effectiveness.
* The algorithm is stateless, which reduces complexity and server costs.
* Easy to implement and understand.

## Possible improvements

Implementation improvements:
* Collect and visualize metrics:
    * Number of successful and failed requests.
    * Rate limit hits.
    * Resources utilization, such as CPU and memory.
* Instrument the workload to stream logs. Log warnings and errors; log less critical events if needed.
* Add more tests:
    * More unit tests (for example, cover middlewares).
    * Stress/performance tests, for example:
        * Mass GET.
        * Nonce flood.

Features:
* Extend PoW protocol:
    * Add protocol versioning.
    * It's possible to switch to an existing standard (for example, Hashcash).
* Implement more granular or advanced rate limiting. For example, based on IP adress instead of global.
* Allow a subset of clients (authenticated, or privileged) to skip the PoW challenge.
