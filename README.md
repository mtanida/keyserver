# keyserver
random key generating API service


## Running unit tests
Run unit tests with
```
go test
```

## Building Docker image
To build the docker container image
```
docker build -t keyserver .
```
To run the docker container:
```
docker run -p 1123:1123 keyserver:latest
```

## Installing Helm chart
The helm chart with:
```
helm install -n <namespace> <release name> ./helm/keyserver [--set "max_key_size=1024"] [--set "service.port=1123"]
```
