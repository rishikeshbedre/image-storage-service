#!/bin/bash

#---------------------------------------Build Image Storage Service------------------------------------------------------------

#rm -rf *.tar

docker build --build-arg http_proxy="$http_proxy" --build-arg https_proxy="$https_proxy" -t image-storage-service:0.0.1 .

#docker run -it -p 3333:3333 -e HOSTIP=`hostname -I | awk '{print $1}'` image-storage-service:0.0.1   ------  to run with swagger

#docker run -it -p 3333:3333 -e HOSTIP=`hostname -I | awk '{print $1}'` -e PRODMODE=true image-storage-service:0.0.1  ----- to disable swagger

#docker save -o image-storage-service.tar image-storage-service:0.0.1

#chmod 777 *.tar

#---------------------------------------Pull Eclipse Mosquito-------------------------------------------------------------------

docker pull eclipse-mosquitto:1.6.9

#---------------------------------------Build Notification Subscriber (helper service)------------------------------------------

cd extras/notifierSubscriber

docker build --build-arg http_proxy="$http_proxy" --build-arg https_proxy="$https_proxy" -t notification-subscriber-service:0.0.1 .

cd ../..