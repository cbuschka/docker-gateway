# docker-gateway

### Exports docker containers labeled as ingress to the public

## Features

* nginx
* watches docker container changes via docker events
* generates nginx config for docker containers labeled with ingress-domain=www.example.com
* limits config updates to at max one every 5 seconds
* template for nginx config customizable

## Build

```
make build-with-docker build-container
```

## License

Copyright (c) 2022 by [Cornelius Buschka](https://github.com/cbuschka).

[MIT](./license.txt)
