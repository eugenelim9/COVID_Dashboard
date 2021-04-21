# exporting env var
# comment if runs in local
export TLSCERT=/etc/letsencrypt/live/api.t-mokaramanee.me/fullchain.pem \
export TLSKEY=/etc/letsencrypt/live/api.t-mokaramanee.me/privkey.pem \
export MYSQL_ROOT_PASSWORD="qazwsx" \
export MY_SQL_HOST="mysqlContainer" \
export MY_SQL_DB="demo" \
export MY_SQL_PORT="3306" \
export DSN="root:qazwsx@tcp(mysqlContainer:3306)/demo" \
export SESSIONKEY="akey" \
export REDDISADDR="redisServer:6379" \
export ADDR=":443"
export DASHBOARDADDR="myDashboardServer:8080"
export DASHBOARDPORT="8080"
export MONGO_ENDPOINT="mongodb://customMongoContainer:27017/test"


# create private network for everything to be run on
docker network create network-441

# run redis
docker rm -f redisServer

docker run -d \
--network network-441 -p 6379:6379 \
--name redisServer redis

# running my sql
docker rm -f mysqlContainer
docker pull towm1204/mysql

docker run -d \
--name mysqlContainer -p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=demo \
--network network-441 \
towm1204/mysql

# running mongo
docker rm -f customMongoContainer

docker run -d \
--network network-441 \
-p 27017:27017 \
--name customMongoContainer \
mongo

# sleep to wait for sql to be up
sleep 20

docker rm -f myDashboardServer
docker pull towm1204/dashboardservice

docker run -d --name myDashboardServer \
-e MONGO_ENDPOINT=$MONGO_ENDPOINT \
-e DASHBOARDPORT=$DASHBOARDPORT \
--network network-441 towm1204/dashboardservice



# remove last one and pull image
docker rm -f gatewayServer
docker pull towm1204/mygateway

# docker run, mounting cert, open port, env variables
docker run -d \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
--name gatewayServer -p 443:443 \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e DSN=$DSN \
-e SESSIONKEY=$SESSIONKEY \
-e ADDR=$ADDR \
-e DASHBOARDADDR=$DASHBOARDADDR \
-e REDDISADDR=$REDDISADDR \
--network network-441 \
towm1204/mygateway