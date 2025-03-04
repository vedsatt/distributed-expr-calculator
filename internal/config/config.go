package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// конфигурация оркестратора
type Config struct {
	TimeAddition       time.Duration // сложение
	TimeSubtraction    time.Duration // вычитание
	TimeMultiplication time.Duration // умножение
	TimeDivision       time.Duration // деление

	ComputingPower int
}

// getEnv получает значение переменной среды
func getEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Environment variable %s is not set. Using default value: %d", key, defaultValue)
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil || result <= 0 {
		log.Printf("Invalid value for %s: %s. Using default value: %d", key, value, defaultValue)
		return defaultValue
	}

	return result
}

// загружает конфигурацию из переменных сред
func LoadConfig() Config {
	var cfg Config

	// получение времени выполнения операций с преобразованием в time.Duration
	cfg.TimeAddition = time.Duration(getEnv("TIME_ADDITION_MS", 2000)) * time.Millisecond
	cfg.TimeSubtraction = time.Duration(getEnv("TIME_SUBTRACTION_MS", 2000)) * time.Millisecond
	cfg.TimeMultiplication = time.Duration(getEnv("TIME_MULTIPLICATIONS_MS", 3000)) * time.Millisecond
	cfg.TimeDivision = time.Duration(getEnv("TIME_DIVISIONS_MS", 3000)) * time.Millisecond

	// получение вычислительной мощности
	cfg.ComputingPower = getEnv("COMPUTING_POWER", 3)

	return cfg
}
