Тестовое задание:

По Websocket получать данные о состоявшихся сделках с биржи Poloniex

API описан тут https://docs.poloniex.com/#price-aggregated-book

На входе словарь в виде json-строки:

{"poloniex":["BTC_USDT", "TRX_USDT", "ETH_USDT"]}

Надо подписаться на данные пары(в их API они перевёрнутые, например USDT_BTC)

На выходе(лог в консоль) должна быть структура:

type RecentTrade struct {
    Id        string    // ID транзакции
    Pair      string    // Торговая пара (из списка выше)
    Price     float64   // Цена транзакции
    Amount    float64   // Объём транзакции
    Side      string    // Как биржа засчитала эту сделку (как buy или как sell)
    Timestamp time.Time // Время транзакции
}

Не использовать готовые клиенты (библиотеки) для работы с API полоникса, разрешено использовать стандартные библиотеки для работы с ws

Пример запуска:
`go build -race . && ./price_book_test '{"poloniex":["BTC_USDT", "TRX_USDT", "ETH_USDT", "COMBO_USDT"]}'`
