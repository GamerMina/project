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

// ValidLuhn Проверка по алгоритму Luhn
func (s *Services) ValidLuhn(number string) (string, error) {
	sum := 0
	alternate := false
	for i := len(number) - 1; i >= 0; i-- {
		n := int(number[i] - '0')
		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alternate = !alternate
	}
	if sum%10 != 0 {
		return "", errors.New("проверка Luhn не пройдена")
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var number string

	for i := 0; i < 15; i++ {
		number += fmt.Sprintf("%d", r.Intn(10))
	}

	checkDigit, err := s.generateLuhnCheckDigit(number)
	if err != nil {
		return "", err
	}

	num := number + fmt.Sprintf("%d", checkDigit)

	_, err = s.ValidLuhn(num)
	if err != nil {
		return "", err
	}

	return num, nil
}

// GenerateCVV создание СВВ
func GenerateCVV() string {
	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(900) + 100 // 100–999
	numberString := strconv.Itoa(number)
	return numberString
}

// Детерминированный хеш (для PAN, phone, поиск)
func HashData(text, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// Хеш пароля
func HashPassword(text, secret string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(text+secret), bcrypt.DefaultCost)
	return string(hash)
}

// Проверка пароля
func ComparePassword(hash, text, secret string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text+secret)) == nil
}

// Хеш CVV
func HashCVV(cvv string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Проверка CVV
func CompareCVV(hash, cvv string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(cvv)) == nil
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

// TransferMoney Transfer баланс именно 1 фунцией чтобы не было проблем по типо того что на 1 аккауете деньги снялись а на 2 не прибавились при каких либо проблемах с нашей стороны
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

func ValidateTransfer(sender string, receiver string, amount types.Balance) error {
	if sender == "" {
		return errors.New("invalid sender")
	}
	if receiver == "" {
		return errors.New("invalid receiver")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if sender == receiver {
		return errors.New("cannot transfer to yourself")
	}
	return nil
}

// ValidateTransactionAfterGettingData статус сендер статус ресивер, валюта сендер валюта ресивер для проверок
func ValidateTransactionAfterGettingData(sender string, receiver string, curencySender string, curencyReceiver string) error {
	if sender != "active" {
		return errors.New("sender account or card is not active")
	}
	if receiver != "active" {
		return errors.New("receiver account or  card is not active")
	}
	if curencySender != curencyReceiver {
		return errors.New("currencies do not match")
	}
	return nil
}
