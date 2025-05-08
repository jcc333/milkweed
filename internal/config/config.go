package config

import (
	"fmt"
	"os"
)

// The configuration for milkweed consists of:
type Config struct {
	// The BlueSky Username
	Username string

	// The BlueSky App Password
	Password string

	// The PDS Server
	Server string

	// The URL of the RSS feed endpoint
	RSS string

	// The path to the local SQLite state
	SQLite string

	// The CRON Schedule for Poasting
	Schedule string
}

func FromEnv() (config *Config, err error) {
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
	server, present := os.LookupEnv("MILKWEED_PDS")
	if !present {
		server = "https://bsky.social"
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
		Username: username,
		Password: password,
		Server:   server,
		RSS:      rss,
		SQLite:   sqlite,
		Schedule: schedule,
	}
	return
}
