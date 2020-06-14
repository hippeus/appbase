package main

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/hippeus/appbase/pkg/buildinfo"
	"github.com/hippeus/appbase/pkg/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configFile = pflag.String("config", "app.yaml", "relative path to the configuration file. Must include extension.")

	defaultConfigLookupPath = "/config/"
)

type Config struct {
	App     AppConfig
	HTTP    HTTPServerConfig
	Logging logger.Config
}

type AppConfig struct {
	Name string
	buildinfo.Build
}

type HTTPServerConfig struct {
	Host string
	Port uint32
	TLS  *TLSConfig
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

func parseConfigFlag(pConfFlag *string) (dir, file, ext string) {
	dir, file = filepath.Split(*pConfFlag)
	ext = filepath.Ext(file)
	file = strings.TrimSuffix(file, ext)
	if dir == "" {
		dir = "."
	}
	return dir, file, ext
}

func getConfig() Config {
	pflag.Parse()
	vp := viper.New()
	_ = vp.BindPFlag("config-file", pflag.Lookup("config"))
	dir, filename, ext := parseConfigFlag(configFile)
	vp.SetConfigName(filename)
	vp.SetConfigType(string(ext[1:]))
	vp.AddConfigPath(defaultConfigLookupPath)
	vp.AddConfigPath(dir)
	vp.SetEnvPrefix(filename)
	vp.AutomaticEnv()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	err := vp.ReadInConfig()
	if err != nil {
		cfgErr := &viper.ConfigFileNotFoundError{}
		if errors.As(err, cfgErr) {
			log.Println("🚨 Missing config file, using default configuration.... 🚨")
			cfg = Config{
				App: AppConfig{Name: "application"},
				HTTP: HTTPServerConfig{
					Host: "localhost",
					Port: 8080,
				},
				Logging: logger.DefaultConfig(),
			}
			return cfg
		} else {
			log.Fatalf("Fatal error config file: %s \n", err)
		}
	}

	if err = vp.Unmarshal(&cfg); err != nil {
		log.Fatalf("Fatal error while loading config: %s \n", err)
	}

	return cfg
}
