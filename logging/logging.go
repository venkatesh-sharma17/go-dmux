package logging

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

//LogConf holds logging configuration
type LogConf struct {
	Path        string     `json:"path"`
	EnableDebug bool       `json:"enable_debug"`
	Rotate      RotateConf `json:"rotation"`
}

//RotateConf holds logging Rotation configuration
type RotateConf struct {
	FileSize     int  `json:"size_in_mb"`
	RotationDays int  `json:"retention_days"`
	NoOfFiles    int  `json:"retention_count"`
	Compress     bool `json:"compress"`
}

//DMuxLogging logging
type DMuxLogging struct{}

//Start starting logging
func (c *DMuxLogging) Start(logconf LogConf) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   logconf.Path,
		MaxSize:    logconf.Rotate.FileSize, // megabytes
		MaxBackups: logconf.Rotate.NoOfFiles,
		MaxAge:     logconf.Rotate.RotationDays, //days
		Compress:   logconf.Rotate.Compress,     // disabled by default
	})
}
