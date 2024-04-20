package boot

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	g "main/app/global"
	"main/utils/file"
	"os"
	"time"
)

// Options 定义了日志选项
type Options struct {
	SavePath     string // 日志保存路径
	EncoderType  string // 编码器类型("json","console")
	EncodeLevel  string // 编码级别
	EncodeCaller string // 调用者类型
}

// 定义编码器类型和编码级别的常量
const (
	// JsonEncoder 编码器类型
	JsonEncoder   = "json"
	ConsoleEncode = "console"

	// LowercaseLevelEncoder 编码级别
	LowercaseLevelEncoder      = "LowercaseLevelEncoder"
	LowercaseColorLevelEncoder = "LowercaseColorLevelEncoder"
	CapitalLevelEncoder        = "CapitalLevelEncoder"
	CapitalColorLevelEncoder   = "CapitalColorLevelEncoder"

	// ShortCallerEncoder 调用者选项
	ShortCallerEncoder = "ShortCallerEncoder"
	FullCallerEncoder  = "FullCallerEncoder"
)

// LoggerSetup 函数设置日志
func LoggerSetup() {
	options := Options{
		SavePath:     g.Config.Logger.SavePath,
		EncoderType:  g.Config.Logger.EncoderType,
		EncodeLevel:  g.Config.Logger.EncodeLevel,
		EncodeCaller: g.Config.Logger.EncodeCaller,
	}
	LoggerSetupWithOptions(options)
}

// LoggerSetupWithOptions 函数使用指定的选项设置日志
func LoggerSetupWithOptions(options Options) {
	// 创建日志目录
	err := file.IsNotExistMkDir(options.SavePath)
	if err != nil {
		panic(err)
	}

	// 创建动态级别
	dynamicLevel := zap.NewAtomicLevel()
	// 创建各级别的优先级
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	// 获取编码器
	encoder := getEncoder(options)
	// 创建各级别的核心 //...表示让Go编译器根据初始值的数量来确定数组的长度
	cores := [...]zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, dynamicLevel), // 控制台输出
		// 根据级别分文件
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/all/server_all.log", options.SavePath)), zapcore.DebugLevel),
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/debug/server_debug.log", options.SavePath)), debugPriority),
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/info/server_info.log", options.SavePath)), infoPriority),
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/warn/server_warn.log", options.SavePath)), warnPriority),
		zapcore.NewCore(encoder, getWriteSyncer(fmt.Sprintf("./%s/error/server_error.log", options.SavePath)), errorPriority),
	}
	// 创建日志记录器 //...将cores切片中的元素展开并传递给zapcore.NewTee函数。zapcore.NewTee函数接受一个可变数量的zapcore.Core类型的参数
	zapLogger := zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
	defer func(zapLogger *zap.Logger) {
		_ = zapLogger.Sync()
	}(zapLogger)
	// 设置当前日志级别为"Debug"
	dynamicLevel.SetLevel(zap.DebugLevel)
	// 设置全局日志记录器
	g.Logger = zapLogger.Sugar()
	g.Logger.Info("initialize logger successfully!")
	//sugar.Debug("test")
	//sugar.Warn("test")
	//sugar.Error("test")
	//sugar.DPanic("test")
	//sugar.Panic("test")
	//sugar.Fatal("test")
}

// getEncoder函数根据选项返回一个编码器
func getEncoder(options Options) zapcore.Encoder {
	if options.EncoderType == JsonEncoder {
		// 如果编码器类型是JsonEncoder，返回一个JSON编码器
		return zapcore.NewJSONEncoder(getEncoderConfig(options))
	}
	// 否则，返回一个控制台编码器
	return zapcore.NewConsoleEncoder(getEncoderConfig(options))
}

// getEncoderConfig函数返回一个编码器配置
func getEncoderConfig(options Options) (config zapcore.EncoderConfig) {
	// 初始化编码器配置
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // format: short（package/filepath.go:line） full (filepath.go:line)
	}
	// 根据选项设置编码级别
	switch {
	case options.EncodeLevel == LowercaseLevelEncoder: // default
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case options.EncodeLevel == LowercaseColorLevelEncoder:
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case options.EncodeLevel == CapitalLevelEncoder:
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case options.EncodeLevel == CapitalColorLevelEncoder:
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	// 如果调用者类型是ShortCallerEncoder，设置编码调用者为短格式
	if options.EncodeCaller == ShortCallerEncoder {
		config.EncodeCaller = zapcore.ShortCallerEncoder
	}
	return config
}

// getWriteSyncer函数返回一个写同步器
func getWriteSyncer(file string) zapcore.WriteSyncer {
	// 创建一个lumberjack日志记录器
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, // 日志文件位置
		MaxSize:    1,    // 日志文件最大大小（MB）
		MaxBackups: 100,
		MaxAge:     30, // 天
		Compress:   true,
	}
	// 返回一个写同步器
	return zapcore.AddSync(lumberJackLogger)
}

// CustomTimeEncoder 函数是一个自定义的时间编码器
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 将时间格式化为"[2006-01-02 15:04:05.000]"格式，并添加到编码器中
	enc.AppendString(t.Format("[2006-01-02 15:04:05.000]"))
}
