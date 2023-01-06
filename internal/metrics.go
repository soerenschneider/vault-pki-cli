package internal

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strings"
)

const (
	namespace = "vault_pki_cli"
)

var (
	MetricSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "success_bool",
		Help:      "Whether the tool ran successful",
	}, []string{"domain"})

	MetricCertExpiry = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "cert_expiry_seconds",
		Help:      "The date after the cert is not valid anymore",
	}, []string{"domain"})

	MetricCertLifetimeTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "cert_lifetime_seconds_total",
		Help:      "The total number of seconds this certificate is valid",
	}, []string{"domain"})

	MetricCertParseErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "cert_parse_errors_total",
		Help:      "The total number of parsing errors of a cert",
	}, []string{"domain"})

	MetricCertLifetimePercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "cert_lifetime_percent",
		Help:      "The passed lifetime of the certificate in percent",
	}, []string{"domain"})

	MetricRunTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "run_timestamp_seconds",
		Help:      "The date after the cert is not valid anymore",
	}, []string{"domain"})
)

func WriteMetrics(path string) error {
	log.Info().Msgf("Dumping metrics to %s", path)
	metrics, err := dumpMetrics()
	if err != nil {
		log.Info().Msgf("Error dumping metrics: %v", err)
		return err
	}

	err = os.WriteFile(path, []byte(metrics), 0644)
	if err != nil {
		log.Info().Msgf("Error writing metrics to '%s': %v", path, err)
	}
	return err
}

func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal().Msgf("can not start metrics server at %s: %v", addr, err)
	}
}

func dumpMetrics() (string, error) {
	var buf = &bytes.Buffer{}
	enc := expfmt.NewEncoder(buf, expfmt.FmtText)

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return "", err
	}

	for _, f := range families {
		// Writing these metrics will cause a duplication error with other tools writing the same metrics
		if !strings.HasPrefix(f.GetName(), "go_") {
			if err := enc.Encode(f); err != nil {
				log.Info().Msgf("could not encode metric: %s", err.Error())
			}
		}
	}

	return buf.String(), nil
}
