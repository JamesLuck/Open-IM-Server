#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg

#Check if the service exists
#If it is exists,kill this process
check=$(ps aux | grep -w ./${msg_transfer_name} | grep -v grep | wc -l)
if [ $check -eq 1 ]; then
  oldPid=$(ps aux | grep -w ./${msg_transfer_name} | grep -v grep | awk '{print $2}')
  kill -9 $oldPid
fi
#Waiting port recycling
sleep 1
cd ${msg_transfer_binary_root}
nohup ./${msg_transfer_name} >>../logs/${msg_transfer_name}.log 2>&1 &
#Check launched service process
check=$(ps aux | grep -w ./${msg_transfer_name} | grep -v grep | wc -l)
if [ $check -eq 1 ]; then
  newPid=$(ps aux | grep -w ./${msg_transfer_name} | grep -v grep | awk '{print $2}')
  ports=$(netstat -netulp | grep ${newPid} | awk '{print $4}' | awk -F '[:]' '{print $NF}')
  allPorts=""

  for i in $ports; do
    allPorts=${allPorts}"$i "
  done
  echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS !!!"${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${msg_transfer_name}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${newPid}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
  echo -e ${YELLOW_PREFIX}${msg_transfer_name}${COLOR_SUFFIX}${RED_PREFIX}"SERVICE START ERROR !!! PLEASE CHECK ERROR LOG"${COLOR_SUFFIX}
fi
