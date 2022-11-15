package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

var (
	reg               = prometheus.NewRegistry()
	gPredefinedMatrix map[string][]string
)

func main() {
	cfg := flag.String("config", "config.yaml", "config file for script_exporter")
	addr := flag.String("addr", ":9090", "export address")
	matrix := flag.String("matrix", "", "matrix config. format key1=val1,val2:key2=val3,val4")
	flag.Parse()

	gPredefinedMatrix = parsePredefinedMatrix(*matrix)

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
