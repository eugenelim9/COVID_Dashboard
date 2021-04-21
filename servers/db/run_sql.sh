bash build_sql.sh
export MYSQL_ROOT_PASSWORD="qazwsx"

docker rm -f mysqldemo

docker run -d \
-p 3306:3306 \
--name mysqldemo \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=demo \
towm1204/mysql