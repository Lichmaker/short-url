一个简单的短链接服务，支持使用kafka进行异步浏览统计。

开发中 feature : 
- 接入 sentry
- 加入布隆过滤器
- 尝试使用 map reduce

目前文档中介绍方式，使用 `Dockerfile-Dev` 创建整个golang环境，方便理解整个环境运作，且更好的支持热更新。
实际生产环境使用中，可以使用 `Dockerfile` ，创建编译后最小运行环境。

# 项目依赖
- mysql
- redis
- kafka (可选)

# 使用 Docker 安装

- 克隆项目
```
git clone https://github.com/lichmaker/short-url.git
```

- 更改项目配置 `.env`
```
cp .env.example .env
```

.env
```
# 可选项： product / local
APP_ENV=local
# 可暂时留空， 后续使用项目生成后再填入
APP_KEY=ZDkocQEicULJvQrlvojzIFlizwwxbCBF
# 日志显示 DEBUG 信息
APP_DEBUG=true
# 项目URL，涉及到短链接生成
APP_URL=http://localhost:8006
# 项目暴露端口
APP_PORT=8006

# mysql链接，请自行注意容器内网络
DB_CONNECTION=mysql
#DB_HOST=127.0.0.1
DB_HOST=host.docker.internal
DB_PORT=3306
# 库名可自定义
DB_DATABASE=shorturl
DB_USERNAME=root
DB_PASSWORD=password
DB_DEBUG=2

# 日志类型，可选 ： single - 单个文件， daily - 每日一个文件
LOG_TYPE=daily
# 日志等级
LOG_LEVEL=debug
# 自定义日志文件路径，不选则默认在项目 storage 目录中
#LOG_NAME=/sourcecode/testlogs/logs.log

# redis配置，必须
REDIS_HOST=host.docker.internal
REDIS_PORT=6379

# 短链接缓存数量上限
SHORT_CACHE_MAX=1000

# 启动统计，1-开启（需要kafka），0-关闭
STATISTIC_ENABLE=1

# kafka配置。STATISTIC_ENABLE=1时为必须
KAFKA_ADDRESS=host.docker.internal:29093,host.docker.internal:29094,host.docker.internal:29095
```

- 创建镜像
```
docker build -f ./Dockerfile-Dev -t shorturl
```

- 创建、启动容器
```
# 启动容器
docker run --name="shorturl" -p 宿主端口:容器端口 -v 本地项目绝对路径:/sourcecode -d shorturl /bin/sh -c "while true; do  sleep 1; done"
# 举个栗子🌰 配置文件中 APP_PORT=8006 , 需要暴露 8005 端口
docker run --name="shorturl" -p 8005:8006 -v /Users/wuguozhang/go/src/github.com/lichmaker/short-url:/sourcecode -d shorturl /bin/sh -c "while true; do  sleep 1; done"
```

- 环境配置
```
# 设置golang环境
docker exec -w /sourcecode shorturl go env -w GO111MODULE=on
docker exec -w /sourcecode shorturl go env -w GOPROXY=https://goproxy.cn,direct
docker exec -w /sourcecode shorturl go mod vendor
# 编译
docker exec -w /sourcecode shorturl go build
```

- 启动 web 服务
```
# 启动
docker exec -w /sourcecode shorturl ./shorturl serve -d

# 平滑重启
docker exec -w /sourcecode shorturl ./shorturl serve restart
```

- 启动统计服务
    - kafka创建topic。topic名称为 shorturl_stat ，具体参数可根据实际情况来
    ```
    kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 2 --partitions 3 --topic shorturl_stat
    ```
    - 启动消费者
    ```
    # 需输入启动实例数量，例如启动3个
    docker exec -w /sourcecode shorturl ./shorturl consumer -c 3
    ```
    - 平滑关闭所有消费者
    ```
    docker exec -w /sourcecode shorturl ./shorturl consumer shutdown
    ```