# banner-rotation
Ротация баннеров. Проектная работа

[![Go Report Card](https://goreportcard.com/badge/github.com/usmartpro/banner-rotation)](https://goreportcard.com/report/github.com/usmartpro/banner-rotation)
![Tests](https://github.com/usmartpro/banner-rotation/actions/workflows/main.yml/badge.svg)
# ТЗ проекта:

https://github.com/OtusGolang/final_project/blob/master/02-banners-rotation.md

Общее описание
Сервис "Ротация баннеров" предназначен для выбора наиболее эффективных (кликабельных) баннеров, в условиях меняющихся предпочтений пользователей и набора баннеров.

Предположим, что на сайте есть место для показа баннеров (слот) и есть набор баннеров, которые конкурируют за право показа в этом месте. Набор баннеров может меняться.

Сервис осуществляет "ротацию" баннеров, показывая те, которые наиболее вероятно приведут к переходу. Для этого используется алгоритм "Многорукий бандит": https://habr.com/ru/company/surfingbird/blog/168611/

# Архитектура
Сервис состоит из REST API и базы данных

# Команды

Сборка приложения
```
make build
```

Запуск приложения (включает сборку)
```
make run
```

Остановка приложения
```
make down
```

# Тесты

Запуск unit-тестов
```
make test
```

Запуск интеграционных тестов
```
make test-int
```

Запуск линтера
```
make lint
```