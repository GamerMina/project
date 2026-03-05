package service

import (
	"fmt"
	"math/rand"
	"projcet/repository"
	"time"
)

type Services struct {
	Repository *repository.Repository
}

func (s *Services) NewServices(rep *repository.Repository) *Services {
	return &Services{Repository: rep}
}

// Проверка по алгоритму Luhn
func ValidLuhn(number string) bool {
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
	return sum%10 == 0
}

// Генерация  цифры Luhn
func (s *Services) generateLuhnCheckDigit(number string) int {
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

	return (10 - (sum % 10)) % 10
}

// Генерация 16-значного номера карты
func (s *Services) generateCardNumber() string {
	rand.Seed(time.Now().UnixNano())

	var number string

	// Генерируем первые 15 цифр
	for i := 0; i < 15; i++ {
		number += fmt.Sprintf("%d", rand.Intn(10))
	}
	println("Сгенерированный номер карты:", number)
	// Вычисляем контрольную цифру
	checkDigit := s.generateLuhnCheckDigit(number)

	return number + fmt.Sprintf("%d", checkDigit)
}
