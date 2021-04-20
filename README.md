# yiq

It's [jiq](https://github.com/fiatjaf/jiq) with [yq](https://github.com/mikefarah/yq).

You can drill down interactively by using [yq](https://github.com/mikefarah/yq) filtering queries.

yiq uses [yq](https://github.com/mikefarah/yq) internally, and it **requires** you to have `yq` in your `PATH`.

## Demo

T.B.D.

## Installation

Either [prebuilt binary for your system](https://github.com/fiatjaf/jiq/releases) (and make sure to `chmod +x` it first) or install/compile with Go:

```
go get github.com/zoetrope/yiq/cmd/yiq
```

If you don't have `yq` installed, follow instructions at https://github.com/mikefarah/yq/releases and make sure to put it in your `PATH`.

## Usage

### Quick start

* [simple example](#simple-example)

#### simple example

```
cat > sample.yaml <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample
  labels:
    app: ubuntu
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ubuntu
  template:
    metadata:
      labels:
        app: ubuntu
    spec:
      containers:
      - name: ubuntu
        image: ubuntu:18.04
EOF
```

```
cat sample.yaml | yiq
```
