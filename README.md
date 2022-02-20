# 下载
直接整个clone到本地
```
git clone https://github.com/lichmaker/short-url.git
```

# 配置文件
```
cp .env.example .env
vim .env
```

.env
```
# 本地环境-local 生产环境-production
APP_ENV=local
# 后续生成，可暂时留空
APP_KEY=
# 是否开启debug，debug模式下部分限制会放宽，例如jwt过期时间，在生产环境中谨慎开启
APP_DEBUG=true
# 项目域名
APP_URL=http://localhost:8006
# 端口
APP_PORT=8006

# mysql 配置
DB_CONNECTION=mysql
DB_HOST=host.docker.internal
DB_PORT=3306
DB_DATABASE=shorturl
DB_USERNAME=root
DB_PASSWORD=lichmakerroot
DB_DEBUG=2

# 应用日志生成模式，single-所有日志都放在单个文件中。 daily-所有日志按照日期区分
LOG_TYPE=daily
# 应用日志记录类型
LOG_LEVEL=debug
# 应用日志的存储路径。daily模式下文件名必须为logs.log，否则无非分割日期
LOG_NAME=storage/logs/logs.log

# redis 配置
REDIS_HOST=host.docker.internal
REDIS_PORT=6379

# 短链接缓存数量，请根据运行环境酌情修改
SHORT_CACHE_MAX=1000
```

# 安装

## 最小安装
编译后创建一个最小容器进行运行， 若后续进行热更新，需要自行编译并执行 `./shorturl serve restart`
```
docker build -f ./Dockerfile -t shorturl

# 启动容器
docker run --name="shorturl" --rm -p 宿主端口:容器端口 -d shorturl
# 假设 APP_PORT=8006 , 需要暴露 8005 端口
docker run --name="shorturl" --rm -p 8005:8006 -d shorturl
```

## 编译环境安装

```
# 生成镜像
docker build -f ./Dockerfile-Dev -t shorturl .

# 启动容器
docker run --name="shorturl" -p 宿主端口:容器端口 -v 本地项目绝对路径:/sourcecode -d shorturl /bin/sh -c "while true; do  sleep 1; done"
# 举个栗子🌰 配置文件中 APP_PORT=8006 , 需要暴露 8005 端口
docker run --name="shorturl" -p 8005:8006 -v /Users/wuguozhang/go/src/github.com/lichmaker/short-url:/sourcecode -d shorturl /bin/sh -c "while true; do  sleep 1; done"

# 编译
docker exec -w /sourcecode shorturl go build
```

# 程序配置

## key 生成
```
docker exec -w /app shorturl ./shorturl key
```
将生成出来的key，填入 `.env` 中的 APP_KEY

## 数据表迁移
如果使用 “最小安装” ，需先自行进行迁移
```
go run main.go migrate up
```
如果使用 “编译环境安装”，可执行以下下命令
```
docker exec -w /app shorturl ./shorturl migrate up
```

# 程序启动
使用 “最小安装” ，直接启动容器即可，不再赘述。

使用 “编辑环境安装”，执行以下命令
```
docker exec -w /sourcecode shorturl ./shorturl serve -d
```

# 平滑重启（热更新）
使用 “最小安装”，只能自行编译，放到容器中，然后在容器中执行以下命令：
```
./shorturl serve restart
```

使用 “编辑环境安装”，执行以下命令
```
# 直接重新编译(非热更新可不执行)
docker exec -w /app shorturl go build
# 平滑重启
docker exec -w /app shorturl ./shorturl serve restart
```

# API
swagger todo

# kafka

## 容器
```
docker-compose -f kafka-docker-compose.yml up -d
```

## 创建 topic

```
docker exec c_kafka1 kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 2 --partitions 3 --topic shorturl_stat
```