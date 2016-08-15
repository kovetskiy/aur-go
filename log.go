package aur

func debugf(format string, values ...interface{}) {
	logger.Debugf(format, values...)
}

func debugln(value interface{}) {
	logger.Debug(value)
}
