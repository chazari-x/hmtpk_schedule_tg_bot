package cmd

import (
	"github.com/chazari-x/hmtpk_schedule/domain/telegram"
	"github.com/chazari-x/hmtpk_schedule/redis"
	"github.com/chazari-x/hmtpk_schedule/schedule"
	"github.com/chazari-x/hmtpk_schedule/selenium"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "telegram",
		Short: "telegram",
		Long:  "telegram",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := getConfig(cmd)
			if err != nil {
				log.Error(err)
				return
			}

			log.Trace("telegram starting..")
			defer log.Trace("telegram stopped")

			r, err := redis.Redis(&cfg.Redis)
			if err != nil {
				log.Error(err)
				if !cfg.Log.Dev {
					return
				}
			}

			newSelenium, s, err := selenium.NewSelenium()
			if err != nil {
				log.Error(err)
				return
			}

			defer s.Quit()

			if err = telegram.Start(&cfg.Telegram, r, schedule.NewSchedule(&cfg.Schedule, r, newSelenium)); err != nil {
				log.Error(err)
				return
			}
		},
	}
	rootCmd.AddCommand(cmd)
	PersistentConfigFlags(cmd)
}
