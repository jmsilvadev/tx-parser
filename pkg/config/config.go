package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jmsilvadev/tx-parser/pkg/jsonrpc"
	"github.com/jmsilvadev/tx-parser/pkg/logger"
	"github.com/jmsilvadev/tx-parser/pkg/parser"
	"github.com/jmsilvadev/tx-parser/pkg/parser/leveldb"
	"github.com/jmsilvadev/tx-parser/pkg/parser/memorydb"
)

// Default Values
var (
	parserEngine   = "leveldb"
	dbPath         = "/tmp/parser.db"
	serverPort     = ":5000"
	loggerLevel    = "DEBUG"
	environment    = "dev"
	timeout        = "1s"
	defaultTimeout = time.Second
	cliUrl         = "https://ethereum-rpc.publicnode.com"
)

type Config struct {
	ServerPort string
	Env        string
	Timeout    time.Duration
	Parser     parser.Parser
	Logger     logger.Logger
}

func New(ctx context.Context, port, env string, duration time.Duration, p parser.Parser, logger logger.Logger) *Config {
	return &Config{
		ServerPort: port,
		Env:        env,
		Timeout:    duration,
		Logger:     logger,
		Parser:     p,
	}
}

func GetDefaultConfig() *Config {
	environment = getEnv("ENV", environment)
	serverPort = getEnv("SERVER_PORT", serverPort)
	loggerLevel = getEnv("LOG_LEVEL", loggerLevel)
	dbPath = getEnv("DB_PATH", dbPath)
	parserEngine = getEnv("PARSER_ENGINE", parserEngine)
	cliUrl = getEnv("JSONRPC_URL", cliUrl)

	timeout = getEnv("TIMEOUT", timeout)
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		duration = defaultTimeout
	}

	level := logger.LEVEL_ERROR
	if loggerLevel == "INFO" {
		level = logger.LEVEL_INFO
	}
	if loggerLevel == "WARN" {
		level = logger.LEVEL_WARN
	}
	if loggerLevel == "DEBUG" {
		level = logger.LEVEL_DEBUG
	}
	log := logger.New(level)

	db, err := getDatabase(parserEngine, dbPath, cliUrl, log)
	if err != nil {
		log.Info("invalid database")
		panic("invalid database")
	}

	ctx := context.Background()
	config := New(ctx, serverPort, environment, duration, db, log)

	return config
}

func getDatabase(parserEngine, dbPath, cliUrl string, l logger.Logger) (parser.Parser, error) {
	var (
		p   parser.Parser
		err error
	)

	cli := jsonrpc.NewEthereum(l, cliUrl)
	p = memorydb.New(cli, l)
	if strings.ToLower(parserEngine) == "leveldb" {
		p, err = leveldb.New(dbPath, cli, l)
	}

	if err != nil {
		return nil, fmt.Errorf("DB ERROR: %s", err.Error())
	}

	return p, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
