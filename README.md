# `servecert`

A small local HTTPS proxy using [`mkcert`](https://github.com/FiloSottile/mkcert) certificates. Please use with care.

Usage:

```shell
git clone https://github.com/lgarron/servecert && cd servecert
make build

# Proxy to https://localhost/
build/servecert https://example.com /
```

You could also also use any of the following as the second argument:

| Argument                        | Meaning                                   |
| ------------------------------- | ----------------------------------------- |
| `foo.test`                      | `https://foo.test/`                       |
| `foo.test:8000/`                | `https://foo.test:8000/`                  |
| `/subpath`                      | `https://localhost/subpath/`              |
| `http://localhost/`             | (serves HTTP without a cert on port 80)   |
| `http://foo.test:8000/subpath/` | (serves HTTP without a cert on port 8000) |

When you proxy a new domain using HTTPS for the first time, this invokes `mkcert` to install a root certificate and generates a certificate. [Learn more about `mkcert` certificate installation](https://github.com/FiloSottile/mkcert#installation).
