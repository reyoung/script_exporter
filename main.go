package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

var reg = prometheus.NewRegistry()

func main() {
	cfg := flag.String("config", "config.yaml", "config file for script_exporter")
	addr := flag.String("addr", ":9090", "export address")
	flag.Parse()

	file := panicT(os.Open(*cfg))
	defer func() {
		panicIf(file.Close())
	}()
	decoder := yaml.NewDecoder(file)
	var config Config
	panicIf(decoder.Decode(&config))
	monitor(config)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	_ = http.ListenAndServe(*addr, nil)
}
