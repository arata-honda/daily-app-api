# daily-app-api
日記アプリ用API

# 環境

- docker
- go1.14 darwin/amd64  

# 開発環境起動方法

```bash
$ git clone https://github.com/arata-honda/daily-app-api.git
$ cd arata-honda/daily-app-api
$ docker-compose up -d
$ sh init-mysql.sh
```

起動後 http://localhost:8080 で 返ってくればよい。
