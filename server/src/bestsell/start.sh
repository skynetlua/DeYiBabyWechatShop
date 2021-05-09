#!/bin/bash

cd $(dirname ${0})
ROOT_PATH=$(pwd)
pwd

cd ${ROOT_PATH}/../../
pwd
PROJECT_PATH=$(pwd)
cd ${ROOT_PATH}

build(){
	echo "build ..."
	go build ${ROOT_PATH}/main.go
	chmod 777 ${ROOT_PATH}/main
}

# tool(){ 129.204.119.83 
# 	echo "tool ..."
# 	go run main.go -m tool -c ./config.json
# }

LOG_PATH=${PROJECT_PATH}/../log
mkdir -p ${LOG_PATH}

excel(){
	echo "excel ..."
	echo "${ROOT_PATH}/main.go"
	go run ${ROOT_PATH}/main.go -m excel -c ./config.json -p ${PROJECT_PATH}
}

start(){
	echo "api start ..."
	stop
	build
	TIME_DIR=$(date "+%Y_%m_%d")
	TIME_FILE=$(date "+%H_%M_%S")
	mkdir -p ${LOG_PATH}/${TIME_DIR}
	nohup ${ROOT_PATH}/main -c ./config.json >${LOG_PATH}/${TIME_DIR}/output${TIME_FILE}.log 2>&1 &
	PID=$(ps -ef |grep ./config.json |grep -v grep |awk '{print $2}')
	if [[ ! -z ${PID} ]]; then
		echo "api started, pid = ${PID}"
	else
		echo "api start failed!"
	fi
}

stop(){
	echo "api stop ..."
	PID=$(ps -ef |grep ./config.json |grep -v grep |awk '{print $2}')
	if [[ ! -z ${PID} ]]; then
		kill ${PID}
		echo "api stopped "
	else
		echo "api no run!"
	fi	
}

startr(){
	echo "res start ..."
	stopr
	build
	TIME_DIR=$(date "+%Y_%m_%d")
	TIME_FILE=$(date "+%H_%M_%S")
	mkdir -p ${LOG_PATH}/${TIME_DIR}
	nohup ${ROOT_PATH}/main -c ./configres.json >${LOG_PATH}/${TIME_DIR}/outputr${TIME_FILE}.log 2>&1 &
	PID=$(ps -ef |grep ./configres.json |grep -v grep |awk '{print $2}')
	if [[ ! -z ${PID} ]]; then
		echo "res started, pid = ${PID}"
	else
		echo "res start failed!"
	fi
}

stopr(){
	echo "res stop ..."
	PID=$(ps -ef |grep ./configres.json |grep -v grep |awk '{print $2}')
	if [[ ! -z ${PID} ]]; then
		kill ${PID}
		echo "res stopped "
	else
		echo "res no run!"
	fi	
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
