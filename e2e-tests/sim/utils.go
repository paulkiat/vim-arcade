package sim

import (
	"context"
	"log/slog"
	"time"

	assert "vim-arcade.theprimeagen.com/pkg/assert"
	prettylog "vim-arcade.theprimeagen.com/pkg/pretty-log"
)

func KillContext(cancel context.CancelFunc) {
    go func() {
        time.Sleep(time.Second * 5)
        cancel()
        assert.Never("context should never be killed with KillContext")
    }()
}

func CreateLogger(name string) *slog.Logger {
    logger := prettylog.CreateLoggerFromEnv(nil)
    logger = logger.With("area", name).With("process", "sim")
    slog.SetDefault(logger)

    logger.Error("Test Logger Created")

    return logger
}


