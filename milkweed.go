package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jcc333/milkweed/poasts"
	"github.com/jcc333/milkweed/state"

	_ "github.com/joho/godotenv/autoload"
)

// The configuration for milkweed consists of:
type Config struct {
	// The BlueSky Username
	username string

	// The BlueSky Password
	password string

	// The URL of the RSS feed endpoint
	rss string

	// The path to the local SQLite state
	sqlite string

	// The CRON schedule for Poasting
	schedule string
}

func configFromEnv() (config *Config, err error) {
	username, present := os.LookupEnv("MILKWEED_USERNAME")
	if !present {
		err = fmt.Errorf("Milkweed needs the MILKWEED_USERNAME to be set to a BlueSky username")
		return
	}
	password, present := os.LookupEnv("MILKWEED_PASSWORD")
	if !present {
		err = fmt.Errorf("Milkweed needs the MILKWEED_PASSWORD to be set to a BlueSky password")
		return
	}
	rss, present := os.LookupEnv("MILKWEED_RSS_URL")
	if !present {
		err = fmt.Errorf("Milkweed needs the MILKWEED_RSS_URL to be set to an RSS feed's URL")
		return
	}
	sqlite, present := os.LookupEnv("MILKWEED_SQLITE_PATH")
	if !present {
		sqlite = ":memory:"
	}
	schedule, present := os.LookupEnv("MILKWEED_CRON_SCHEDULE")
	if !present {
		err = fmt.Errorf("Milkweed needs the MILKWEED_CRON_SCHEDULE to be set to a valid refresh schedule in valid crontab syntax")
		return
	}
	config = &Config{
		username: username,
		password: password,
		rss:      rss,
		sqlite:   sqlite,
		schedule: schedule,
	}
	return
}

func main() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(slogger)
	logger := log.Default()

	config, err := configFromEnv()
	if err != nil {
		log.Fatal("failed to read configuration for Milkweed: '%v'", err)
	}

	ps, err := poasts.New(config.rss)
	if err != nil {
		logger.Println(err)
		logger.Fatal("failed to get post stream for Milkweed")
	}
	st, err := state.New(config.sqlite)
	if err != nil {
		logger.Println(err)
		log.Fatal("failed to initialize state for Milkweed")
	}
	for p, err := range ps.After(st.Published) {
		if err != nil {
			logger.Println(err)
			logger.Println("error reading poasts")
		}
		if err == nil {
			st.Publish(p.Published)
		}
		fmt.Println(p)
	}
}
