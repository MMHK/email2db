FROM golang:1.16-alpine as builder

# Add Maintainer Info
LABEL maintainer="Sam Zhou <sam@mixmedia.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go version \
 && export GO111MODULE=on \
 && export GOPROXY=https://goproxy.io \
 && go mod vendor \
 && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o email2db \
 && chmod +x email2db

######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /app

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64 \
 && chmod +x /usr/local/bin/dumb-init \
 && apk add --update libintl \
 && apk add --virtual build_deps gettext \
 && apk add --no-cache tzdata \
 && cp /usr/bin/envsubst /usr/local/bin/envsubst \
 && apk del build_deps \
 && echo "{}" > /app/config.json

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/email2db .
COPY ./web /app/web

ENV HTTP_LIST="0.0.0.0:8843" \
 S3_KEY="" \
 S3_SECRET="" \
 S3_SECRET="" \
 S3_REGION="ap-southeast-1" \
 S3_PREFIX="email2db" \
 S3_BUCKET="s3.test.mixmedia.com" \
 MYSQL_HOST="" \
 MYSQL_PORT=3306 \
 MYSQL_DATABASE="email2db" \
 MYSQL_USERNAME="" \
 MYSQL_PASSWORD="" \
 PARSER_TYPE="pop3" \
 ZOHO_POP3_HOST="" \
 ZOHO_POP3_PORT=995 \
 ZOHO_EMAIL="" \
 ZOHO_APP_SECRET="" \
 ZOHO_POP3_TLS=true \
 TZ="Asia/Hong_Kong" \
 CHECK_INTERVAL=10800 \
 FETCH_LIMIT=100 \
 LOG_LEVEL=INFO \
 WEB_ROOT=/app/web

EXPOSE 8843

ENTRYPOINT ["dumb-init", "--"]

CMD envsubst < /app/config.json > /app/temp.json \
 && /app/email2db -c /app/temp.json
