package rlog

import (
	"sync"
	"io"
	"github.com/fatih/color"
	"time"
	"fmt"
	"os"
	"runtime"
)

const (
	LevelError = iota
	LevelDebug
	LevelWarning
	LevelInfo
)

const (
	textDebug = "DEBUG"
	textInfo = "INFO"
	textWarning = "WARNING"
	textError = "ERROR"
)

type logger struct{
	mu sync.Mutex
	stdout io.Writer
	useColor bool
	logLevel int
	showLine bool
}

var l *logger

func init(){
	l = &logger{
		stdout:os.Stdout,
	}
}

func (log *logger) write(level int, message interface{}, name ...string){
	if level > l.logLevel{
		//return
	}
	if log.stdout != nil{
		var label, output string

		switch level {
		case LevelDebug:
			label = textDebug
		case LevelInfo:
			label = textInfo
		case LevelWarning:
			label = textWarning
		case LevelError:
			label = textError
		}

		if l.showLine{
			_, fn, line, _ := runtime.Caller(2)
			output = fmt.Sprintf("%s %s:%d  [%v]: %v", time.Now().Format("2006-01-02 15:04:05"), fn, line, label, message)
		}else{
			output = fmt.Sprintf("%s [%v]: %v", time.Now().Format("2006-01-02 15:04:05"), label, message)
		}

		if log.useColor{
			color.Output = log.stdout
			switch level {
			case LevelDebug:
				color.Green(output)
			case LevelInfo:
				color.White(output)
			case LevelWarning:
				color.Yellow(output)
			case LevelError:
				color.Red(output)
			}
		}else{
			log.mu.Lock()
			log.stdout.Write([]byte(output + "\n"))
			log.mu.Unlock()
		}

	}
}

func Info(message interface{}){
	l.write(LevelInfo, message)
}

func Debug(message interface{}){
	l.write(LevelDebug, message)
}

func Warn(message interface{}){
	l.write(LevelWarning, message)
}

func Err(message interface{}){
	l.write(LevelError, message)
}

func UseColor(val bool){
	l.mu.Lock()
	l.useColor = val
	l.mu.Unlock()
}

func SetOutput(w io.Writer){
	l.mu.Lock()
	l.stdout = w
	l.mu.Unlock()
}

func ShowRuntimeInfo( val bool){
	l.mu.Lock()
	l.showLine = val
	l.mu.Unlock()
}
