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
- [ ] Подумать как можно добавить временной разрез для shows/clicks.
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

- [ ] Посмотреть какой еще алгоритм ротации можно добавить.
- [ ] Реализована отправка статистики в очередь
- [ ] Написаны интеграционные тесты
- [ ] Наличие валидных Dockerfile и Makefile для сервиса. Проект возможно собрать чере make build, запустить через make run и протестировать через make test


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