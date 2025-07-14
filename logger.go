package main

import (
	"context"
	"gRPC-Learner-Project/proto/logger"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() error {

	logPath := "logs/learner.log"
	logDir := filepath.Dir(logPath)

	// Ensure the directory exists
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return err
	}

	// Ensure the file exists (no-op if already exists)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		file, err := os.Create(logPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	cfg := zap.NewProductionConfig()

	// Set the minimum log level (Debug, Info, Warn, Error)
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel) // Using InfoLevel for production

	// Set JSON encoding and log output paths
	cfg.Encoding = "json"

	cfg.OutputPaths = []string{
		"stdout", // Logs to console
		logPath,  // Logs to file
	}

	// Customize encoder to include caller information
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // Full caller path for detailed location
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

	// Keep production settings but explicitly add caller information
	cfg.Development = false      // Production mode
	cfg.DisableStacktrace = true // Only include stacktraces for panics in production

	// Build the logger with production settings but explicitly including caller info
	log, err := cfg.Build(zap.AddCaller()) // Explicitly adding caller even in production mode
	if err != nil {
		return err
	}

	Logger = log
	return nil
}

func (s *Server) Log(_ context.Context, req *logger.LogRequest) (*logger.LogResponse, error) {

	if Logger == nil {
		return &logger.LogResponse{
			Success: false,
			Err:     "Logger not initialized",
		}, nil
	}

	switch req.Level {
	case "INFO":
		Logger.WithOptions(zap.AddCallerSkip(1)).Info(req.Message, zap.String("service", req.Service))
	case "DEBUG":
		Logger.WithOptions(zap.AddCallerSkip(1)).Debug(req.Message, zap.String("service", req.Service))
	case "ERROR":
		Logger.WithOptions(zap.AddCallerSkip(1)).Error(req.Message, zap.String("service", req.Service))
	}

	Logger.Sync()

	return &logger.LogResponse{
		Success: true,
	}, nil
}
