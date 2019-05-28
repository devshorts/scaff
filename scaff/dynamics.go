package scaff

import (
	"github.com/devshorts/scaff/scaff/config"
)

func ResolveDynamics(conf config.ScaffConfig, bag map[string]string) map[string]string {
	b := make(map[string]string)

	for k, v := range bag {
		b[k] = v
	}
	runner := NewRuleRunner(b)

	for key, value := range conf.Dynamics {
		b[key] = runner.Replace(value, DEFAULT_DELIM)
	}

	return b
}
