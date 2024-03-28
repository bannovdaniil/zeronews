﻿# zeroNews
### Golang

Написать json REST сервер с двумя ручками:

- POST /edit/:Id - изменение новости по Id
- GET /list - список новостей
- БД - postgres

### Требования:

- В качестве сервера использовать fiber.
- Для работы с базой reform.
- Соединение с базой должно использовать connection pool.
- Все настройки через переменные окружения и/или viper.

Start application in docker:
переменные в .evn имеет приоритет над локальными из application.yml

> docker-compose down --volumes

> docker-compose build

> docker-compose up -d

Start application local:
переменные в application.yml

### Дополнительные реализации:
- валидацию полей при редактировании, 
  - проверка ID;
- docker
- логгирование logrus
- добавил fiber ErrorHandler
