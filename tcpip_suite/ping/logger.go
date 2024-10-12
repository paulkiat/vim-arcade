// ping/logger.go

package ping

import (
    "log"
    "os"

    "enterprise_networking_suite/logger"
)

func init() {
    // Redirect ping module logs to the centralized logger
    log.SetOutput(logger.Logger.Writer())
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}