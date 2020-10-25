package logger

import (
	"fmt"
	"sync"
)

type Logger struct {
	mu      sync.Mutex            //锁
	writers map[string]LogWriterI //打印机
}

func NewLogger(confs ...ConfigI) *Logger {
	back := &Logger{
		writers: make(map[string]LogWriterI),
	}
	return back
}

//刷入
func (log *Logger) Flush() {
	for _, w := range log.writers {
		w.Flush()
	}
}

//设置打印器
func (log *Logger) SetWriter(writer LogWriterI) {
	if writer == nil {
		return
	}
	log.writers[writer.Name()] = writer
}

//获得打印器
func (log *Logger) GetWriter(writerName string) LogWriterI {
	if back, ok := log.writers[writerName]; ok {
		return back
	}
	return nil
}

func (log *Logger) doWrite(lvl int, format string, args ...interface{}) {
	log.mu.Lock()
	defer log.mu.Unlock()
	content := fmt.Sprintf(format, args...)
	content += "\n"
	for _, writer := range log.writers {
		writer.Write(lvl, content)
	}
}

func (log *Logger) write(lvl int, format string, args ...interface{}) {
	log.doWrite(lvl, format, args...)
}

func (log *Logger) Fatal(format string, args ...interface{}) {
	log.doWrite(FATAL, format, args...)
	log.Flush()
}

func (log *Logger) Error(format string, args ...interface{}) {
	log.doWrite(ERROR, format, args...)
}

func (log *Logger) Warning(format string, args ...interface{}) {
	log.doWrite(WARNING, format, args...)
}

func (log *Logger) Notice(format string, args ...interface{}) {
	log.doWrite(NOTICE, format, args...)
}

func (log *Logger) Debug(format string, args ...interface{}) {
	log.doWrite(DEBUG, format, args...)
}

func (log *Logger) Info(format string, args ...interface{}) {
	log.doWrite(INFO, format, args...)
}
