package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/mail"
	"project/types"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
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
	//TODO: пересмотреть обработку ошибок
	return num, nil
}

// HashData Хеширование
func HashData(text, secret string) string {

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(text))

	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

// GenerateCVV создание СВВ
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

// HidePAN Hide PAN
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

// AddYearsMonths years that months добовляет в функцию
func AddYearsMonths(years int, months int) (int, int) {
	now := time.Now()
	newDate := now.AddDate(years, months, 0)
	year := newDate.Year() % 100  // последние 2 цифры года
	month := int(newDate.Month()) // месяц от 1 до 12
	return year, month
}

// ParcePhoneNumber парсим номер телефона
func ParcePhoneNumber(input string) string {
	re := regexp.MustCompile(`\D`)
	number := re.ReplaceAllString(input, "")
	if number == "" {
		return ""
	}
	if strings.HasPrefix(number, "992") {
		return "+" + number
	}
	if len(number) > 9 {
		number = number[len(number)-9:]
	}
	return "+992" + number
}

// ParseName Parse name
func ParseName(name, surname string) error {
	if IsValidName(name) == false {
		return errors.New("name contain some symbols")
	}
	if IsValidName(surname) == false {
		return errors.New("name contain some symbols")
	}
	return nil
}

// IsValidName  парсит стринг
func IsValidName(name string) bool {
	for _, r := range name {
		if !(unicode.IsLetter(r) || r == ' ' || r == '-') {
			return false
		}
	}
	return true
}

// ParseMail Parcs mail
func ParseMail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// GetAge из точной даты получаем сколько ему полных лет на данный момент
func GetAge(DateOfBirth string) (int, error) {
	// парсим дату
	birthDate, err := time.Parse("2006-01-02", DateOfBirth)
	if err != nil {
		return 0, fmt.Errorf("неверный формат даты: %v", err)
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// проверяем был ли ДР в этом году
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age, nil
}

// TransferMoney Transfer balance
func TransferMoney(balanceFirst types.Balance, balanceSecend types.Balance, amount types.Balance) (types.Balance, types.Balance, error) {
	var err error
	if balanceFirst < amount {
		err = errors.New("declined: insufficient funds")
		return 0, 0, err
	}
	balanceFirst -= amount
	balanceSecend += amount
	return balanceFirst, balanceSecend, err
}
