package app

import (
	"github.com/go-seidon/local/internal/logging"
)

func NewDefaultLog(config Config) logging.Logger {
	opts := []logging.LogOption{}

	appOpt := logging.WithAppContext(config.AppName, config.AppVersion)
	opts = append(opts, appOpt)

	if config.AppDebug {
		debugOpt := logging.EnableDebugging()
		opts = append(opts, debugOpt)
	}

	if config.AppEnv == ENV_LOCAL || config.AppEnv == ENV_TEST {
		prettyOpt := logging.EnablePrettyPrint()
		opts = append(opts, prettyOpt)
	}

	stackSkipOpt := logging.AddStackSkip("github.com/go-seidon/local/internal/app")
	opts = append(opts, stackSkipOpt)

	return logging.NewLogrusLog(opts...)
}
