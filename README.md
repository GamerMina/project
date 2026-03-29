API Endpoints
Регистрация аккаунта

POST /register_account

{
  "first_name": "Anvar",
  "last_name": "Nazarov",
  "date_of_birth": "2003-05-12",
  "phone_number": "+992900111111",
  "email": "anvar@gmail.com",
  "password": "1234",
  "currency": "TJS"
}
Регистрация карты

POST /register_card

{
  "id": 1
}
Проверка наличия карты

GET /check

{
  "id": 1
}
Управление картой
Заблокировать карту

POST /block_card

{
  "phone_number": "+992900111111",
  "password": "1234"
}
Активировать карту

POST /activate_card

{
  "phone_number": "+992900111111",
  "password": "1234"
}
Переводы денег
Account to Account

POST /transfer/account-to-account

{
  "phone_number_sender": "+992900111111",
  "phone_number_receiver": "+992900222222",
  "password_sender": "1234",
  "amount": 100
}
Account to Card

POST /transfer/account-to-card

{
  "phone_number_sender": "+992900111111",
  "card_receiver": "4444888812345678",
  "password_sender": "1234",
  "amount": 50
}
Card to Account

POST /transfer/card-to-account

{
  "card_sender": "4444888812345678",
  "phone_number_receiver": "+992900222222",
  "cvv_sender": "123",
  "amount": 25
}
Card to Card

POST /transfer/card-to-card

{
  "card_sender": "4444888812345678",
  "card_receiver": "4444888877776666",
  "cvv_sender": "123",
  "amount": 75
}

Безопасность

В проекте реализовано:

Пароль хранится в виде bcrypt hash
CVV хранится в виде bcrypt hash
Номер карты хранится в виде HMAC SHA256 hash
Переводы выполняются через database transaction

Возможные доработки
JWT авторизация
История транзакций
Комиссии
Лимиты по карте
Конвертация валют

Автор

Anvar Nazarov
