# ispctl

Консольная утилита для получения и переопределения конфигруации модулей isp

* [Требования](#Требования)
* [Установка](#Установка)
* [Конфигурация](#Конфигурация)
* [Использование](#Использование)
* [Описание](#Описание)
* [Примеры](#Пример)
    - [ispctl status](#ispctl-status)
    - [ispctl get](#ispctl-get)
    - [ispctl set](#ispctl-set)
    - [ispctl delete](#ispctl-delete)
    - [ispctl schema](#ispctl-schema)
    - [ispctl merge](#ispctl-merge)
    - [ispctl gitget](#ispctl-gitget)


## Требования
* Linux

## Установка
```bash
yum install ispctl
```

## Использование
```bash
ispctl [flag...]  status    [local_flag]
ispctl [flag...]  get      module_name  property_path  [local_flag]
ispctl [flag...]  set      module_name  property_path  [new_object]
ispctl [flag...]  delete   module_name  property_path
ispctl [flag...]  schema   module_name  [local_flag]
```

## Описание

| Команды  | Описание                                                                    |
|----------|-----------------------------------------------------------------------------|
| `status` | возвращает доступные конфигурации модулей, их состояния и подключения       |
| `get`    | возвращает объект конфигурации указанного модуля                            |
| `set`    | изменяет объект конфигурации указанного модуля                              |
| `delete` | удаляет объект конфигурации указанного модуля                               |
| `merge`  | поблочно производит слияние конфигурации модуля с json конфигом через stdin |
| `gitget` | скачивает из репозитория указанный файлы с указанного комита                |
| `schema` | возвращает схему конфигурации указанного модуля                             |
|          |                                                                             |

| Флаги       | Параметры | Описание                                               |
|-------------|-----------|--------------------------------------------------------|
| `-g`        | string    | адрес isp-config-service                               |
| `-unsafe`   |           | отключает проверку схемы перед изменением конфигурации |
|             |           |                                                        |

| Аргументы       | Описание                                                                                                    |
|-----------------|-------------------------------------------------------------------------------------------------------------|
| `module_name`   | Название модуля с которым происходит взаимодействие                                                         |
| `property_path` | Путь к объекту конфигурации, при значении `.` работа происходит со всей конфигурацией модуля                |
| `new_object`    | Новый объект, значение должно быть экранировано с помощью `' '`, при отсутсвии ожидается ввод из stdin      |
|                 |                                                                                                             |

| Локальные флаги | Параметры | Команды | Описание                                                                                                                        |
|-----------------|-----------|---------|---------------------------------------------------------------------------------------------------------------------------------|
| `-o`            | string    | schema  | определяет формат вывода схемы в stdout; по умолчанию `json`; возможные значения `json`, `html`                                 |
| `-full`         | bool      | get     | осуществляет взаимодействие с объектами конфигурации модуля с учетом объектом общих конфигураций, которые имеют связь с модулем |
|                 |           |         |                                                                                                                                 |

## Пример
### ispctl status
`ispctl [flag...]  status [local_flag]`

`[local_flag]` - локальный флаг `-o` определяет формат вывода доступных конфигураций и состояний модулей. Доступные значения: `json`.
По умолчанию вывод осуществляется в виде таблицы.

Запрос
```bash
ispctl status
```
Ответ
```bash
       MODULE NAME       |    STATUS     |        ADDRESSES
+------------------------+---------------+-------------------------+
  admin                  | CONNECTED     | 127.0.0.1
  auth                   | NOT_CONNECTED |
  config                 | CONNECTED     | 127.0.0.1
  converter              | CONNECTED     | 127.0.0.1
  journal                | NOT_CONNECTED |
```

### ispctl get
`ispctl [flag...]  get      module_name  property_path`
* Получение полной конфигурации

Запрос
```bash
ispctl get example .
```
Ответ
```bash
{
    "journal": {
        "bufferSize": 4092,
        "compress": true,
        "enable": false,
        "enableRemoteTransfer": true,
        "filename": "/var/log/example-service/runtime.log",
        "maxSizeMb": 512,
        "rotateTimeoutMs": 86400000
    },
    "metrics": {
        "address": {
            "ip": "127.0.0.1",
            "path": "/metrics",
            "port": 1
        },
        "gc": 2,
        "memory": true
    }
}
```

* Получение конкретного объекта конфигурации

Запрос
```bash
ispctl get example .metrics.address
```
Ответ
```bash
{
    "ip": "127.0.0.1",
    "path": "/metrics",
    "port": 1
}
```

### ispctl set
`ispctl [flag...]  set      module_name  property_path  [new_object]`

При указании `new_object` необходимо его экранировать. При его отсутсвии ожидается ввод из `stdin` до `EOF`

* Вставка нового поля в объект конфигурации

Запрос
```bash
ispctl set example .metrics.newField '"1000"'
```
Ответ
```bash
{
    "journal": {
        "bufferSize": 4092,
        "compress": true,
        "enable": false,
        "enableRemoteTransfer": true,
        "filename": "/var/log/example-service/runtime.log",
        "maxSizeMb": 512,
        "rotateTimeoutMs": 86400000
    },
    "metrics": {
        "address": {
            "ip": "127.0.0.1",
            "path": "/metrics",
            "port": 1
        },
        "gc": 2,
        "memory": true,
        "newField": "1000"
    }
}
```
* Изменение объекта конфигурации

Запрос
```bash
ispctl set example .metrics '{"address":{"ip":"198.0.0.1","newField":"100","port":1},"gc":1,"memory":false}'
```
Ответ
```bash
{
    "journal": {
        "bufferSize": 4092,
        "compress": true,
        "enable": false,
        "enableRemoteTransfer": true,
        "filename": "/var/log/example-service/runtime.log",
        "maxSizeMb": 512,
        "rotateTimeoutMs": 86400000
    },
    "metrics": {
        "address": {
            "ip": "198.0.0.1",
            "newField": "100",
            "port": 1
        },
        "gc": 1,
        "memory": false
    }
}
```
* Полное обновление конфигурации

Запрос
```bash
ispctl set example . '{"journal":{"bufferSize":1111,"compress":false,"enable":true,"filename":"/var/log/example-service/runtime.log","newField":"1000"},"metrics":{"address":{"ip":"198.0.0.1","newField":"100","port":1},"gc":1,"memory":false}}'
```
Ответ
```bash
{
    "journal": {
        "bufferSize": 1111,
        "compress": false,
        "enable": true,
        "filename": "/var/log/example-service/runtime.log",
        "newField": "1000"
    },
    "metrics": {
        "address": {
            "ip": "198.0.0.1",
            "newField": "100",
            "port": 1
        },
        "gc": 1,
        "memory": false
    }
}
```
### ispctl delete
`ispctl [flag...]  delete   module_name  property_path`
* Удаление объекта конфигурации

Запрос
```bash
ispctl delete example .journal
```
Ответ
```bash
{
    "metrics": {
        "address": {
            "ip": "198.0.0.1",
            "newField": "100",
            "port": 1
        },
        "gc": 1,
        "memory": false
    }
}
```
* Удаление поля из конфигурации

Запрос
```bash
ispctl delete example .metrics.address.ip
```
Ответ
```bash
{
    "metrics": {
        "address": {
            "newField": "100",
            "port": 1
        },
        "gc": 1,
        "memory": false
    }
}
```
* Удаление конфигурации

Запрос
```bash
ispctl delete example .
```
Ответ
```bash
{}
```

### ispctl schema
`ispctl [flag...]  schema   module_name  [local_flag]`
* Получение схемы конфигурации

Запрос
```bash
ispctl schema example
```
Ответ
```bash
{
    "title": "example"
    "schema": {
        "title": "Remote config",
        "type": "object",
        "required": [
            "journal",
        ],
         ...
         ...
         ...
    }
}
```
Запрос
```bash
ispctl schema example -o html
```
Ответ
```bash
<html>
<head>
    ...
    ...
    ...
<body>
<div class="results"></div>
<script>
    var schema = [
            {"schema":{"required":["journal"]...
    ...
    ...
</script>
</body>
</html>

```

### ispctl merge
Содержание файла `example.json`
```json
{
  "metrics": {
    "gc": 1,
    "memory": false
  },
  "property": "replaced"
}
```
Конфигурация модуля `example`
```json
{
  "database": {
    "host": "127.0.0.1"
  },
  "property": "original"
}
```

```bash
ispctl merge example < example.json
```
Результат в конфигурации модуля
```json
{
   "metrics": {
     "gc": 1,
     "memory": false
   },
   "property": "replaced",
   "database": {
      "host": "127.0.0.1"
   }
}
```

### ispctl gitget
```bash
ispctl gitget git@github.com:integration-system/ispctl.git main.go d4d7c679aad47ae204c2b3a4587b032a723fe315
```

### Запрос с флагами
Запрос
```bash
ispctl -g '127.0.0.1:9002' set example . '{"metrics":{"address":{"ip":"127.0.0.1","newField":"100","port":1},"gc":1,"memory":false}}'
```
Ответ
```json
{
    "metrics": {
        "address": {
            "ip": "127.0.0.1",
            "newField": "100",
            "port": 1
        },
        "gc": 1,
        "memory": false
    }
}
```
