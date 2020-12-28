package log

import (
	"errors"
	"fmt"
	"github.com/zzlpeter/aps-go/libs/tomlc"
	"github.com/zzlpeter/aps-go/libs/utils"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var loggerN = make(map[string]*zap.Logger)
var once sync.Once
const (
	INFO = "info"
	DEBUG = "debug"
	WARNING = "warning"
	ERROR = "error"
	FATAL = "fatal"
)

func getLogger() {
	once.Do(func() {
		makeLogger()
	})
}

func makeLogger() {
	logConfS := tomlc.Config{}.LogConfS()
	for alias, conf := range logConfS {
		fmt.Println(alias, conf)
		// 获取dir路径
		dir := filepath.Dir(conf["file"].(string))
		// 判断dir路径是否存在
		if _, err := os.Stat(dir); err != nil {
			log.Fatalf("LogPath: %v is not existed", dir)
		}

		hook := lumberjack.Logger{
			Filename:   conf["file"].(string), 		// 日志文件路径
			MaxSize:    conf["max_size"].(int),     // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: conf["max_backups"].(int),  // 日志文件最多保存多少个备份
			MaxAge:     conf["max_age"].(int),      // 文件最多保存多少天
			Compress:   false,                      // 是否压缩
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
			EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder, //
			EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
			EncodeName:     zapcore.FullNameEncoder,
		}

		// 设置日志级别
		atomicLevel := zap.NewAtomicLevel()
		atomicLevel.SetLevel(zap.DebugLevel)

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),                // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)),  // 打印到文件
			atomicLevel,                                          // 日志级别
		)

		// 开启开发模式，堆栈跟踪
		caller := zap.AddCaller()
		// 开启文件及行号
		development := zap.Development()
		// 设置初始化字段
		hostName, _ := utils.HostManager{}.LocalHostName()
		hostIp, _ := utils.HostManager{}.LocalIp()
		filed := zap.Fields(zap.String("host_name", hostName), zap.String("host_ip", hostIp))
		// 构造日志
		loggerN[alias] = zap.New(core, caller, development, filed)
	}
}

func init() {
	getLogger()
}

// 记录日志
// 具体使用方法参考 LogRecord("default", "trace-id", log.INFO, "测试信息", "name", "alex", "age", 24)
// target代表写入哪个日志文件, 该文件alias必须在conf[log]模块中定义，否则默认写入default
// tid为trace-id，方法入口处会自动生成
func LogRecord(target, tid, level, msg string, kvs ...interface{}) error {
	if len(kvs) % 2 == 1 {
		return errors.New("kvs should be coupled")
	}
	var logger *zap.Logger
	if _logger, ok := loggerN[target]; ok {
		logger = _logger
	} else {
		logger = loggerN["default"]
	}
	err, zapFields := makeZapFields(tid, kvs...)
	if err != nil {
		return err
	}
	switch level {
	case DEBUG:
		logger.Debug(msg, zapFields...)
	case INFO:
		logger.Info(msg, zapFields...)
	case WARNING:
		logger.Warn(msg, zapFields...)
	case ERROR:
		logger.Error(msg, zapFields...)
	case FATAL:
		logger.Fatal(msg, zapFields...)
	default:
		logger.Info(msg, zapFields...)
	}
	return nil
}

// 组装zap需要的日志参数
func makeZapFields(tid string, kvs ...interface{}) (error, []zap.Field) {
	var arr []zap.Field
	// 校验kvs是否合法
	l := len(kvs)
	var idx = 0
	for {
		if idx >= l {
			break
		}
		k, ok := kvs[idx].(string)
		if !ok {
			return errors.New(fmt.Sprintf("Type of k: %v is not string", kvs[idx])), arr
		}
		arr = append(arr, zap.Any(k, kvs[idx+1]))
		idx+=2
	}
	// 设置trace-ID
	if tid != "" {
		arr = append(arr, zap.String("trace_id", tid))
	}
	// 获取调用栈信息
	_, file, line, _ := runtime.Caller(2)
	arr = append(arr, zap.String("file", file), zap.Int("line", line))

	return nil, arr
}