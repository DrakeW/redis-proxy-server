# Redis Cache Proxy

Redis cache proxy is a simple read-through proxy that provides basic LRU caching capaility.

## How to use

To get the proxy server and a backing redis instance up and running (with docker containers). Do (Note: `docker`, `docker-compose` required)

```bash
make run
```

## Architecture

The project is composed of 3 main components:

1. Proxy server -- The proxy server is a simple HTTP server that exposes a single endpoint that translates the `GET` request from the client to a redis `GET` command.
    - if the key-value pair is in the cache of the proxy, the proxy server will load the value from cache and refresh its metadata
    - if the key-value pair is not yet cached or the cache has expired, the proxy server will translate and send the corresonding `GET` request to the backing redis instsance. After the result is retrieved successfully, the cache will be updated accordingly.

2. Cache -- A simple LRU cache implemented from scratch that can store any type of data other than just `string`, it exposes three methods:
    - `GET` - Loads a key-value pair from cache. The reading part of the operation is **not** an atomic operation to maximize concurrent read performance. And after the value is returned, the following update to the internal cache structure is an atomic operation that's protected by a `mutex` lock.
    - `Add` - Adds a key-value pair into cache. This operation is atomic and protected by a `mutex` lock. And the same key aleady exists in cache the entry that contains its value will be reused to minimize memory overhead of creating new entry.
    - `Remove` - Removes a key-value pair from the cache. This operation is also atomic and protected by a `mutex` lock. This method is not used by the server directly in this project and is mainly used for dealing with removing expired cached entry or cache eviction when capacity is reached. But it's still a public method for its future usage in, for example, cache invalidation.

3. Command line tool -- A simple command line tool that allows you to pass parameters to the server before it starts running.

More details are documented along with the code itself :)

## Algorithmic complexity of cache operation

`Get` - O(1)

`Add` - O(1)

`Remove` - O(1)

## How to run tests

Unit tests - `make test:unit`

Integration tests - `make test:integration`

## Time spent

It's hard to measure since the work is not done in a single chunk of time... so roughly

- research ~ 2 hrs (total)
- command line + proxy server ~ 1 hr
- cache - 1.5 hrs
- doc - 30min

## Some future TODOs

1. Currently logs are written to stdout by defualt, allow passing in a log file location and log level upon server starts.
2. RESP protocol -- didn't do it because it's kind of late in the night and it doesn't seem like something that'll help with my sleep :P
