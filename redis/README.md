# Redis Clone

https://codingchallenges.fyi/challenges/challenge-redis

Data is stored in `data/data.rdb`. Note that this format is **not** strictly following the standard RDB format. Some changes have been made for easier implementation.

```sh
# prints the make commands
make

# be sure to install redis-cli first
# enter the REPL (note: this will cause some un-implemented commands to be sent to the server, but it's fine)
redis-cli

# send single commands
# redis-cli <COMMAND>
redis-cli incr k
redis-cli save
```

## Benchmarks

Benchmarks done on M1 macbook air.

tl;dr: my implementation achieves ~80% performance of actual redis, which I think isn't too bad considering I didn't try to optimise it much.

### this implementation

Benchmarked with `PANIC` level logs (aka almost no logs).
Enabling `TRACE` level logs slows the program down almost 10x...

```sh
❯ redis-benchmark -t set,get, -n 100000 -q
WARNING: Could not fetch server CONFIG
SET: 137741.05 requests per second, p50=0.175 msec
GET: 142653.36 requests per second, p50=0.175 msec
```

### redis-server

Redis version = 7.2.4

```sh
❯ redis-benchmark -t set,get, -n 100000 -q
SET: 162601.62 requests per second, p50=0.167 msec
GET: 168067.22 requests per second, p50=0.167 msec
```
