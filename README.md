## 一.项目配置
### 1. 安装依赖
```shell
go mod tidy
```
### 2. 配置./configs文件
```shell
echo '{"addr": ""}' > ./configs/server.json
echo '{"url": ""}' > ./configs/amqp.json
echo '{"dsn": ""}' > ./configs/db.json
echo '{"addr": "", "password": ""}' > ./configs/rdb.json
```
### 3. 配置私钥和公钥(用于登录或注册的敏感信息传输)  
```shell
openssl genrsa -out private.pem 2048
```
```shell
openssl rsa -in private.pem -pubout -out public.pem
```
### 4.配置docker
```shell
docker pull redis
docker pull rabbitmq:management
docker pull groonga/mroonga:latest
```
```shell
#redis和mysql需要提供密码
#redis的.conf文件需配置
#mysql需要运行db.sql文件中的命令
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
docker run --name redis -p 6379:6379 -d redis
docker run --name mysql -e MYSQL_ROOT_PASSWORD=rootpassword -p 3306:3306 -d groonga/mroonga:latest
```
### 5. 项目运行
```shell
go run ./cmd/main.go
```
## 二.项目介绍

