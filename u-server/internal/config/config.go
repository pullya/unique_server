package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	configFile = "config.yaml"

	logLevelDebug = "Debug"
	logLevelInfo  = "Info"
	logLevelWarn  = "Warn"
	logLevelError = "Error"
	logLevelFatal = "Fatal"
	logLevelPanic = "Panic"
	logLevelTrace = "Trace"

	envPort        = "U_SERVER_PORT"
	envTimeout     = "U_SERVER_TIMEOUT"
	envServiceName = "U_SERVER_SERVICE_NAME"
	envEndpoint    = "U_SERVER_ENDPOINT"
	envStrictMode  = "U_SERVER_STRICT_MODE"
	envMaxIpConn   = "U_SERVER_MAX_IP_CONN"
	envLogLevel    = "U_SERVER_LOG_LEVEL"

	shardsCount = 8
)

var Config Configuration

var logLevelsMap = map[LogLevel]log.Level{
	logLevelDebug: log.DebugLevel,
	logLevelInfo:  log.InfoLevel,
	logLevelWarn:  log.WarnLevel,
	logLevelError: log.ErrorLevel,
	logLevelFatal: log.FatalLevel,
	logLevelPanic: log.PanicLevel,
	logLevelTrace: log.TraceLevel,
}

var envArray = []string{
	envPort,
	envTimeout,
	envServiceName,
	envEndpoint,
	envStrictMode,
	envMaxIpConn,
	envLogLevel,
}

type LogLevel string

type Configuration struct {
	Port        int    `yaml:"port"`
	Timeout     int    `yaml:"timeout"`
	ServiceName string `yaml:"serviceName"`

	Endpoint   string `yaml:"endpoint"`
	StrictMode string `yaml:"strictMode"`
	MaxIpConn  int    `yaml:"maxIpConnection"`

	ShardsCnt int `yaml:"shardsCnt"`

	LogLevel LogLevel `yaml:"logLevel"`
}

func ReadConfig() {
	log.SetLevel(log.DebugLevel)

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read configuration from file '%s', error: %v", configFile, err)
		return
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("failed to unmarshall config file '%s', error: %v", configFile, err)
		return
	}
	Config.ShardsCnt = shardsCount

	log.Debugf("Default configuration read: %v", Config)

	checkEnv()
}

func (l LogLevel) ToLogrusFormat() log.Level {
	res, ok := logLevelsMap[l]
	if !ok {
		return log.ErrorLevel
	}
	return res
}

func checkEnv() {
	for _, name := range envArray {
		if envVal := os.Getenv(name); envVal != "" {
			switch name {
			case envPort:
				port, err := validatePort(envVal)
				if err == nil {
					Config.Port = port
					log.Debugf("tcpPort set to %d", Config.Port)
				}
			case envTimeout:
				timeout, err := validateTimeout(envVal)
				if err == nil {
					Config.Timeout = timeout
					log.Debugf("timeout set to %d", Config.Timeout)
				}
			case envServiceName:
				Config.ServiceName = envVal
				log.Debugf("serviceName set to '%s'", Config.ServiceName)
			case envEndpoint:
				Config.Endpoint = envVal
				log.Debugf("endpoint set to %s", Config.Endpoint)
			case envStrictMode:
				sm, err := validateMode(envVal)
				if err == nil {
					Config.StrictMode = sm
					log.Debugf("strictMode set to %s", Config.StrictMode)
				}
			case envMaxIpConn:
				mic, err := validateMaxIpConn(envVal)
				if err == nil {
					Config.MaxIpConn = mic
					log.Debugf("maxIpConn set to %d", Config.MaxIpConn)
				}
			case envLogLevel:
				ll, err := validateLogLevel(envVal)
				if err == nil {
					Config.LogLevel = ll
					log.Debugf("logLevel set to '%v'", Config.LogLevel)
				}
			}
		}
	}
}

func validatePort(in string) (int, error) {
	num, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if num < 0 || num > 65535 {
		return 0, errors.New("incorrect port number")
	}
	return num, nil
}

func validateTimeout(in string) (int, error) {
	num, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if num < 0 {
		return 0, errors.New("incorrect timeout")
	}
	return num, nil
}

func validateMode(in string) (string, error) {
	if in != "on" && in != "off" {
		return "", errors.New("incorrect mode")
	}
	return in, nil
}

func validateMaxIpConn(in string) (int, error) {
	num, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if num < 1 {
		return 0, errors.New("incorrect max ip connection")
	}
	return num, nil
}

func validateLogLevel(in string) (LogLevel, error) {
	if _, ok := logLevelsMap[LogLevel(in)]; !ok {
		return "", errors.New("incorrect log level")
	}
	return LogLevel(in), nil
}

func BuildPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
