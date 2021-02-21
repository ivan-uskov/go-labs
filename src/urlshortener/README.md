# urlshortener

Сервис коротких ссылок, запускает веб-сервис,
который предоставляет возможность использовать короткие ссылки

```
>urlshortener -f config.json
```

Пример конфигурационного файлв
```json
{
  "not_found_message": "Page %s is not found",
  "server_port": 8080,
  "paths": {
    "/go-hd": "http://www.harley-davidson.com",
    "/go-code": "https://github.com/ivan-uskov/go-labs"
  }
}
```