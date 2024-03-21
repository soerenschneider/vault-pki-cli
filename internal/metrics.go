package internal

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	"github.com/rs/zerolog/log"
)

const (
	metricsNamespace = "vault_pki_cli"
)

var (
	MetricSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "success_bool",
		Help:      "Whether the tool ran successful",
	}, []string{"domain"})

	MetricCertExpiry = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cert_expiry_seconds",
		Help:      "The date after the cert is not valid anymore",
	}, []string{"domain"})

	MetricCertLifetimeTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cert_lifetime_seconds_total",
		Help:      "The total number of seconds this certificate is valid",
	}, []string{"domain"})

	MetricCertParseErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cert_parse_errors_total",
		Help:      "The total number of parsing errors of a cert",
	}, []string{"domain"})

	MetricCertLifetimePercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cert_lifetime_percent",
		Help:      "The passed lifetime of the certificate in percent",
	}, []string{"domain"})

	MetricRunTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
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

	return os.WriteFile(path, []byte(metrics), 0644) // #nosec G306
}

func StartMetricsServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:              addr,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func dumpMetrics() (string, error) {
	var buf = &bytes.Buffer{}
	fmt := expfmt.NewFormat(expfmt.TypeTextPlain)
	enc := expfmt.NewEncoder(buf, fmt)

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return "", err
	}

	for _, f := range families {
		// Writing these metrics will cause a duplication error with other tools writing the same metrics
		if strings.HasPrefix(f.GetName(), metricsNamespace) {
			if err := enc.Encode(f); err != nil {
				log.Info().Msgf("could not encode metric: %s", err.Error())
			}
		}
	}

	return buf.String(), nil
}
