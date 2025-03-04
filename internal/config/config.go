package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config представляет конфигурацию
type Config struct {
	TimeAddition       time.Duration // время для сложения
	TimeSubtraction    time.Duration // время для вычитания
	TimeMultiplication time.Duration // время для умножения
	TimeDivision       time.Duration // время для деления
	ComputingPower     int           // вычислительная мощность
}

// LoadConfig загружает конфигурацию из файла .env или использует значения по умолчанию
func LoadConfig() Config {
	// значения по умолчанию
	cfg := Config{
		TimeAddition:       2000 * time.Millisecond,
		TimeSubtraction:    2000 * time.Millisecond,
		TimeMultiplication: 3000 * time.Millisecond,
		TimeDivision:       3000 * time.Millisecond,
		ComputingPower:     3,
	}

	// Открываем файл .env
	file, err := os.Open(".env")
	if err != nil {
		log.Println("File .env not found. Using default values.")
		return cfg
	}
	defer file.Close()

	// Читаем файл построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Пропускаем пустые строки и комментарии
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		// Разделяем строку на ключ и значение
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Присваиваем значения в зависимости от ключа
		switch key {
		case "TIME_ADDITION_MS":
			if v, err := strconv.Atoi(value); err == nil && v > 0 {
				cfg.TimeAddition = time.Duration(v) * time.Millisecond
			}
		case "TIME_SUBTRACTION_MS":
			if v, err := strconv.Atoi(value); err == nil && v > 0 {
				cfg.TimeSubtraction = time.Duration(v) * time.Millisecond
			}
		case "TIME_MULTIPLICATIONS_MS":
			if v, err := strconv.Atoi(value); err == nil && v > 0 {
				cfg.TimeMultiplication = time.Duration(v) * time.Millisecond
			}
		case "TIME_DIVISIONS_MS":
			if v, err := strconv.Atoi(value); err == nil && v > 0 {
				cfg.TimeDivision = time.Duration(v) * time.Millisecond
			}
		case "COMPUTING_POWER":
			if v, err := strconv.Atoi(value); err == nil && v > 0 {
				cfg.ComputingPower = v
			}
		}
	}

	// Проверяем ошибки сканера
	if err := scanner.Err(); err != nil {
		log.Println("Error reading .env file:", err)
	}

	return cfg
}
