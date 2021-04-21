docker rm -f redisServer
docker run -d -p 6379:6379 --name redisServer redis