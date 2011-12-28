All keys go in the namespace `cstore:` to prevent collisions.

## Getting a unique server ID

    id = INCR cstore:server_id_generator

## Heartbeat

    SETX cstore:server:$id 20 $host_and_port

## Registering a blob

    LPUSH cstore:blob:$digest $id

## Finding a blob

    ids = SMEMBERS cstore:blob:$digest
    addrs = MGET cstore:server:$ids[N]...

Try http://$addr[N]/$digest in sequence until we get a hit.
