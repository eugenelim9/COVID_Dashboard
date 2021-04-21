# build
bash build.sh

# push image to dockerhub
docker push towm1204/mygateway

# ssh and run script
ssh ec2-user@api.t-mokaramanee.me < vm-script.sh
