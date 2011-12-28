## cstore

**Warning: This is a random toy, not a production-ready service.  It will
eat your data.  It stores all data in RAM, and uses the same Redis backend
for unit tests and standalone servers.  This repository may go away.**

The `cstore` server is a distributed, content-based storage system.  Like
`git`, it stores binary blobs indexed by cryptographic hash codes.  The
files are read and written using a simple REST API.

The storage servers keep a global index in Redis, which they use to
internally forward HTTP requests to other servers.  This means that you can
access any file on any server.

    $ echo -n 'Hello, world!' | sha256sum
    315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3
    $ HASH=315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3

In two other terminals, start some cstore servers:

    ./cstore 127.0.0.1:12345
    ./cstore 127.0.0.1:12346

For now, these will use the default Redis store on the local host.  Next,
add a file to the store:

    curl -i -X PUT -d 'Hello, world!' "http://localhost:12345/$HASH"

...and fetch it from the other store:

    curl "http://localhost:12346/$HASH"

### How it works

See REDIS.md for a description of how the servers communicate.

### Future ideas

In the future, I might add support for POST:

    curl -i -X POST -d 'Hello, world!' http://localhost:12345/

This would return a 'Location:' header with the SHA256 sum.
