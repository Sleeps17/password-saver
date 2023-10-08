# password-saver
Это сервер, который позволяет сохранить ваши пароли от каких либо сайтов.<br>
Для использования нужно:
- Склонировать репозиторий к себе на компьютер.
- Установить компилятор Go.
- Зайти в папку с проектом.
- Сделать ``` go run ./cmd/main.go```

Возможно сделать 5 действий:
## Метод save
Сохраняет ваш пароль.<br>
Параметры: url, alias, password.<br>
url - обязательный параметр, alias и password - необязательные. <br>
Если alias или password не будут указаны или будут указаны пустыми, то их значения сгенерируются автоматически.<br>
Пример запроса: ``` curl -X POST -H "Content-Type application/json" -d '{"url": "https://github.com", "alias": "github", "password": "1234567890"} http://localhost:8082/saver/save'```
## Метод get
Выдает вам пароль по ключу. Ключом может быть url или alias сайта, как вам удобно.<br>
Параметры: key - alias или url.<br>
key - обязательный параметр.<br>
Пример запроса: ``` curl -X POST -H "Content-Type application/json" -d '{"key": "https://github.com"} http://localhost:8082/saver/get'```<br>
## Метод getAll
Возвращает вам все пароли, которые сейчас хранятся.<br>
Параметры: -<br>
Пример запроса: ``` curl -X POST -H "Content-Type application/json" http://localhost:8082/saver/getAll'```
## Метод reset
Изменяет пароль по ключу. Ключом может быть alias или url сайта, как вам удобно.<br>
Параметры: key, new_password.<br>
key - обязательный параметр, new_password - необязательный. <br>
Если new_password не будет указан или будет указан пустым, то он сгенерируется автоматически.<br>
Пример запроса: ``` curl -X POST -H "Content-Type application/json" -d '{"key": "github", "new_password": "0987654321"} http://localhost:8082/saver/reset'```
## Метод delete
Удаляет пароль по ключу. Ключом может быть url или alias сайта, как вам удобно.<br>
Параметры: key.<br>
key - обязательный параметр.<br>
Пример запроса: ``` curl -X POST -H "Content-Type application/json" -d '{"key": "github"} http://localhost:8082/saver/delete'```
