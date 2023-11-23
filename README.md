# email2db
現在有兩個服務需要使用 mailbox 接收email 並分析保存。

這兩項服務都需要佔用 我們一個獨立emailbox (現在的 mailbox 都價值不菲)，
所以需要一個 Email 2 DB 的Service 將一些 mailbox的郵件保存在 DB，方便後續的處理。
現在大部分的 email 發送服務都會附帶一下 上行email 的 webhook，例如 [SendGrid的這個](https://docs.sendgrid.com/for-developers/parsing-email/inbound-email)

所以基於這個webhook，我們可以將一些 mailbox 的電郵轉存至 DB

### Features

- 支持接收處理 sendgrid 的inbound mail webhook 請求
- 支持將 Email 數據(包括附件) 保存到 MySQL
- 支持將 附件 upload S3


### 配置

```json
{
  "listen": "",
  "web_root": "",
  "storages": {
    "s3": {
      "access_key": "",
      "secret_key": "",
      "bucket": "",
      "region": "",
      "prefix": ""
    }
  },
  "db": {
    "mysql": {
      "dsn": ""
    }
  },
  "tmp_path": ""
}
```

### Environment

```ini
S3_KEY=
S3_SECRET=
S3_REGION=ap-southeast-1
S3_BUCKET=

MYSQL_HOST=
MYSQL_PORT=3306
MYSQL_DATABASE=email2db
MYSQL_USERNAME=
MYSQL_PASSWORD=

TZ=Asia/Hong_Kong

HTTP_LIST=0.0.0.0:8843
WEB_ROOT=/var/www/email2db
```


#### Docker

此项目已经打包成docker 镜像

- 签出docker 镜像

```shell
docker pull mmhk/email2db:latest
```

- 运行
```shell
docker run --name email2db -p 8843:8843 mmhk/email2db:latest
```