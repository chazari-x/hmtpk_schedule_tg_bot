package cmd

import (
	"fmt"
	"github.com/chazari-x/hmtpk_schedule/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

var configFile = "etc/config.yaml"

type Config struct {
	Telegram config.Telegram `yaml:"bot"`
	Redis    config.Redis    `yaml:"redis"`
	DB       config.DataBase `yaml:"db"`
	Schedule config.Schedule `yaml:"schedule"`
	Log      config.Log      `yaml:"log"`
}

func getConfig(cmd *cobra.Command) (*Config, error) {
	var cfg Config

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		PadLevelText:              true,
		EnvironmentOverrideColors: true,
	})

	file, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, fmt.Errorf("get flag err: %s", err)
	}

	if file != "" {
		configFile = fmt.Sprintf("etc/config.%s.yaml", file)
	}

	f, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("open config file \"%s\": %s", configFile, err)
	}

	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config file: %s", err)
	}

	level, err := log.ParseLevel(cfg.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("parse level err: %s", err)
	}

	if cfg.Log.Level == "" {
		cfg.Log.Level = "trace"
	}
	log.SetLevel(level)

	return &cfg, nil
}

func PersistentConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("config", "", "dev")
}
