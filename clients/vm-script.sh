# remove running image and pull image
docker rm -f client
docker pull towm1204/dash-client

# docker run, ports 80, 443, mount cert
docker run -d \
--name client \
-p 80:80 \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
towm1204/dash-client