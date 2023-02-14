# P2P.opg interview task

## Description
Написать экспортер метрик (лучше всего на golang, но если не умеете, то можно на python) забирающий из ноды текущий номер блока, рассинхрон времени этого блока в секундах (текущее время в секундах минус время блока в секундах, чтобы было видно здоровье этой ноды) и количество пиров. Т.е. три метрики всего.

### Текущий номер блока
Забираем вес последнего блока через REST API по пути `cosmos/base/tendermint/v1beta1/blocks/latest` из ключа `Block.Header.Height`

### Рассинхрон времени этого блока в секундах
Забираем время синхронизации последнего блока через REST API по пути `cosmos/base/tendermint/v1beta1/blocks/latest` из ключа `Block.Header.Time` и вычитаем текущее время для вычисления разницы времени синхронизации `Block.Header.Time.Unix() - time.Now().Unix()`

### Количество пиров
Я что-то в документации не нашел термина peer, и не могу понять, откуда забирать значение.

## How to run

Для отладки и запуска локально я использовал проброс портов через SSH:

```
ssh -L 11317:localhost:1317 ubuntu@43.131.34.230
go run main.go
curl localhost:113717/metrics

# HELP cosmos_latest_block_height The latest block id hash
# TYPE cosmos_latest_block_height untyped
cosmos_latest_block_height 118815
# HELP cosmos_latest_block_timestamp Unsync node in ms
# TYPE cosmos_latest_block_timestamp untyped
cosmos_latest_block_timestamp -9
```

