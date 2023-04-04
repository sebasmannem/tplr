package internal

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/sebasmannem/tplr/pkg/tekton_handler"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

/*
 * This module reads the config file and returns a config object with all entries from the config yaml file.
 */

const (
	envDebug     = "TPLR_DEBUG"
	envNameSpace = "TPLR_NAMESPACE"
	envPipeline  = "TPLR_PIPELINE"
	envTimeout   = "TPLR_TIMEOUT"
	configFile   = "/etc/tplr/tplr.yaml"
)

var (
	debug     bool
	version   bool
	pipeline  string
	namespace string
	timeout   string
)

func ProcessFlags() (err error) {
	if os.Getenv(envDebug) != "" {
		debug = true
	}
	flag.BoolVar(&debug, "d", debug, "Add debugging output")
	flag.BoolVar(&version, "v", false, "Show version information")
	flag.StringVar(&pipeline, "p", os.Getenv(envPipeline), "Pipeline to run")
	flag.StringVar(&namespace, "n", os.Getenv(envNameSpace), "Namespace to find the pipeline to run")
	flag.StringVar(&timeout, "t", os.Getenv(envTimeout),
		"Timeout for the pipeline")

	flag.Parse()

	if version {
		//nolint
		fmt.Println(appVersion)
		os.Exit(0)
	}

	return err
}

type Config struct {
	Tekton  tekton_handler.Config `yaml:"tekton"`
	Debug   bool                  `yaml:"debug"`
	Timeout string                `yaml:"pipeline_timeout"`
}

func (c Config) GetTimeoutContext(parentContext context.Context) (context.Context, context.CancelFunc) {
	if c.Timeout == "" {
		return parentContext, nil
	}
	lockDuration, err := time.ParseDuration(c.Timeout)
	if err != nil {
		log.Fatal(err)
	}
	return context.WithTimeout(parentContext, lockDuration)
}

func NewConfig() (config Config, err error) {
	if err = ProcessFlags(); err != nil {
		return
	}
	if _, err = os.Stat(configFile); !errors.Is(err, os.ErrNotExist) {
		var yamlConfig []byte
		yamlConfig, err = os.ReadFile(configFile)
		if err != nil {
			return config, err
		} else if err = yaml.Unmarshal(yamlConfig, &config); err != nil {
			return config, err
		}
	}
	if namespace != "" {
		config.Tekton.Namespace = namespace
	}
	if pipeline != "" {
		config.Tekton.Pipeline = pipeline
	}
	if timeout != "" {
		config.Tekton.Timeout = timeout
	}
	if debug {
		config.Debug = true
	}
	return config, nil
}
