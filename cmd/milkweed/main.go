package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jcc333/milkweed/internal/config"
	"github.com/jcc333/milkweed/internal/poasts"
	"github.com/jcc333/milkweed/internal/state"
	"github.com/jcc333/milkweed/send"
	_ "github.com/joho/godotenv/autoload"
	// "github.com/robfig/cron"
)

// Finds new skeets and publishes them
func publishNewSkeets(ctx context.Context, ps *poasts.Poasts, st *state.State, sender send.Sender) {
	logger := log.Default()
	for p, err := range ps.All() {
		if err != nil {
			logger.Println(err)
			logger.Println("error reading poasts")
		}
		isPublished, err := st.IsPublished(p.GUID)
		if err != nil {
			logger.Println(err)
			logger.Println("error checking if poast is published")
			continue
		}
		modifier := ""
		if !isPublished {
			modifier = "not"
		}
		logger.Println(p.Title, " is ", modifier, " published")
		if !isPublished {
			logger.Println("Publishing ", p.GUID, " from ", p.Published)
			st.Publish(p.GUID, p.Published)
			err = sender.Send(ctx, p, st)
			if err != nil {
				logger.Println(err)
				logger.Println("error sending skeet")
			}
		}
	}
	return
}

func main() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(slogger)
	logger := log.Default()

	cfg, err := config.FromEnv()
	if err != nil {
		logger.Fatal("failed to read configuration for Milkweed: '%v'", err)
	}

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	ps, err := poasts.New(cfg.RSS, logger)
	if err != nil {
		logger.Println(err)
		logger.Fatal("failed to get post stream for Milkweed")
	}
	st, err := state.New(cfg.SQLite, logger, ctx)
	defer st.Close()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		stop()
	}()

	if err != nil {
		logger.Println(err)
		log.Fatal("failed to initialize state for Milkweed")
	}
	snd := send.New(cfg)
	err = snd.Connect(ctx)
	if err != nil {
		logger.Println(err)
		log.Fatal("failed to connect sender for Milkweed")
	}
	defer func() {
		err := snd.Close()
		if err != nil {
			logger.Println(err)
		}
	}()
	publishNewSkeets(ctx, ps, st, snd)
	return
	// c := cron.New()
	// err = c.AddFunc(cfg.schedule, func() { publishNewSkeets(ps, st) })
	// if err != nil {
	// 	logger.Fatal("failed to add CRON func on schedule '", cfg.Schedule, "'")
	// }
	// c.Run()
}
