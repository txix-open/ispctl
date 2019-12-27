# ispctl

Консольная утилита для получения и переопределения конфигруации модулей isp

* [Требования](#Требования)
* [Установка](#Установка)
* [Конифгурация](#Конифгурация)
* [Использование](#Использование)
* [Описание](#Описание)
* [Примеры](#Пример)
    - [ispctl status](#ispctl-status)
    - [ispctl get](#ispctl-get)
    - [ispctl set](#ispctl-set)
    - [ispctl delete](#ispctl-delete)
    - [ispctl schema](#ispctl-schema)
    - [ispctl common](#ispctl-common)
        * [set](#ispctl-common-set)
        * [get](#ispctl-common-get)
        * [delete](#ispctl-common-delete)
        * [remove](#ispctl-common-remove)
        * [link](#ispctl-common-link)
        * [unlink](#ispctl-common-unlink)
        * [contain](#ispctl-common-contain)
    - [запрос с флагами](#Запрос-с-флагами)


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
gateHost: 127.0.0.1:9002
instanceUuid: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
* `gateHost` - адрес и порт GRPC любого isp-config-service в кластере
* `instanceUuid` - UUID  - идентификатор экземпляра приложений

## Использование
```bash
ispctl [flag...]  status    [local_flag]
ispctl [flag...]  get      module_name  property_path  [local_flag]
ispctl [flag...]  set      module_name  property_path  [new_object]
ispctl [flag...]  delete   module_name  property_path
ispctl [flag...]  schema   module_name  [local_flag]
ispctl [flag...]  common   sub_command
```

## Описание

| Команды     | Описание                                                                                    |
|-------------|---------------------------------------------------------------------------------------------|
| `status`    | возвращает доступные конфигурации модулей, их состояния и подключения                       |
| `get`       | возвращает объект конфигурации указанного модуля                                            |
| `set`       | изменяет объект конфигурации указанного модуля                                              |
| `delete`    | удаляет объект конфигурации указанного модуля                                               |
| `schema`    | возвращает схему конфигурации указанного модуля                                             |
| `common`    | комманда для взаимодействие с общими конфигурациями                                         |
|             |                                                                                             |

| Флаги       | Параметры | Описание                                                                                                            |
|-------------|-----------|---------------------------------------------------------------------------------------------------------------------|
| `-g`        | string    | переопределяет gateHost из конфигурации утилиты, значение должно быть экранировано с помощью `' '`                  |
| `-u`        | string    | переопределяет instanceUuid из конфигурации утилиты, значение должно быть экранировано с помощью `' '`              |
| `-c`        |           | раскрашивает json перед выводом на экран                                                                            |
| `-unsafe`   |           | отключает проверку схемы перед изменением конфигурации                                                              |
|             |           |                                                                                                                     |

| Аргументы       | Описание                                                                                                    |
|-----------------|-------------------------------------------------------------------------------------------------------------|
| `module_name`   | Название модуля с которым происходит взаимодействие                                                         |
| `proprtry_path` | Путь к объекту конфигурации, при значении `.` работа происходит со всей конфигурацией модуля                |
| `new_object`    | Новый объект, значение должно быть экранировано с помощью `' '`, при отсутсвии ожидается ввод из stdin      |
|                 |                                                                                                             |

| Локальные флаги | Параметры | Команды | Описание                                                                                                                        |
|-----------------|-----------|---------|---------------------------------------------------------------------------------------------------------------------------------|
| `-o`            | string    | schema  | определяет формат вывода схемы в stdout; по умолчанию `json`; возможные значения `json`, `html`                                 |
| `-full`         | bool      | get     | осуществляет взаимодействие с объектами конфигурации модуля с учетом объектом общих конфигураций, которые имеют связь с модулем |
|                 |           |         |                                                                                                                                 |

### sub_command
#### common
```bash
ispctl [flag...]  common   set      config_name      property_path   [new_object]
ispctl [flag...]  common   get      [config_name]    property_path
ispctl [flag...]  common   delete   config_name      property_path
ispctl [flag...]  common   remove   config_name
ispctl [flag...]  common   link     config_name      module_name
ispctl [flag...]  common   unlink   config_name      module_name
ispctl [flag...]  common   contain  config_name
```

| Аргументы       | Описание                                                                                                    |
|-----------------|-------------------------------------------------------------------------------------------------------------|
| `config_name`   | Названия конфига с которым происходит взаимодействие                                                        |
| `module_name`   | Название модуля с которым происходит взаимодействие                                                         |
| `proprtry_path` | Путь к объекту конфигурации, при значении `.` работа происходит со всей конфигурацией модуля                |
| `new_object`    | Новый объект, значение должно быть экранировано с помощью `' '`, при отсутсвии ожидается ввод из stdin      |
|                 |                                                                                                             |


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

### ispctl common
#### ispctl common get
Возвращает список названий общих конфигураций или, если в первом аргументе указано название общей конфигурации, возвращает объект общей конфигурации по указанному пути

Запрос
```bash
ispctl common get 
```
Ответ
```bash
test b c 
```
Запрос
```bash
ispctl common get test .
```
Ответ
```bash
{
    "test": null
}
```
#### ispctl common set
Добавляет объект в общую конфигурацию по указанному пути

Запрос
```bash
ispctl common set test .test '{"a":"a","b":"b"}'
```
Ответ
```bash
{
    "test": {
            "a": "a",
            "b": "b"
    }
}
```
#### ispctl common delete
Удаляет объект из общей конфигурации по указанному пути

Запрос
```bash
ispctl common delete test .test.a 
```
Ответ
```bash
{
    "test": {
            "b": "b"
    }
}
```
#### ispctl common link
Связывает 'common_config' к конфигурации модуля по 'module_name'. Возвращает список общих конфигураций, которые связаны с модулем

Запрос
```bash
ispctl common link test module_example
```
Ответ
```bash
[test] [second_common_config]
```
#### ispctl common remove
Удаляет объект общий конфигурации, если он не связан с конфигурациями модулей. Если общая конфигурация имеет связи, выводит список модулей с которыми установлена связь

Запрос
```bash
ispctl common remove test
```
Ответ
```bash
config [test] not deleted, need unlink in next modules:
[module_example]
```
#### ispctl common unlink
Отвязывает 'common_config' от конфигурации модуля по 'module_name'. Возвращает список общих конфигураций, которые связаны с модулем

Запрос
```bash
ispctl common unlink test module_example
```
Ответ
Ответ
```bash
[second_common_config]
```
#### ispctl common contain
Возвращает список названий модулей, с которыми имеет связь общая конфигурация

Запрос
```bash
isp-ctl common contain second_common_config
```
Ответ
Ответ
```bash
[module_example]
```

### Запрос с флагами
Запрос
```bash
ispctl -u '00000000-1111-2222-3333-444444444444' -g '127.0.0.1:0000' -c set example . '{"metrics":{"address":{"ip":"127.0.0.1","newField":"100","port":1},"gc":1,"memory":false}}'
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
