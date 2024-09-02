package employeside

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"

	cloudconfig "github.com/realbucksavage/spring-config-client-go/v2"

	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

type Configuration struct {
	DB      DatabaseConfiguration `yaml:"db" validate:"required"`
	Bind    BindingConfiguration  `yaml:"bind" validate:"required"`
	Profile string                `yaml:"-"`
}

type BindingConfiguration struct {
	HTTP string `yaml:"http" validate:"required,tcp_addr"`
}

type DatabaseConfiguration struct {
	Database string `validate:"required" yaml:"database"`
	Username string `validate:"required" yaml:"username"`
	Password string `validate:"required" yaml:"password"`
	Hostname string `validate:"required" yaml:"hostname"`
	Port     int    `validate:"required" yaml:"port"`
}

func (cfg DatabaseConfiguration) CreateDSN() string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Hostname, cfg.Port, cfg.Database)
}

func LoadConfiguration() Configuration {

	var (
		profile    string
		configFile string

		configServer   = os.Getenv("CONFIG_SERVER")
		configUsername = os.Getenv("CONFIG_USERNAME")
		configPassword = os.Getenv("CONFIG_PASSWORD")
	)

	klog.InitFlags(nil)
	flag.StringVar(&profile, "profile", "development", "The profile to run the app in")
	flag.StringVar(&configFile, "cfg", "./config.yaml", "The configuration file")
	flag.Parse()

	klog.EnableContextualLogging(true)

	var cfg Configuration
	if configServer == "" {
		file, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}

		defer file.Close()

		if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
			panic(err)
		}
	} else {
		cfgClient, err := cloudconfig.NewClient(
			configServer,
			"employe-side",
			profile,
			cloudconfig.WithBasicAuth(configUsername, configPassword),
			cloudconfig.WithFormat(cloudconfig.YAMLFormat),
		)
		if err != nil {
			panic(err)
		}

		if err := cfgClient.Decode(&cfg); err != nil {
			panic(err)
		}
	}

	err := validator.New().Struct(cfg)
	var vErr validator.ValidationErrors
	if errors.As(err, &vErr) {
		panic(vErr)
	}

	cfg.Profile = profile
	return cfg
}
