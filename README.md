# URL Shortener

RESTFull go application for url shortening

## Description

This URL Shortener is implemented as a technical assignment for Ozon Bank reqruitment. There are several supported storages in application. The storage can be changed via argument <STORAGE>. All the repositories in project support concurrent updates and avoid link duplication with minimal lock overheads.

## Built with

* [go v1.23.6](https://go.dev/)
* [gin v1.10.0](https://gin-gonic.com/)

## OpenApi

Can be found in `docs` and accessed via `GET /swagger/index.html/`

## 🔄 Requests

### Generate short link

#### Method and path

`POST /api/v1/link/`

#### Request body

```json
{
  "url": "https://given.url.com/topic/2?a=213"
}
```

#### Response body

```json
{
  "status": "OK",
  "link": "https://shrt.com/<SHORT_LINK>"
}
```

### Get initial URL

#### Method and path

`GET /api/v1/link/<SHORT_LINK>`

#### Response body

```json

{
  "status": "OK",
  "link": "https://given.url.com/topic/2?a=213"
}
```

## 🙌 How to Start 

1. Clone and open this repo

   ```shell
   git clone https://github.com/zaqbez39me/UrlShortener.git
   cd UrlShortener
   ```

2. Copy `.env.example` in `.env`

   ```shell
   cp .env.example .env
   ```

3. Edit or delete environment variables in `.env`

4. Build and execute using instructions below

## ⚙️ Execution flags

* `--storage-type=<STORAGE_TYPE>`
  Possible values: (postgres, memory)
* `--cache-type=<CACHE_TYPE>`
  Possible values: (redis, none)

## 🛠️ How to build

```shell
make build
```

## ⚡ How to run

### Using Makefile
```shell
make execute STORAGE_TYPE=<STORAGE_TYPE>
```

### Build and run with docker

```shell
docker build --tag 'image_name' .
docker run -p 8080:8080 --storage-type <STORAGE_TYPE> --env-file .env 'image_name'
```

### Run using docker-compose with all deps

```shell
STORAGE_TYPE=<STORAGE_TYPE> docker compose up
```

### Build from sources

```shell
docker build -t 'image_name' .
```

## 🔎 QA

### Run Formatter

```shell
make fmt
```

### Run Tests

```shell
make test
```
