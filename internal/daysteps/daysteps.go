package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}

	stepsStr := parts[0]
	durStr := parts[1]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
}
