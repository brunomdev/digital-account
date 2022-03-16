package newrelic

import (
	"github.com/brunomdev/digital-account/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
)

func NewNewRelic(cfg *config.Config) (*newrelic.Application, error) {
	option := []newrelic.ConfigOption{
		newrelic.ConfigAppName(cfg.NewRelicAppName),
		newrelic.ConfigLicense(cfg.NewRelicLicenseKey),
		newrelic.ConfigEnabled(cfg.NewRelicLicenseKey != ""),
	}

	if cfg.AppDebug {
		option = append(option, newrelic.ConfigDebugLogger(os.Stdout))
	}

	return newrelic.NewApplication(option...)
}
