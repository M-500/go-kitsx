package zapx

type LoggerX interface {
	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
}
