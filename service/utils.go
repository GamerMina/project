package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

// Проверка по алгоритму Luhn
func (s *Services) ValidLuhn(number string) (string, error) {
	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		n := int(number[i] - '0') // а также можно минус 48 сделать потомучто ASCII таблица
		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alternate = !alternate
	}

	if sum != 0 {
		return "", errors.New("проверку Luhn не пройденна ")
	}
	return number, nil
}

// Генерация  цифры Luhn
func (s *Services) generateLuhnCheckDigit(number string) (int, error) {
	sum := 0
	alternate := true
	for i := len(number) - 1; i >= 0; i-- {
		n := int(number[i] - 48) //-48 потомучто мы в ASCII таблицу вытащит а нам этого не надо
		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alternate = !alternate
	}
	luhnDigit := (10 - (sum % 10)) % 10
	return luhnDigit, nil
}

// Генерация 15-значного номера карты
func (s *Services) generateCardNumber() (string, error) {
	rand.Seed(time.Now().UnixNano())

	var number string

	// Генерируем первые 15 цифр
	for i := 0; i < 15; i++ {
		number += fmt.Sprintf("%d", rand.Intn(10))
	}
	println("Сгенерированный номер карты:", number)
	// Вычисляем контрольную цифру
	checkDigit, err := s.generateLuhnCheckDigit(number)
	if err != nil {
		return "", err
	}
	num := number + fmt.Sprintf("%d", checkDigit)
	num, err = s.ValidLuhn(num)

	return num, nil
}

// HashCardData Хеширование
func HashCardData(pan, secret string) string {

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(pan))

	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
func GenerateCVV() string {
	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(900) + 100 // 100–999
	numberString := strconv.Itoa(number)
	return numberString
}

// HashCVV Хеширование для CVV через bcrypt
func HashCVV(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), 4) // просто чтобы быстро работало
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareHash проверяет хеш от bcrypt сравнивает настоящий хеш_пароля с паролем через bcrypt
func CompareHash(hash string, code string) (bool, error) {
	check := bcrypt.CompareHashAndPassword([]byte(hash), []byte(code))
	if check != nil {
		return false, errors.New("password is incorrect ")
	}
	return true, nil
}

func HidePAN(s string) string {
	runes := []rune(s)
	digitIndex := 0

	for i := 0; i < len(runes); i++ {
		if runes[i] == ' ' {
			continue
		}

		digitIndex++

		if digitIndex >= 7 && digitIndex <= 12 {
			runes[i] = '*'
		}
	}
	hidenPAN := string(runes)
	return hidenPAN
}
func AddYearsMonths(years int, months int) (int, int) {
	now := time.Now()
	newDate := now.AddDate(years, months, 0)
	year := newDate.Year() % 100  // последние 2 цифры года
	month := int(newDate.Month()) // месяц от 1 до 12
	return year, month
}
