# Калькулятор-сервер

Простой калькулятор, работающий через простой однопоточный сервер. ОБрабатывает знаки `+-*/`, скобки, рациональные числа через `.`.<br>

## Начало работы
Для запуска сервера используйте команду из корневой директории проекта:

```
go run cmd/main.go
```

По умолчанию сервер будет использовать порт `8080`, вы можете указать свой:

```
go run cmd/main.go 4321
```

Для доступа непосредственно к функционалу используйте любимый инструмент, который может отправить HTTP-запрос следующего содержания:

```
POST /api/v1/calculate HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{"expression": "52"}
```

### Пример команды curl

Windows cmd:
```cmd
curl -X POST -H "Content-Type: application/json" -d "{\"expression\":\"1\"}" localhost:8080/api/v1/calculate
```
Bash:
```bash
curl -X POST \
-H "Content-Type: application/json" \
-d '{"expression":"1"}' \
localhost:8080/api/v1/calculate
```

### Пример скрипта на Python:

```python
import requests

url = 'http://localhost:8080/api/v1/calculate'
headers = {'Content-Type': 'application/json'}
data = {'expression': '1+2 - 3 / (5 - 6)'}

r = requests.post(url, headers=headers, json=data)
print(r.json())
```

## Возможные ответы сервера

В случае правильного запроса вы получите:
(код 200)

```
{
    "result": "результат выражения"
}
```

В случае ошибки в запросе или ключе `expression` (напр `"шышел-мышел пернул вышел"`):<br>
(код 422)

```
{
    "error": "Подробности касаемо валидности данных"
}
```

В случае непредвиденной ошибки сервера вы ничего не получите, точнее:<br>
(код 500)

```
{
    "error": "Internal server error"
}
```

## Принцип работы

### Калькулятор

Работает с рациональными числами. Находится в `pkg/calculator/calculator.go`. Далее многобукв про работу кода.<br>При получении строки с выражением, она сначала проверяется на разорванные числа, например `"1 2 + 3   . 4"`. Затем из строки убираются пробелы и выражение валидируется по полной: на оссмысленность, последовательность скобок, операторы. Стоит отметить, что последовательности операторов, вроде `"1 + -2"` не проходят валидацию, а `1/-(2)` проходят -- так же работает яндексовский калькулятор (по запросу "калькулятор"). Далее строка разбивается на *слайс интерфейсов* из чисел (в том числе *отрицательных*) и символов `()+-*/` и калькулятор работает с этой последовательностью. Он ищет скобки и сначала считает их, возвращая на место выражения скобок результат (задействуется рекурсия). Стоит отметить, что для этого я каждый раз создаю копии слайса, т.е. аллоцирую их, вместо работы над одним. Непосредственно операции `+-*/` выполняются в порядке приоритета. Также `-` может быть самым приоритетным, если он стоит перед выражением в скобках (к этому моменту результату этого выражения). На этом всё.

### HTTP-сервер
Работает с 1 клиентом в 1 время. <br>
Разбит по следующей структуре.

```
internal/
|--application/
  |--application.go      -- обслуживание хэндлеров
  |--application_test.go -- тесты
  |--httpHandlers.go     -- хэнддеры
  |--utils.go            -- вспомогательные функции
  |--errors.go           -- ошибки для сервера (отличные от errors.go в calculator)
```

Доступен по `localhost:8080/api/v1/calculate` (подробнее см. пункт 1). По сути состоит из хэндлера `CalcHandler` и 2-х мидлварей (middleware) `Validation`, `Outer`. При правильном запросе он направляется в хэндлер, и если что-то не прошло, например, валидацию калькулятора, то возвращается ошибка в `defer` предыдущей мидлвари или миддлваря (middleware). Сам `Validation` обрабатывает ошибку валидации и возвращает ее со статусом `422`, а если что-то совсем пошло не так (ошибка тут или непредопределена), то ошибка поднимается в `Outer` и там записывается в лог файл `errors.log`, который создается в корне проекта. Клиент же получит статус `500` и соответствующий ответ. Такую ошибку придется дебажить тому кто это делал.
