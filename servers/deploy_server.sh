# build and push sql
cd db && bash build_sql.sh && cd ..
docker push towm1204/mysql

# build and deploy message service
cd dashboards && bash deploy.sh && cd ..

# start deploying gateway
cd gateway && bash deploy.sh && cd ..

