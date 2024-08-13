package models

type Accrual struct {
	MaxWorker            int    //Количество рабочих процессов горутин
	TimeTicker           int    //Количество секунд повторения запроса на обновление статуса необработанных заказов
	MaxRetries           int    //Количество попыток получить данные от сервиса accrual в случае ошибки 429
	AccrualServerAddress string //Адрес сервера accrual
}

type AccrualOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}
