# Forum Z Auth Server

## Download

```bash
git clone git@github.com:oyamo/forumz-auth-server.git
cd forumz-auth-server
```

## Build
```shell
docker build -t auth-server:1.0.0 .
```

## Run
```shell
docker run -d \
  --network sandbox \
  -e  AUTH_SERVICE_DATABASE_DSN='postgresql://postgres:5432/auth?user=dev&password=Test@12345' \
  -e  AUTH_SERVICE_P12_CERTIFICATE='./keystore.p12' \
  -e  AUTH_SERVICE_PUBLIC_KEY='./public_key.pub' \
  -e  AUTH_SERVICE_CERT_PASSWORD='Testing@123456' \
  -e  AUTH_SERVICE_KAFKA_CONSUMER='' \
  -e  AUTH_SERVICE_KAFKA_PRODUCER='' \
  -e  AUTH_SERVICE_REDIS_SERVER='postgres:6379' \
  auth-server:1.0.0

```