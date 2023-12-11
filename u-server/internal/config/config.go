package config

import log "github.com/sirupsen/logrus"

const (
	ServiceName     = "u-server" // Имя сервиса для отображения в логах
	Endpoint        = "/tws"     // Имя эндпоинта
	WsPort          = 8081       // Номер порта для соединения с tcp-сервером
	MaxIpConnection = 1          // Максимальное количество соединений с одного IP адреса

	LogLevel = log.DebugLevel // Уровень логирования
)
