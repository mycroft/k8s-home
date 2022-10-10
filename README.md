# k8s-home

This repository builds & stores k8s configuration.

## Installation

Please start `fluxcd`:

```sh
flux bootstrap git \
  --url=ssh://git@git.mkz.me/mycroft/k8s-home.git \
  --private-key-file=/tmp/rsa \
  --branch=generated \
  --path=generated/
```
