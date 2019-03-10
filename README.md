# webhook-service-updater

Служит для обновления сервисов в docker stack после сборки билдов
Для работы нужно просто дернуть внешнюю страничку и получить результат

Path for yml files: /app/files/yml/filename1.yml

Config path(in container): /app/files/config.json

Default headers:

WSU_USER: username

WSU_TOKEN: mySuperToken

WSU_PREFIX: ""(empty)

You can change your WSU_USER or WSU_TOKEN or WSU_PREFIX from env

Example config:
```json
{
  "users": {
    "username": {
      "tokens": [
        "mySuperToken",
        "mySuperToken2"
      ],
      "services": [
        "service1",
        "service2",
        "service3"
      ],
      "stacks": [
        {
          "service": "service1",
          "filename": "filename1"
        }
      ],
      "admin": true,
      "lock_ip": true,
      "ips": [
        "127.0.0.1"
      ]
    }
  }
}
```