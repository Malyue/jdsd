#!/bin/bash

# 停止工作容器
echo "stop jdsd container"
docker stop jdsd
echo "jdsd stoped"

echo "delete jdsd container"
docker rm -f jdsd
echo "jdsd deleted"

docker rmi $(docker images | grep "none" |awk '{print $3}')
echo "image deleted"

docker build -t jdsd:v1 .
docker run -d --name=jdsd jdsd:v1
echo "build success"


## 查看镜像是否存在
#imageName=jdsd
#name=jdsd:v1
#
## 查询得到指定名称的容器id
#ARG1=$(docker ps -aqf "name=${imageName}")
#
## 查询得到指定名称的镜像id
#ARG2=$(docker images -q --filter reference=${name})
#
#docker stop jdsd
#docker rm jdsd
#
### 如果查询结果不为空，先停止容器再删除
##if [  -n "$ARG1" ]; then
##  echo $ARG1
##  docker stop ${imageName}
##  docker rm  ${imageName}
### docker rm -f $(docker stop $ARG1)
## echo "$name容器停止删除成功.....！！！"
##fi
#
#docker rmi $(docker images | grep "none" |awk '{print $3}')
#
##如果查询结果不为空，先删除镜像
##删除镜像
#if [  -n "$ARG2" ]; then
# docker rmi -f $ARG2
# echo "$name镜像删除成功.....！！！"
#fi
#
## build项目
#docker build -t jdsd:v1 .
#docker run -d --name=${imageName} ${name}
#echo "build success"