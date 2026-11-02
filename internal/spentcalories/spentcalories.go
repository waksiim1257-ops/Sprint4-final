package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных тренировки")
	}

	stepsStr := parts[0]
	activity := parts[1]
	durStr := parts[2]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	stepLen := height * stepLengthCoefficient
	distMeters := float64(steps) * stepLen
	return distMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	if hours <= 0 {
		return 0
	}
	return dist / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры для расчета калорий при беге")
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("некорректные параметры для расчета калорий при беге")
	}

	durationMinutes := duration.Minutes()
	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры для расчета калорий при ходьбе")
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, errors.New("некорректные параметры для расчета калорий при ходьбе")
	}

	durationMinutes := duration.Minutes()
	calories := (weight * speed * durationMinutes) / minInH
	calories = calories * walkingCaloriesCoefficient
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	switch activity {
	case "Бег":
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}

		result := fmt.Sprintf("Тип тренировки: Бег\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", duration.Hours(), dist, speed, calories)
		return result, nil
	case "Ходьба":
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}

		result := fmt.Sprintf("Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", duration.Hours(), dist, speed, calories)
		return result, nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}
