# Monitor API Golang with Grafana, Prometheus, Loki

![](https://github.com/toannd96/go-monitor-grafana/blob/master/grafana.png)

![](https://github.com/toannd96/go-monitor-grafana/blob/master/promtail.png)

An example of using Loki with Grafana for Golang app monitoring.

## Running
```bash
make run-docker

make run-app
```

After that, you could go to `http://127.0.0.1:3000` and set up Loki data source in the Grafana UI.
