package types

type DbConf struct {
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string
}

type Card struct {
	ID             int     `json:"id"`
	IDAccount      int     `json:"id_account"`       // ID к которому она привязана
	CardNumber     string  `json:"card_number" `     // номер карты (16 цифр)
	CardNumberHash string  `json:"card_number_Hash"` // кеш номера карты сгенерированный из даты пана и cvv
	Holder         string  `json:"holder" `          // владелец карты на латинском
	ExpMonth       int     `json:"exp_month" `       // месяц окончания
	ExpYear        int     `json:"exp_year" `        // год окончания
	CVV            string  `json:"cvv" `             // CVV код
	CVVHash        string  `json:"cvv_Hash" `        // CVV код который мы отправляем в БД
	Balance        float64 `json:"balance" `         // баланс
	Currency       string  `json:"currency" `        // валюта
	Status         string  `json:"status" `          // active / blocked
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
