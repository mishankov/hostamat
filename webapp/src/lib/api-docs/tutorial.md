# Публичный API Hostamat: быстрый старт

Этот документ показывает рабочий путь от создания API-ключа до публикации zip-архива и подключения пользовательского домена.

Публичный API v1 доступен по умолчанию здесь:

```bash
export HOSTAMAT_API_BASE_URL="https://api.hostamat.ru/v1"
```

Во всех примерах ниже используется переменная `HOSTAMAT_API_BASE_URL`. Если вы работаете с другим окружением, задайте в ней URL нужного API v1.

## Создание API-ключа

API-ключи создаются в dashboard, а не через публичный API.

1. Войдите в Hostamat dashboard.
2. Откройте страницу управления API-ключами, например `/dashboard/api-keys`.
3. Нажмите создание нового ключа.
4. Укажите понятное имя, например `CI deploy`.
5. Выберите срок действия. Для бессрочного ключа выберите вариант `never expires`.
6. Выберите только нужные scopes:
   - `sites:read` для чтения списка и деталей сайтов.
   - `sites:write` для создания, переименования и удаления сайтов.
   - `sites:upload` для загрузки raw zip-архива в существующий сайт.
   - `domains:read` для чтения пользовательских доменов.
   - `domains:write` для создания, проверки, перемещения и удаления пользовательских доменов.
7. Создайте ключ и сразу сохраните полный secret в формате `hmat_<id>_<secret>`.

Полный API-ключ показывается только один раз. Позже dashboard показывает только не секретный идентификатор или префикс, чтобы можно было понять, какой ключ используется интеграцией.

Для примеров задайте ключ в переменную окружения:

```bash
export HOSTAMAT_API_KEY="hmat_00000000-0000-0000-0000-000000000000_your-secret"
```

## Bearer authentication

Каждый запрос к публичному API передает ключ в заголовке `Authorization`:

```bash
curl -sS "$HOSTAMAT_API_BASE_URL/me" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY"
```

`GET /me` не требует resource scope. Он нужен для проверки, что ключ валиден, не истек, не отозван и принадлежит активному пользователю.

Пример успешного ответа:

```json
{
  "user": {
    "id": "018f68e5-c654-7740-a1d2-eaa8b13f22d8",
    "username": "user"
  },
  "apiKey": {
    "id": "018f68ef-1ac0-7408-9b87-db22fd61d4a2",
    "name": "CI deploy",
    "keyPrefix": "hmat_018f68ef",
    "scopes": ["sites:read", "sites:upload"],
    "expiresAt": null
  }
}
```

Все ошибки публичного API возвращаются в единой форме:

```json
{
  "error": {
    "code": "insufficient_scope",
    "message": "API key does not have the required scope."
  }
}
```

Сообщение `message` предназначено для чтения человеком и остается на английском. Для логики клиента используйте `error.code`. Идентификатор запроса возвращается в заголовке ответа `Hostamat-Trace-Id`, включая ошибочные ответы.

## Работа с сайтами

Для чтения сайтов нужен `sites:read`. Для создания, переименования и удаления нужен `sites:write`. Эти scopes не включают друг друга.

### Создать сайт

```bash
curl -sS -X POST "$HOSTAMAT_API_BASE_URL/sites" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name":"Landing"}'
```

Пример ответа:

```json
{
  "site": {
    "id": "018f68f5-3296-7a0c-9b1f-57b59a07d918",
    "name": "Landing",
    "hostamatDomain": "bright-river.hostamat.ru",
    "url": "https://bright-river.hostamat.ru",
    "createdAt": "2026-04-28T09:30:00Z",
    "updatedAt": "2026-04-28T09:30:00Z"
  }
}
```

Сохраните `site.id` для следующих команд:

```bash
export HOSTAMAT_SITE_ID="018f68f5-3296-7a0c-9b1f-57b59a07d918"
```

### Получить список сайтов

```bash
curl -sS "$HOSTAMAT_API_BASE_URL/sites" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY"
```

В v1 список не поддерживает пагинацию, фильтрацию или сортировку.

### Получить сайт

```bash
curl -sS "$HOSTAMAT_API_BASE_URL/sites/$HOSTAMAT_SITE_ID" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY"
```

### Переименовать сайт

```bash
curl -sS -X PATCH "$HOSTAMAT_API_BASE_URL/sites/$HOSTAMAT_SITE_ID" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name":"Production landing"}'
```

В v1 через публичный API у сайта можно менять только `name`.

### Удалить сайт

```bash
curl -i -X DELETE "$HOSTAMAT_API_BASE_URL/sites/$HOSTAMAT_SITE_ID" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY"
```

Успешный ответ возвращает `204 No Content`.

## Загрузка raw zip-архива

Для загрузки нужен `sites:upload`. Scope `sites:write` не дает право загружать файлы.

Endpoint принимает только raw body с `Content-Type: application/zip`:

```bash
zip -r site.zip dist

curl -sS -X POST "$HOSTAMAT_API_BASE_URL/sites/$HOSTAMAT_SITE_ID/upload" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY" \
  -H "Content-Type: application/zip" \
  --data-binary "@site.zip"
```

Загрузка синхронная. Успешный запрос заменяет все существующие файлы сайта результатом распаковки архива. В v1 нет upload status или upload history endpoints.

Пример ответа:

```json
{
  "summary": {
    "uploadedCount": 12,
    "skippedCount": 1,
    "totalBytes": 245760
  },
  "files": [
    {
      "path": "/index.html",
      "size": 4096,
      "status": "uploaded"
    },
    {
      "path": "/.DS_Store",
      "size": 6148,
      "status": "skipped",
      "reason": "unsupported_file_type"
    }
  ]
}
```

Правила путей внутри архива, допустимые файлы, лимиты размера и формат uploaded/skipped результатов совпадают с текущей логикой загрузки файлов в dashboard.

## Пользовательские домены

Domains API управляет только пользовательскими доменами. Сгенерированный домен Hostamat находится в ответах Site API в поле `hostamatDomain`.

Для чтения нужен `domains:read`. Для создания, проверки, перемещения и удаления нужен `domains:write`. Эти scopes не включают друг друга.

### Добавить домен

```bash
curl -sS -X POST "$HOSTAMAT_API_BASE_URL/domains" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"hostname\":\"www.example.com\",\"siteId\":\"$HOSTAMAT_SITE_ID\"}"
```

Пример ответа:

```json
{
  "domain": {
    "id": "018f6901-8db7-73a7-a6c9-c675f0f6767e",
    "hostname": "www.example.com",
    "status": "pending",
    "site": {
      "id": "018f68f5-3296-7a0c-9b1f-57b59a07d918",
      "name": "Production landing",
      "hostamatDomain": "bright-river.hostamat.ru"
    },
    "dnsInstructions": [
      {
        "type": "TXT",
        "name": "_hostamat.www.example.com",
        "value": "domain-verification-token"
      },
      {
        "type": "CNAME",
        "name": "www.example.com",
        "value": "dns.hostamat.ru"
      }
    ],
    "lastVerificationError": null,
    "lastCheckedAt": null,
    "activatedAt": null,
    "createdAt": "2026-04-28T09:35:00Z",
    "updatedAt": "2026-04-28T09:35:00Z"
  }
}
```

Сохраните `domain.id`:

```bash
export HOSTAMAT_DOMAIN_ID="018f6901-8db7-73a7-a6c9-c675f0f6767e"
```

После создания настройте DNS по `dnsInstructions`.

### Получить список доменов

```bash
curl -sS "$HOSTAMAT_API_BASE_URL/domains" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY"
```

В v1 список доменов глобальный для пользователя. Per-site domain routes нет.

### Получить домен

```bash
curl -sS "$HOSTAMAT_API_BASE_URL/domains/$HOSTAMAT_DOMAIN_ID" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY"
```

### Проверить DNS

```bash
curl -sS -X POST "$HOSTAMAT_API_BASE_URL/domains/$HOSTAMAT_DOMAIN_ID/verify" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY"
```

Если проверка не прошла, API вернет домен со статусом `verification_failed`. Актуальные инструкции по DNS остаются в `dnsInstructions`, а последняя ошибка проверки в `lastVerificationError`.

### Переместить домен на другой сайт

```bash
export HOSTAMAT_NEW_SITE_ID="018f6909-22bd-7006-9b2f-a88619c0c9d1"

curl -sS -X PATCH "$HOSTAMAT_API_BASE_URL/domains/$HOSTAMAT_DOMAIN_ID" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"siteId\":\"$HOSTAMAT_NEW_SITE_ID\"}"
```

Новый сайт должен принадлежать тому же пользователю, которому принадлежит API-ключ.

### Удалить домен

```bash
curl -i -X DELETE "$HOSTAMAT_API_BASE_URL/domains/$HOSTAMAT_DOMAIN_ID" \
  -H "Authorization: Bearer $HOSTAMAT_API_KEY"
```

Успешный ответ возвращает `204 No Content`.

## Частые ошибки

`missing_bearer_token`: заголовок `Authorization` отсутствует.

`malformed_bearer_token`: bearer token имеет неправильный формат.

`expired_api_key`: срок действия ключа истек.

`revoked_api_key`: ключ был удален или отозван в dashboard.

`unknown_api_key`: ключ не найден или secret не прошел проверку.

`api_key_owner_inactive`: пользователь-владелец ключа неактивен.

`insufficient_scope`: ключ валиден, но у него нет нужного scope.

`invalid_body`: JSON body не удалось прочитать.

`invalid_upload`: zip body не удалось прочитать.

`not_found`: ресурс не найден или не принадлежит владельцу ключа.

`forbidden`: доступ запрещен.

`site_limit_reached`: пользователь достиг лимита сайтов.

`unsupported_media_type`: upload endpoint получил не `application/zip`.

`upload_too_large`: zip-архив превышает допустимый размер.

`invalid_file_type`: файл из архива не поддерживается текущей логикой загрузки.

`hostname_taken`: домен уже подключен к другому сайту.

`name_too_long`: имя слишком длинное.

`missing_required_field`: обязательное поле не передано.

`invalid_hostname`: домен не прошел проверку формата или типа.

`internal_error`: внутренняя ошибка сервера.
