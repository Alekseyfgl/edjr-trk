package service

import (
	"edjr-trk/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sync"
	"time"
)

type RateLimiter struct {
	requests  map[string][]time.Time // Карта для хранения временных меток запросов по IP
	blocked   map[string]time.Time   // Карта для хранения времени блокировки IP
	mu        sync.Mutex             // Мьютекс для конкурентного доступа
	limit     int                    // Количество запросов
	window    time.Duration          // Окно времени
	blockTime time.Duration          // Время блокировки
}

// Создание нового лимитера
func NewRateLimiter(limit int, window, blockTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:  make(map[string][]time.Time),
		blocked:   make(map[string]time.Time),
		limit:     limit,
		window:    window,
		blockTime: blockTime,
	}

	// Запуск очистки карты каждые 1 дней
	go rl.cleanup(1 * 24 * time.Hour)

	return rl
}

// Проверка запроса
func (rl *RateLimiter) ValidateRequest(c *fiber.Ctx) error {
	clientIP := utils.GetClientIP(c)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Проверка на блокировку
	if unblockTime, ok := rl.blocked[clientIP]; ok {
		if time.Now().Before(unblockTime) {
			return fmt.Errorf("IP %s заблокирован до %v", clientIP, unblockTime)
		}
		delete(rl.blocked, clientIP) // Убираем блокировку, если срок истёк
	}

	// Удаляем старые записи из окна
	now := time.Now()
	requestTimes := rl.requests[clientIP]
	var filtered []time.Time
	for _, t := range requestTimes {
		if now.Sub(t) <= rl.window {
			filtered = append(filtered, t)
		}
	}
	rl.requests[clientIP] = filtered

	// Проверка количества запросов
	if len(filtered) >= rl.limit {
		rl.blocked[clientIP] = now.Add(rl.blockTime)
		delete(rl.requests, clientIP) // Убираем историю запросов для заблокированного IP
		return fmt.Errorf("IP %s заблокирован на %v", clientIP, rl.blockTime)
	}

	// Добавляем текущий запрос
	rl.requests[clientIP] = append(rl.requests[clientIP], now)
	return nil
}

// Очистка карты
func (rl *RateLimiter) cleanup(interval time.Duration) {
	for {
		time.Sleep(interval)

		rl.mu.Lock()
		now := time.Now()

		// Удаляем данные о запросах и блокировках
		for ip, t := range rl.blocked {
			if now.After(t) {
				delete(rl.blocked, ip)
			}
		}

		for ip, times := range rl.requests {
			var filtered []time.Time
			for _, t := range times {
				if now.Sub(t) <= rl.window {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = filtered
			}
		}
		rl.mu.Unlock()
	}
}
