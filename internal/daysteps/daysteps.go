package daysteps

import (
	"errors"
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
	if data == "" {
		return 0, time.Duration(0), errors.New("data не может быть пустой строкой")
	}

	arrData := strings.Split(data, ",")
	if len(arrData) != 2 {
		return 0, time.Duration(0), errors.New("data ошибка формата")
	}

	steps, err := strconv.Atoi(arrData[0])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if steps <= 0 {
		return 0, time.Duration(0), errors.New("количество шагов не может быть 0")
	}

	duration, err := time.ParseDuration(arrData[1])
	if err != nil {
		return 0, time.Duration(0), err
	}

	if duration.Milliseconds() <= 0 {
		return 0, time.Duration(0), errors.New("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	distation := (stepLength * float64(steps)) / float64(mInKm)

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf(`Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.`, steps, distation, calories) + "\n"
}
