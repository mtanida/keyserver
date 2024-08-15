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

## Monitoring
We can monitor the keyserver service using Prometheus or Victoriametrics ServiceMonitor. The ServiceMonitor is implemented with a Kubernetes Operator that monitors services with a specific label. In addition, we can add Grafana and set up Dashboards for visualizing the health of the service. For alerting, we can make use of Prometheus/Victoriametrics Alertmanager to send an alert under certain conditions. Alerts can be integrated with PagerDuty, Slack, SendGrid, etc. to send timely alerts.
