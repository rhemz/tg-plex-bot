package main

import (
    "go.uber.org/zap"
    "github.com/rhemz/tg-plex-bot/server"
)


func main() {
    // replace the sytem logger
    logger, _ := zap.NewProduction()
    undoLoggerReplace := zap.ReplaceGlobals(logger)
    defer undoLoggerReplace()

    // start doing server-y things
    server.Start()
}
