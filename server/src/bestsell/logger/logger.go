package logger

import (
	"fmt"
	"bestsell/common"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	DbgLogger *log.Logger
	InfLogger *log.Logger
	ErrLogger *log.Logger
	infOutput io.Writer
	errOutput io.Writer
	_logLevel byte
	_logPath  string
)

var (
	errLoggerAppendMutex sync.Mutex
	errLoggerReferences  []**log.Logger
)

func Init(logPath string, logLevel byte) {
	fmt.Println("init logger")
	_logPath = logPath
	_logLevel = logLevel
	if logPath != "" && logPath != "/" {
		_ = os.MkdirAll(_logPath, 0666)
		setLoggers()
	} else {
		if _logLevel >= 3 {
			DbgLogger = log.New(os.Stdout, "", 0)
		}
		if _logLevel >= 2 {
			InfLogger = log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
			infOutput = os.Stdout
		}
		ErrLogger = log.New(os.Stderr, "", log.Ltime|log.Lshortfile)
		errOutput = os.Stderr
	}
}

func setLoggers() {
	nowTime := common.Now()
	oldFile := _logPath + "/" + nowTime.AddDate(0, 0, -30).Format("2006_01_02")
	newFile := _logPath + "/" + nowTime.Format("2006_01_02")
	if _logLevel >= 3 {
		_ = os.Remove(oldFile + ".dbg")
		fileDbg, _ := os.OpenFile(newFile+".dbg", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		DbgLogger = log.New(fileDbg, "", log.Ltime|log.Lshortfile)
	} else {
		DbgLogger = nil
	}
	if _logLevel >= 2 {
		_ = os.Remove(oldFile + ".inf")
		fileInf, _ := os.OpenFile(newFile+".inf", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		InfLogger = log.New(fileInf, "", log.Ltime|log.Lshortfile)
		infOutput = fileInf
	} else {
		InfLogger = nil
		infOutput = nil
	}
	_ = os.Remove(oldFile + ".err")
	fileErr, _ := os.OpenFile(newFile+".err", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	ErrLogger = log.New(fileErr, "", log.Ltime|log.Lshortfile)
	errOutput = fileErr
	os.Stderr = fileErr
	for _, r := range errLoggerReferences {
		*r = ErrLogger
	}
	hour, minute, second := nowTime.Clock()
	second = hour*3600 + minute*60 + second
	if second >= 86400 {
		second = 0
	}
	time.AfterFunc(time.Duration(86400-second)*time.Second, func() {
		setLoggers()
	})
}

func SetLogLevel(logLevel byte) {
	_logLevel = logLevel
	if _logPath != "" && _logPath != "/" {
		now := common.Now()
		newFile := _logPath + "/" + now.Format("2006_01_02")
		if _logLevel >= 3 {
			if DbgLogger == nil {
				fileDbg, _ := os.OpenFile(newFile+".dbg", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
				DbgLogger = log.New(fileDbg, "", log.Ltime|log.Lshortfile)
			}
		} else {
			DbgLogger = nil
		}
		if _logLevel >= 2 {
			if InfLogger == nil {
				fileInf, _ := os.OpenFile(newFile+".inf", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
				InfLogger = log.New(fileInf, "", log.Ltime|log.Lshortfile)
				infOutput = fileInf
			}
		} else {
			InfLogger = nil
			infOutput = nil
		}
	} else {
		if _logLevel >= 3 {
			DbgLogger = log.New(os.Stdout, "", 0)
		}
		if _logLevel >= 2 {
			InfLogger = log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
			infOutput = os.Stdout
		}
	}

}

func Debug(str string) {
	if _logLevel >= 3 && DbgLogger != nil {
		_ = DbgLogger.Output(2, str)
	}
}

func Debugf(format string, v ...interface{}) {
	if _logLevel >= 3 && DbgLogger != nil {
		_ = DbgLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Info(str string) {
	if _logLevel >= 2 && InfLogger != nil {
		_ = InfLogger.Output(2, str)
	}
}

func Infof(format string, v ...interface{}) {
	if _logLevel >= 2 && InfLogger != nil {
		_ = InfLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func PrintInfo(buffs ...[]byte) {
	if infOutput != nil {
		var l int
		for _, buff := range buffs {
			l += len(buff)
		}
		var b = make([]byte, 0, l+10)
		formatHeader(&b, common.Now())
		for _, buff := range buffs {
			b = append(b, buff...)
		}
		if b[len(b)-1] != '\n' {
			b = append(b, '\n')
		}
		_, _ = infOutput.Write(b)
	}
}

func Error(str string) {
	_ = ErrLogger.Output(2, str)
}

func Errorf(format string, v ...interface{}) {
	_ = ErrLogger.Output(2, fmt.Sprintf(format, v...))
}

func PrintStack(v interface{}) {
	var s string
	if v != nil {
		s = fmt.Sprintf("%v\n", v)
	}
	l := len(s)
	b := make([]byte, l+1034)
	b2 := b[:0]
	formatHeader(&b2, common.Now())
	copy(b[9:], s)
	n := runtime.Stack(b[9+l:], false)
	if 10+l+n <= len(b) {
		b = b[:10+l+n]
	}
	b[len(b)-1] = '\n'
	_, _ = errOutput.Write(b)
}

func AppendErrLogger(l **log.Logger) {
	errLoggerAppendMutex.Lock()
	errLoggerReferences = append(errLoggerReferences, l)
	*l = ErrLogger
	errLoggerAppendMutex.Unlock()
}

func formatHeader(buf *[]byte, t time.Time) {
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	*buf = append(*buf, ' ')
}

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
