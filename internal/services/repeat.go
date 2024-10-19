package services

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"todo-rest/internal/config"
)

// NextDate вычисляет следующую дату задачи на основе правила повторения
// now — текущее время, date — начальная дата в формате "20060102", repeat — правило повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсим начальную дату задач
	taskDate, err := time.Parse(config.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	// Если repeat пуст, возвращаем ошибку
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	repeatParts := strings.Fields(repeat)

	switch repeatParts[0] {
	case "d":
		return nextDay(taskDate, now, repeatParts)
	case "y":
		return nextYear(taskDate, now)
	case "w":
		return nextWeekday(now, repeatParts)
	case "m":
		return nextMonth(taskDate, now, repeatParts)
	default:
		return "", fmt.Errorf("invalid repeat type: %s", repeatParts[0])
	}
}

// nextDay вычисляет следующую дату по правилу "d <дней>"
func nextDay(taskDate, now time.Time, repeatParts []string) (string, error) {
	if len(repeatParts) != 2 {
		return "", errors.New("invalid day repeat format")
	}

	days, err := strconv.Atoi(repeatParts[1])
	if err != nil || days > 400 {
		return "", errors.New("invalid day interval")
	}

	// Добавляем указанное количество дней, пока дата задачи не станет больше текущей
	taskDate = taskDate.AddDate(0, 0, days)
	for !taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	}

	return taskDate.Format(config.DateFormat), nil
}

// nextYear вычисляет следующую дату по правилу "y" (ежегодно)
func nextYear(taskDate, now time.Time) (string, error) {
	taskDate = taskDate.AddDate(1, 0, 0)
	for taskDate.Before(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}

	return taskDate.Format(config.DateFormat), nil
}

// nextWeekday вычисляет следующую дату по правилу "w <дни недели>"
func nextWeekday(now time.Time, repeatParts []string) (string, error) {
	if len(repeatParts) != 2 {
		return "", errors.New("invalid weekday repeat format")
	}

	// Получаем числовое значение текущего дня недели
	weekday := int(now.Weekday())
	// Преобразуем воскресенье (0) в 7
	if weekday == 0 {
		weekday = 7
	}

	// Делим строку с днями недели
	daysOfWeek := strings.Split(repeatParts[1], ",")
	repeatDays := make([]int, len(daysOfWeek))

	// Проверяем каждый день недели и конвертируем в int
	for i, day := range daysOfWeek {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt > 7 || dayInt < 1 {
			return "", errors.New("invalid day of the week")
		}
		if dayInt <= weekday {
			dayInt += 7 // Переносим день на следующую неделю
		}
		repeatDays[i] = dayInt
	}

	// Сортируем и вычисляем ближайшую дату повторения задачи
	sort.Ints(repeatDays)
	nextTaskDate := now.AddDate(0, 0, repeatDays[0]-weekday)
	return nextTaskDate.Format(config.DateFormat), nil
}

// nextMonth вычисляет следующую дату по правилу "m <дни> [месяцы]"
func nextMonth(taskDate, now time.Time, repeatParts []string) (string, error) {
	if len(repeatParts) < 2 || len(repeatParts) > 3 {
		return "", errors.New("invalid month repeat format")
	}

	// Проверяем дни месяца и конвертируем в int
	daysOfMonth := strings.Split(repeatParts[1], ",")
	dayInt := make([]int, 0, len(daysOfMonth))
	for _, day := range daysOfMonth {
		dInt, err := strconv.Atoi(day)
		if err != nil || dInt < -2 || dInt == 0 || dInt > 31 {
			return "", errors.New("invalid day of the month")
		}
		dayInt = append(dayInt, dInt)
	}

	sort.Ints(dayInt)

	// Получаем месяцы, если они указаны, проверяем и конвертируем в int
	var monthInt []int
	if len(repeatParts) == 3 {
		months := strings.Split(repeatParts[2], ",")
		for _, month := range months {
			mInt, err := strconv.Atoi(month)
			if err != nil || mInt < 1 || mInt > 12 {
				return "", errors.New("invalid month")
			}
			monthInt = append(monthInt, mInt)
		}
	} else {
		// Если месяцы не указаны, по умолчанию рассматриваются все
		for i := 1; i <= 12; i++ {
			monthInt = append(monthInt, i)
		}
	}

	for {
		if !slices.Contains(monthInt, int(taskDate.Month())) {
			taskDate = taskDate.AddDate(0, 1, 0)
			if taskDate.Day() > 1 {
				taskDate = taskDate.AddDate(0, 0, -taskDate.Day()+1)
			}
			continue
		}

		validDays := validDaysInMonth(taskDate, dayInt)
		currentMonth := taskDate.Month()
		for {
			if currentMonth != taskDate.Month() {
				break
			}
			if slices.Contains(validDays, taskDate.Day()) && taskDate.After(now) {
				return taskDate.Format(config.DateFormat), nil
			}
			taskDate = taskDate.AddDate(0, 0, 1)
		}
	}
}

// validDaysInMonth создает список допустимых дней для указанного месяца
func validDaysInMonth(taskDate time.Time, dayInt []int) []int {
	daysInMonth := time.Date(taskDate.Year(), taskDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	result := make([]int, 0, len(dayInt))
	for _, d := range dayInt {
		if d > daysInMonth {
			continue
		}
		if d > 0 {
			result = append(result, d)
			continue
		}
		result = append(result, daysInMonth+d+1)
	}
	return result
}
