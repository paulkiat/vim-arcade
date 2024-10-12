// traceroute/logger.go

package traceroute

import (
    "log"

    "enterprise_networking_suite/logger"
)

func init() {
    // Redirect traceroute module logs to the centralized logger
    log.SetOutput(logger.Logger.Writer())
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}