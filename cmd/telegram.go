package cmd

import (
	"github.com/chazari-x/hmtpk_schedule/domain/telegram"
	"github.com/chazari-x/hmtpk_schedule/redis"
	"github.com/chazari-x/hmtpk_schedule/schedule"
	"github.com/chazari-x/hmtpk_schedule/storage"
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

			newRedis, err := redis.NewRedis(&cfg.Redis)
			if err != nil {
				log.Error(err)
				if !cfg.Log.Dev {
					return
				}
			}

			newStorage, db, err := storage.NewStorage(&cfg.DB, cmd.Context())
			if err != nil {
				return
			}
			defer func() {
				_ = db.Close()
			}()

			if err = telegram.Start(&cfg.Telegram, newRedis, schedule.NewSchedule(&cfg.Schedule, newRedis), newStorage); err != nil {
				log.Error(err)
				return
			}
		},
	}
	rootCmd.AddCommand(cmd)
	PersistentConfigFlags(cmd)
}
