# Monitor API Golang with Grafana, Prometheus, Loki

![](https://github.com/toannd96/go-monitor-grafana/blob/master/grafana.png)

![](https://github.com/toannd96/go-monitor-grafana/blob/master/promtail.png)

An example of using Loki with Grafana for Golang app monitoring.

## Running
```bash

cd ./deploy
docker-compose up -d

cd ..
go run ./cmd/main.go

```

After that, you could go to `http://127.0.0.1:3000` and set up Loki data source in the Grafana UI.
