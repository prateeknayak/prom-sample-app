# Prometheus with Sample App

## How to use the compose setup

1. Clone the repo on you machine
1. Run `make prom-up` command to create stack. This should standup
   1. Local prometheus on http://localhost:9090
   1. Grafana on http://localhost:3000
   1. Sample App on http://localhost:8080
   1. Push Gateway on http://localhost:9091
1. The Sample App as 4 endpoints
   1. `/`: returns random response with random status code
   1. `/metrics`: returns prometheus compatible metrics from the `sample-app`
   1. `/e500`: returns random ~500 error
   1. `/push`: endpoint that accepts {"push_url":"<URL_HERE>"} as body and either pushes or delete the metrics based on method
      1. use curl -X POST -d {"push_url":"http://pushgateway:9091"} http://localhost:9091/push to push metric
      1. use curl -X DELETE -d {"push_url":"http://pushgateway:9091"} http://localhost:9091/push to delete metric
   1. use `make prom-down` to destory the stack
   1. use `make sloth-gen` to generate SLO based on a sloth spec located in `compose/sloth` directory
   1. download hey ( `brew install hey` ) and use it to test sample app endpoints and generate metrics
