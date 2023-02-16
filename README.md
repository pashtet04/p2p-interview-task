# P2P.opg interview task

## Description
Написать экспортер метрик (лучше всего на golang, но если не умеете, то можно на python) забирающий из ноды текущий номер блока, рассинхрон времени этого блока в секундах (текущее время в секундах минус время блока в секундах, чтобы было видно здоровье этой ноды) и количество пиров. Т.е. три метрики всего.

### Текущий номер блока
Забираем вес последнего блока через REST API по пути `cosmos/base/tendermint/v1beta1/blocks/latest` из ключа `Block.Header.Height`

### Рассинхрон времени этого блока в секундах
Забираем время синхронизации последнего блока через REST API по пути `cosmos/base/tendermint/v1beta1/blocks/latest` из ключа `Block.Header.Time` и вычитаем его из текущего времени для вычисления разницы времени синхронизации `time.Now().Unix() - Block.Header.Time.Unix()`

### Количество пиров
Я что-то в документации не нашел термина peer, и не могу понять, откуда забирать значение.

## How to run

- Для указания Cosmos API необходимо задать переменную окружения `COSMOS_API`, по-умолчанию `http://localhost:113717/`
- Для указания Tendermint API необходимо задать переменную окружения `TENDERMINT_API`, по-умолчанию `http://localhost:26657/`

Для отладки и запуска локально я использовал проброс портов через SSH:

```
ssh -L 11317:localhost:1317 ubuntu@43.131.34.230
go run main.go
curl localhost:113717/metrics

# HELP cosmos_latest_block_height The latest block id hash
# TYPE cosmos_latest_block_height untyped
cosmos_latest_block_height 150877
# HELP cosmos_latest_block_timestamp Unsync node in ms
# TYPE cosmos_latest_block_timestamp untyped
cosmos_latest_block_timestamp 8
# HELP cosmos_number_of_peers Number of peers
# TYPE cosmos_number_of_peers untyped
cosmos_number_of_peers 0
```

