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
	if data == "" {
		return 0, "", time.Duration(0), errors.New("data не может быть пустой строкой")
	}

	arrData := strings.Split(data, ",")
	if len(arrData) != 3 {
		return 0, "", time.Duration(0), errors.New("data ошибка формата")
	}

	name := arrData[1]

	steps, err := strconv.Atoi(arrData[0])
	if err != nil {
		return 0, "", time.Duration(0), err
	}

	if steps <= 0 {
		return 0, "", time.Duration(0), errors.New("шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(arrData[2])
	if err != nil {
		return 0, "", time.Duration(0), err
	}

	if duration.Milliseconds() <= 0 {
		return 0, "", time.Duration(0), errors.New("продолжительность должна быть больше 0")
	}

	return steps, name, duration, nil
}

func distance(steps int, height float64) float64 {
	lenStep := height * stepLengthCoefficient
	return (float64(steps) * lenStep) / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)
	return distance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, name, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)

	meanSpeed := meanSpeed(steps, height, duration)

	var calories float64
	switch name {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		err := errors.New("неизвестный тип тренировки")
		log.Println(err)
		return "", err
	}
	return fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f`, name, duration.Hours(), distance, meanSpeed, calories) + "\n", nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов не может быть 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес не может быть 0")
	}
	if height <= 0 {
		return 0, errors.New("рост не может быть 0")
	}
	if duration.Milliseconds() <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	return (weight * meanSpeed * duration.Minutes()) / float64(minInH), nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов не может быть 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес не может быть 0")
	}
	if height <= 0 {
		return 0, errors.New("рост не может быть 0")
	}
	if duration.Milliseconds() <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	return (weight * meanSpeed * duration.Minutes()) / float64(minInH) * walkingCaloriesCoefficient, nil
}
