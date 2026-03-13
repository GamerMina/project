package types

type DbConf struct {
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string
}
type Card struct {
	ID             int     `json:"id" gorm:"column:id"`
	IDAccount      int     `json:"id_account" gorm:"column:account_id"`      // ID к которому она привязана
	CardNumber     string  `json:"card_number" gorm:"-"`                     // номер карты (16 цифр)
	CardNumberHash string  `json:"card_number_hash" gorm:"column:pan_hash"`  // кеш номера карты сгенерированный из даты пана и cvv
	Holder         string  `json:"holder" gorm:"column:cardholder_name"`     // владелец карты на латинском
	ExpMonth       int     `json:"exp_month" gorm:"column:expiration_month"` // месяц окончания
	ExpYear        int     `json:"exp_year" gorm:"column:expiration_year"`   // год окончания
	CVV            string  `json:"cvv" gorm:"-"`                             // CVV код
	CVVHash        string  `json:"cvv_hash" gorm:"column:cvv_hash"`          // CVV код который мы отправляем в БД
	CardType       string  `json:"card_type" gorm:"column:card_type"`
	Balance        float64 `json:"balance" gorm:"column:balance"`    // баланс
	Currency       string  `json:"currency" gorm:"column:currency"`  // валюта
	Status         string  `json:"status" gorm:"column:card_status"` // active / blocked
}
type Account struct {
	ID            int64   `gorm:"id"             json:"id"`             //id
	FirstName     string  `gorm:"st_name"        json:"first_name"`     //имя на кириллице
	LastName      string  `gorm:"last_name"      json:"last_name"`      //фамилия на кириллице
	DateOfBirth   string  `gorm:"date_of_birth"  json:"date_of_birth"`  // Др
	PhoneNumber   string  `gorm:"phone_number"   json:"phone_number"`   //номер телефона в виде 12 чисел с +ом
	Email         string  `gorm:"email"          json:"email"`          // gmail
	Balance       float64 `gorm:"balance"        json:"balance"`        // баланс в виде 15,2 тоесть 13 чисел и 2 после зяпятой
	Currency      string  `gorm:"currency"       json:"currency"`       //TJS USD EUR
	Password      string  `gorm:"password"       json:"password"`       // пароль   4 значный
	AccountStatus string  `gorm:"account_status" json:"account_status"` // active или blocked
}
