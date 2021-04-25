![tests](https://github.com/FedoseevAlex/banner-rotation/actions/workflows/tests.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/FedoseevAlex/banner-rotation)](https://goreportcard.com/report/github.com/FedoseevAlex/banner-rotation)
--------------

# banner-rotation
This is a final project for course OTUS Golang Professional

## Чек лист проекта
- [ ] Ветка master успешно проходит пайплайн в CI-CD системе (на ваш вкус, Circle CI, Travis CI, Jenkins, GitLab CI и пр.).
- [x] Добавить в CI запуск последней версии golangci-lint на весь проект с конфигом, представленным в данном репозитории.
- [x] Добавить в CI запуск юнит тестов командой вида go test -race -count 100.
- [ ] Настроить бэйджи для визуализации состояния кода репы.
- [x] Реализован алгоритм "многорукого бандита" - 2 балла.
    Можно прочитать поподробнее [тут](https://lilianweng.github.io/lil-log/2018/01/23/the-multi-armed-bandit-problem-and-its-solutions.html), [или тут](https://www.optimizely.com/optimization-glossary/multi-armed-bandit/), [полное погружение](https://arxiv.org/pdf/1904.07272.pdf)
- [x] Наличие юнит-тестов на алгоритм многорукого бандита.
    - [x] Проверка на перебор всех баннеров. Каждый баннер должен быть показан хотя бы один раз.
    - [x] Проверка, что у популярного баннера должно быть наибольшее число показов.
- [ ] Реализовано разделение на "слоты" и "соц.дем. группы" - 2 балла.
- [ ] Реализовано API сервиса - 2 балла.
- [ ] Реализована отправка статистики в очередь - 1 балл.
- [ ] Написаны юнит-тесты - 1 балл.
- [ ] Написаны интеграционные тесты - 2 балла.
- [ ] Тесты адекватны и полностью покрывают фукнционал - 1 балл.
- [ ] Наличие валидных Dockerfile и Makefile для сервиса. Проект возможно собрать чере make build, запустить через make run и протестировать через make test - 1 балл.
- [ ] Понятность и чистота кода - до 3 баллов.

Максимум 15 баллов, зачет от 10.