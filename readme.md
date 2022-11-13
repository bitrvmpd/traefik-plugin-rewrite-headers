# Rewrite Headers

Inspired by: github.com/XciD/traefik-plugin-rewrite-headers

Rewrite headers is a middleware plugin for [Traefik](https://traefik.io) which replace a header in the request and/or response

## Configuration

### Static

```yaml

experimental:
  plugins:
    rewriteHeaders:
      modulename: "github.com/bitrvmpd/traefik-plugin-rewrite-headers"
      version: "v0.0.1"
```

### Dynamic

To configure the Rewrite Request Header plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/).
The following example creates and uses the rewriteHeaders middleware plugin to modify the Location header

```yaml
http:
  routes:
    my-router:
      rule: "Host(`localhost`)"
      service: "my-service"
      middlewares : 
        - "rewriteHeaders"
  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1"
  middlewares:
    rewriteHeaders:
      plugin:
        rewriteHeaders:
          rewrites:
            request:
              - header: Location
                regex: "^http://(.+)$"
                replacement: "https://$1"
            response:
              - header: Location
                regex: "^http://(.+)$"
                replacement: "https://$1"
```

Label based configuration

``` yaml
- traefik.http.middlewares.rewriteHeaders.plugin.rewriteHeaders.rewrites.request[0].header = Location
- traefik.http.middlewares.rewriteHeaders.plugin.rewriteHeaders.rewrites.request[0].regex = ^http://(.+)$
- traefik.http.middlewares.rewriteHeaders.plugin.rewriteHeaders.rewrites.request[0].replacement = https://$1
- traefik.http.middlewares.rewriteHeaders.plugin.rewriteHeaders.rewrites.response[0].header = Location
- traefik.http.middlewares.rewriteHeaders.plugin.rewriteHeaders.rewrites.response[0].regex = ^http://(.+)$
- traefik.http.middlewares.rewriteHeaders.plugin.rewriteHeaders.rewrites.response[0].replacement = https://$1
```
