docker rm -f mongoContainer

docker run -d \
-p 27017:27017 \
--name mongoContainer \
mongo