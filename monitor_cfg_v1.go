package main

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

var gCfgMetricMain = map[string]func(namespace, name string, metric *CfgV1Metric) (func(), error){}

func monitorCfgV1(cfg *CfgV1) {
	for name, metric := range cfg.Metrics {
		main, ok := gCfgMetricMain[metric.Kind]
		if !ok {
			panic(fmt.Sprintf("metric kind %s not registered", metric.Kind))
		}
		routine := panicT(main(cfg.Namespace, name, &metric))
		go routine()
	}
}

func stripStdout(fn func(spec map[string]string, stdout string)) func(spec map[string]string, stdout string) {
	return func(spec map[string]string, stdout string) {
		stdout = strings.TrimSpace(stdout)
		fn(spec, stdout)
	}
}

func init() {
	gCfgMetricMain[`gauge`] = func(namespace, name string, metric *CfgV1Metric) (func(), error) {
		if len(metric.Matrix) != 0 {
			return nil, errors.NewPlain("gauge should not contain matrix")
		}
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      metric.Help,
		})
		panicIf(reg.Register(gauge))

		return (&commandExecutor{
			command:  metric.Command,
			interval: metric.Interval,
			matrix:   metric.Matrix,
			valueSetter: stripStdout(func(spec map[string]string, stdout string) {
				fp := panicT(strconv.ParseFloat(stdout, 64))
				gauge.Set(fp)
			}),
		}).exec, nil
	}

	gCfgMetricMain[`gauge_vec`] = func(namespace, name string, metric *CfgV1Metric) (func(), error) {
		if err := checkAndFillVecMatrix(metric.Matrix); err != nil {
			return nil, errors.WrapIff(err, "gauge_vec")
		}

		gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      metric.Help,
		}, mapKeys(metric.Matrix))
		panicIf(reg.Register(gaugeVec))

		return (&commandExecutor{
			command:  metric.Command,
			interval: metric.Interval,
			matrix:   metric.Matrix,
			valueSetter: stripStdout(func(spec map[string]string, stdout string) {
				fp := panicT(strconv.ParseFloat(stdout, 64))
				gaugeVec.With(spec).Set(fp)
			}),
		}).exec, nil
	}

	gCfgMetricMain[`counter_vec`] = func(namespace, name string, metric *CfgV1Metric) (func(), error) {
		if err := checkAndFillVecMatrix(metric.Matrix); err != nil {
			return nil, errors.WrapIff(err, "counter_vec")
		}

		keys := mapKeys(metric.Matrix)

		counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      name,
			Help:      metric.Help,
		}, keys)
		panicIf(reg.Register(counterVec))

		prevCounters := map[string]float64{}

		return (&commandExecutor{
			command:  metric.Command,
			interval: metric.Interval,
			matrix:   metric.Matrix,
			valueSetter: stripStdout(func(spec map[string]string, stdout string) {
				fp := panicT(strconv.ParseFloat(stdout, 64))

				var counterKeyBuilder strings.Builder
				for _, k := range keys {
					counterKeyBuilder.WriteString(spec[k])
					counterKeyBuilder.WriteByte(0)
				}
				counterKey := counterKeyBuilder.String()

				counterVec.With(spec).Add(fp - prevCounters[counterKey])
				prevCounters[counterKey] = fp
			}),
		}).exec, nil
	}
}
