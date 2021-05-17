![tests](https://github.com/FedoseevAlex/banner-rotation/actions/workflows/tests.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/FedoseevAlex/banner-rotation)](https://goreportcard.com/report/github.com/FedoseevAlex/banner-rotation)
--------------

# banner-rotation
This is a final project for course OTUS Golang Professional

## Чек лист проекта
- [x] Ветка master успешно проходит пайплайн в CI-CD системе (на ваш вкус, Circle CI, Travis CI, Jenkins, GitLab CI и пр.).
- [x] Добавить в CI запуск последней версии golangci-lint на весь проект с конфигом, представленным в данном репозитории.
- [x] Добавить в CI запуск юнит тестов командой вида go test -race -count 100.
- [x] Настроить бэйджи для визуализации состояния кода репы.
- [x] Реализован алгоритм "многорукого бандита"
    Можно прочитать поподробнее [тут](https://lilianweng.github.io/lil-log/2018/01/23/the-multi-armed-bandit-problem-and-its-solutions.html), [или тут](https://www.optimizely.com/optimization-glossary/multi-armed-bandit/), [полное погружение](https://arxiv.org/pdf/1904.07272.pdf)
- [x] Наличие юнит-тестов на алгоритм многорукого бандита.
    - [x] Проверка на перебор всех баннеров. Каждый баннер должен быть показан хотя бы один раз.
    - [x] Проверка, что у популярного баннера должно быть наибольшее число показов.
- [x] Реализован слой работы с базой данных
- [x] Написаны тесты для проверки работы с базой данных
- [x] Реализовано разделение на "слоты" и "соц.дем. группы"
- [x] Добавить интерфейс для многорукого бандита в пакет с основными типами.
- [x] Реализовано API сервиса
- [ ] Написаны юнит-тесты
    - [x] Написаны тесты на логику базы данных
    - [x] Написаны тесты на многорукого бандита
    - [ ] Написаны тесты на приложение
    - [ ] Написаны тесты на http сервер
- [ ] Написаны интеграционные тесты
- [x] Подумать как можно добавить временной разрез для shows/clicks.
    Можно добавить таблицу, в которой будут храниться timestamps для rotations с типом события:
    Название: event_timestamps
    | column        | type                      | modifiers     |
    |---------------|---------------------------|---------------|
    |rotation_id    | serial                    | primary key   |
    |stamp          | timestamp without timezone| not null      |
    |event_type     | text                      | not null      |

    Для добавления такой таблицы нужно будет изменить таблицу rotations:
    | column     | type                        | modifiers      |
    |------------|-----------------------------|----------------|
    | id         | serial                      |  primary key   |
    | banner_id  | uuid                        |  not null      |
    | slot_id    | uuid                        |  not null      |
    | group_id   | uuid                        |  not null      |
    | shows      | integer                     |                |
    | clicks     | integer                     |                |
    | deleted    | boolean                     |  default false |
    | deleted_at | timestamp without time zone |                |

    Так как в rotations был составной primary key (banner_id, slot_id, group_id), то при введении поля id нужно будет
    добавить unique constraint на (banner_id, slot_id, group_id).
    Можно еще добавить индекс для полей (banner_id, slot_id, group_id), но это может быть избыточно сейчас. Плюс, его можно будет просто добавить позже.
- [ ] Реализована отправка статистики в очередь
- [x] Проект возможно собрать чере make build, запустить через make run и протестировать через make test
    - [x] Валидный Makefile
- [ ] Валидный Dockerfile

## Замечания после первой проверки
На данный момент набрано 5 баллов

- [x] Оформить дз как pull-request из develop в master.
- [x] Переименовать файл в cmd/rotator/rotator.go в main.go.
- [x] Передавать в printVersion аргумент, реализующий io.Writer, чтобы можно было универсально использовать.
- [x] Вынести структуры с тегами базы данных в отдельный файл models.go.
- [x] Выделить основные используемые в приложении типы и интерфейсы в отдельный пакет.
- [x] Убрать каскадное удаление данных из rotations.
- [x] Добавить поля deleted и deleted_at и изменить код удаления объектов из базы. Статистику желательно оставлять в базе и просто помечать строки как удаленные.

Максимум 15 баллов, зачет от 10.

## Примеры запросов

### Версия приложения
Получить информацию о работающей версии приложения.  
URL: `/version`  
METHOD: `GET`  
Request:  
```
curl --location --request GET 'localhost:8080/version'
```
Response:  
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 16 May 2021 19:30:26 GMT
Content-Length: 109

{"Release":"develop","BuildDate":"2021-05-16T18:13:29","GitHash":"e97669754c6ca54e500c1986709b5e9629a2d5d3"}
```
### Баннеры, слоты и группы
Эндпоинты для работы с баннерами, слотами и группами.  
В приложении для идентификации баннера слота или группы используются UUID v.4.  

#### Создание баннера, слота или группы
URL: `/{banners|slots|groups}`  
METHOD: `POST`  
При создании обязательно передавать описание в теле запроса.  
Request:  
```
curl --location --request POST 'localhost:8080/{banners|slots|groups}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "description": "{banner|slot|group} created from api"
}'
```
Response:  
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 16 May 2021 19:00:50 GMT
Content-Length: 85

{"ID":"0beac2d5-05dd-4bca-9052-9ccb11a715b3","Description":"{banner|slot|group} created from api"}
```

#### Получение баннера,слота или группы
URL: `/{banners|slots|groups}/{:banner_id|:slot_id|:group_id}`  
METHOD: `GET`  
Request:  
```
curl --location --request GET 'localhost:8080/{banners|slots|groups}/0beac2d5-05dd-4bca-9052-9ccb11a715b3'
```
Response:  
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 16 May 2021 19:01:49 GMT
Content-Length: 85

{"ID":"0beac2d5-05dd-4bca-9052-9ccb11a715b3","Description":"{banner|slot|group} created from api"}
```

Если был передан невалидный UUID, то будет выдана ошибка.  
Request:  
```
curl --location --request GET 'localhost:8080/{banners|slots|groups}/0beac2d5-05dd-4bca-9052-9ccb11aasdfasdf'
```
Response:  
```
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Sun, 16 May 2021 19:03:09 GMT
Content-Length: 71

{"Error":"invalid UUID length: 39","Msg":"failed to parse {banner|slot|group} uuid"}
```

#### Удаление баннера, cлота или группы
URL: `/{banners|slots|groups}/{:banner_id|:slot_id|:group_id}`  
METHOD: `DELETE`  
Request:  
```
curl --location --request DELETE 'localhost:8080/{banners|slots|groups}/0beac2d5-05dd-4bca-9052-9ccb11aasdfasdf'
```
Response:  
```
HTTP/1.1 204 No Content
Content-Type: application/json
Date: Sun, 16 May 2021 19:09:02 GMT
```

### Ротации
Эндпоинты для управления ротациями

#### Создать ротацию
Создать связь между соц. дем. группой, слотом и баннером.  
URL: `/group/:group_id/slots/:slot_id/banner/:banner_id`  
METHOD: `POST`  
Request:  
```
curl --location --request POST 'localhost:8080/group/649647a7-6c4c-4044-843f-e48a9748ab90/slots/cc8a98c0-80a6-4e34-b8db-f5377c2897bf/banners/0beac2d5-05dd-4bca-9052-9ccb11a715b3'
```
Response:  
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 16 May 2021 19:20:42 GMT
Content-Length: 169

{"BannerID":"c511c792-a880-4a86-93da-239b12bb6b3e","SlotID":"99165522-e304-4dfc-95e3-1fe326c48f6e","GroupID":"493148ec-0b08-4eb8-afd1-60b608a6a6d2","Shows":0,"Clicks":0}
```

#### Выбрать баннер
Выбрать баннер для отображения данной группе в указанном слоте.  
При передаче запроса в этот эндпоинт баннеру автоматически увеличивается количество показов.  
URL: `/group/:group_id/slots/:slot_id/banner`  
METHOD: `GET`  
Request:  
```
curl --location --request GET 'localhost:8080//group/493148ec-0b08-4eb8-afd1-60b608a6a6d2/slots/99165522-e304-4dfc-95e3-1fe326c48f6e/banner'
```
Response:  
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 16 May 2021 19:20:42 GMT
Content-Length: 169

{"BannerID":"c511c792-a880-4a86-93da-239b12bb6b3e","SlotID":"99165522-e304-4dfc-95e3-1fe326c48f6e","GroupID":"493148ec-0b08-4eb8-afd1-60b608a6a6d2","Shows":2,"Clicks":0}
```

#### Регистрация перехода
Зафиксировать факт перехода по баннеру.  
URL: `/group/:group_id/slots/:slot_id/banner/:banner_id/click`  
METHOD: `POST`  
Request:  
```
curl --location --request GET 'localhost:8080//group/493148ec-0b08-4eb8-afd1-60b608a6a6d2/slots/99165522-e304-4dfc-95e3-1fe326c48f6e/banner'
```
Response:  
```
HTTP/1.1 204 No Content
Content-Type: application/json
Date: Sun, 16 May 2021 19:27:31 GMT
```

#### Получение статистики по ротации
Статистика выдается в виде массива с событиями.  
Событие имеет два поля: тип (click или show) и временную метку (timestamp).  
URL: `/group/:group_id/slots/:slot_id/banner/:banner_id/stats`  
METHOD: `GET`  
Request:  
```
curl --location --request GET 'localhost:8080/group/493148ec-0b08-4eb8-afd1-60b608a6a6d2/slots/99165522-e304-4dfc-95e3-1fe326c48f6e/banners/c511c792-a880-4a86-93da-239b12bb6b3e/stats'
```
Response:  
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 16 May 2021 19:34:28 GMT
Content-Length: 410

[{"Type":"show","Timestamp":"2021-05-16T19:23:31.727697Z"},{"Type":"show","Timestamp":"2021-05-16T19:23:35.14956Z"},{"Type":"show","Timestamp":"2021-05-16T19:24:46.918126Z"},{"Type":"click","Timestamp":"2021-05-16T19:26:55.950379Z"},{"Type":"click","Timestamp":"2021-05-16T19:27:11.941678Z"},{"Type":"click","Timestamp":"2021-05-16T19:27:17.765968Z"},{"Type":"click","Timestamp":"2021-05-16T19:27:31.827802Z"}]
```