package internal

import (
	"bytes"
	"crypto/x509"
	"errors"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/expfmt"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/vault-pki-cli/pkg"
)

const (
	metricsNamespace = "vault_pki_cli"
)

var (
	MetricSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "success_bool",
		Help:      "Boolean that reflects whether the tool ran successful",
	}, []string{"cn"})

	MetricCertExpiry = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cert_expiry_seconds",
		Help:      "The date after the cert is not valid anymore",
	}, []string{"cn"})

	MetricCertLifetimeTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cert_lifetime_seconds_total",
		Help:      "The total number of seconds this certificate is valid",
	}, []string{"cn"})

	MetricCertErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cert_errors_total",
		Help:      "The total number of errors while handling a cert",
	}, []string{"cn", "error"})

	MetricCertLifetimePercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cert_lifetime_percent",
		Help:      "The passed lifetime of the certificate in percent",
	}, []string{"cn"})

	MetricRunTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "run_timestamp_seconds",
		Help:      "The date after the cert is not valid anymore",
	}, []string{"cn"})
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

func UpdateCertificateMetrics(cert *x509.Certificate) {
	if cert == nil {
		log.Warn().Msg("can not update cert metrics, nil cert passed")
		return
	}

	secondsTotal := cert.NotAfter.Sub(cert.NotBefore).Seconds()
	MetricCertLifetimeTotal.WithLabelValues(cert.Subject.CommonName).Set(secondsTotal)
	secondsUntilExpiration := time.Until(cert.NotAfter).Seconds()

	percentage := math.Max(0, secondsUntilExpiration*100./secondsTotal)

	MetricCertExpiry.WithLabelValues(cert.Subject.CommonName).Set(float64(cert.NotAfter.UnixMilli()))
	MetricCertLifetimePercent.WithLabelValues(cert.Subject.CommonName).Set(percentage)
}

func TranslateErrToPromLabel(err error) string {
	if errors.Is(err, pkg.ErrWriteCert) {
		return "write_issued_cert"
	}
	if errors.Is(err, pkg.ErrNoCertFound) {
		return "no_cert_found"
	}
	if errors.Is(err, pkg.ErrIssueCert) {
		return "issue_error"
	}
	if errors.Is(err, pkg.ErrCertInvalidData) {
		return "issued_cert_invalid_data"
	}
	return "unknown"
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
