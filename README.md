# ispctl

## Описание
Консольная утилита для получения и переопределения конфигруации модулей isp

## Требования
* Linux

## Установка
```bash
yum install ispctl
```

## Конифгурация
Путь к файлу: `/etc/ispctl/config.yml`

Содержимое файла:
```yaml
gateHost: 127.0.0.1:0000
instanceUuid: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
* `gateHost` - адрес любого isp-convert-service в кластере
* `instanceUuid` - UUID  - идентификатор экземпляра приложений

## Использование
```bash
ispctl status  [flag...]
ispctl get     [flag...]  module_name  property_path
ispctl set     [flag...]  module_name  property_path  [new_object]
ispctl delete  [flag...]  module_name  property_path
```
## Описание

| Команды   | Описание                                                                                    |
|-----------|---------------------------------------------------------------------------------------------|
| `status`  | возвращает доступные конфигурации модулей, их состояния и подключения                       |
| `get`     | возвращает объект конфигурации указанного модуля                                                   |
| `set`     | изменяет объект конфигурации указанного модуля                                                     |
| `delete`  | удаляет объект конфигурации указанного модуля                                                      |
|           |                                                                                             |

| Флаги  | Параметры | Описание                                            |
|--------|-----------|-----------------------------------------------------|
| `-g`   | string    | переопределяет gateHost из конфигурации утилиты     |
| `-u`   | string    | переопределяет instanceUuid из конфигурации утилиты |
| `-c`   |           | раскрашивает json перед выводом на экран            |
|        |           |                                                     |

| Аргументы       | Описание                                                                                          |
|-----------------|---------------------------------------------------------------------------------------------------|
| `module_name`   | Модуль с которым происходит взаимодействие                                                        |
| `proprtry_path` | Путь к объекту конфигурации, при значении `.` работа происходит со всей конфигурацией модуля      |
| `new_object`    | Новый объект, должен быть экранирован с помощью `' '`, при отсутсвии ожидается ввод из stdin      |
|                 |                                                                                                   |


## Пример
### ispctl status
`ispctl status  [flag...]`

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
`ispctl get     [flag...]  module_name  property_path`
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
`ispctl set     [flag...]  module_name  property_path  [new_object]`

При указании `new_object` необходимо его экранировать. При его отсутсвии ожидается ввод из stdin

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
`ispctl delete  [flag...]  module_name  property_path`
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

### Запрос с флагами
Запрос
```bash
ispctl -u 00000000-1111-2222-3333-444444444444 -g 127.0.0.1:0000 -c set example . '{"metrics":{"address":{"ip":"127.0.0.1","newField":"100","port":1},"gc":1,"memory":false}}'
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