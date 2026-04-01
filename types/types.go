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
	Balance        Balance `json:"balance" gorm:"column:balance"`            // баланс
	Currency       string  `json:"currency" gorm:"column:currency"`          // валюта
	Status         string  `json:"status" gorm:"column:card_status"`         // active / blocked
}
type CardResponse struct {
	IDAccount  int     `json:"id_account"`
	CardNumber string  `json:"card_number"`
	Holder     string  `json:"holder"`
	ExpMonth   int     `json:"exp_month"`
	ExpYear    int     `json:"exp_year"`
	CVV        string  `json:"cvv"`
	Balance    Balance `json:"balance"`
	Currency   string  `json:"currency"`
	Status     string  `json:"status"`
}
type Account struct {
	ID          int     `gorm:"column:id"             json:"id"`
	FirstName   string  `gorm:"column:first_name"     json:"first_name"`
	LastName    string  `gorm:"column:last_name"      json:"last_name"`
	DateOfBirth string  `gorm:"column:date_of_birth"  json:"date_of_birth"`
	PhoneNumber string  `gorm:"column:phone_number"   json:"phone_number"`
	Email       string  `gorm:"column:email"          json:"email"`
	Balance     Balance `gorm:"column:balance"        json:"balance"`
	Currency    string  `gorm:"column:currency"       json:"currency"`
	Password    string  `gorm:"column:password"       json:"password"`
	Status      string  `gorm:"column:account_status" json:"account_status"`
}
type BlockCardByPhone struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
type MoneyTransferAccountToAccount struct {
	PhoneNumberSender   string  `json:"phone_number_sender"`
	PhoneNumberReceiver string  `json:"phone_number_receiver"`
	PasswordSender      string  `json:"password_sender"`
	Amount              Balance `json:"amount"`
}
type MoneyTransferAccountToCard struct {
	PhoneNumberSender string  `json:"phone_number_sender"`
	CardReceiver      string  `json:"card_receiver"`
	PasswordSender    string  `json:"password_sender"`
	Amount            Balance `json:"amount"`
}
type MoneyTransferCardToAccount struct {
	CardSender          string  `json:"card_sender"`
	PhoneNumberReceiver string  `json:"phone_number_receiver"`
	CVVSender           string  `json:"cvv_sender"`
	Amount              Balance `json:"amount"`
}
type MoneyTransferCardToCard struct {
	CardSender   string  `json:"card_sender"`
	CardReceiver string  `json:"card_receiver"`
	CVVSender    string  `json:"cvv_sender"`
	Amount       Balance `json:"amount"`
}
