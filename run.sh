#!/bin/bash

cd $(dirname ${0})
ROOT_PATH=$(pwd)
echo $ROOT_PATH

start(){
	cd ${ROOT_PATH}/server/src/bestsell
	pwd
	chmod 777 ./start.sh
	./start.sh start
}

stop(){
	cd ${ROOT_PATH}/server/src/bestsell
	pwd
	chmod 777 ./start.sh
	./start.sh stop
}

startr(){
	cd ${ROOT_PATH}/server/src/bestsell
	pwd
	chmod 777 ./start.sh
	./start.sh startr
}

stopr(){
	cd ${ROOT_PATH}/server/src/bestsell
	pwd
	chmod 777 ./start.sh
	./start.sh stopr
}

excel(){
	cd ${ROOT_PATH}/server/src/bestsell
	pwd
	chmod 777 ./start.sh
	./start.sh excel
}

case "${1}" in
	build)build;;
	start)start;;
	stop)stop;;
	startr)startr;;
	stopr)stopr;;
	excel)excel;;
	*)echo "Usage : ${0} start|stop, not ${1}"
esac
