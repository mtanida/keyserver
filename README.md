# keyserver
random key generating API service


# Running unit tests
Run test with
```
go test
```

# Building Docker image
To build the docker image
```
docker build .
```

# Installing Helm chart
The helm chart with:
```
helm install -n <namespace> ./helm/keyserver  [--set "max_key_size=1024"] [--set "service.port=1123"]
```
