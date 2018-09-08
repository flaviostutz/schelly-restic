# schelly-restic
This exposes the common functions of Restic with Schelly REST APIs so that it can be used as a backup backend for Schelly (webhook).

See more at http://github.com/flaviostutz/schelly

# Usage

docker-compose .yml

```
version: '3.5'

services:

  restic-api:
    image: flaviostutz/schelly-restic
    ports:
      - 7070:7070
    environment:
      - RESTIC_PASSWORD=123
      - LOG_LEVEL=debug
```

# REST Endpoints

As in https://github.com/flaviostutz/schelly#webhook-spec
