package log

import (
	"log"
	"log/syslog"
	"os"
)

var errLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)

// InitLogger start a syslog connection, and set logging output to syslog.
func InitLogger(foreground bool) {
	log.Println("init logger")

	if foreground {
		initForeground()
	} else {
		initBackground()
	}
}

func initForeground() {
	errLogger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func initBackground() {
	errorWriter, err := syslog.New(syslog.LOG_ERR|syslog.LOG_DAEMON, "legsc")
	if err != nil {
		errLogger = log.New(os.Stdout, "legsc:", log.Ldate|log.Ltime|log.Lmicroseconds)
		errLogger.Println("failed in dial to syslog daemon")
		errLogger.Println(err)
	} else {
		errLogger = log.New(errorWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	infoWriter, err := syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "legsc")
	if err != nil {
		errLogger.Println("failed in dial to syslog daemon")
		errLogger.Println(err)
	} else {
		// overwrite default logger
		log.SetOutput(infoWriter)
		log.SetPrefix("[INFO] ")
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	}
}

// Info outputs log at info level.
// This method same with log.Println(v).
func Info(v ...interface{}) {
	log.Println(v...)
}

// Error outputs log at error level.
func Error(v ...interface{}) {
	errLogger.Println(v...)
}
