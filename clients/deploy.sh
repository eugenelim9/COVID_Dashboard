# call build
bash build.sh

# push image to dockerhub
docker push towm1204/dash-client

# ssh and run vm-script
ssh ec2-user@t-mokaramanee.me < vm-script.sh