package rpcChat

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/kafka"
	log2 "Open_IM/src/common/log"
	"Open_IM/src/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/src/proto/chat"
	"Open_IM/src/utils"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type rpcChat struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	producer        *kafka.Producer
}

func NewRpcChatServer(port int) *rpcChat {
	rc := rpcChat{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImOfflineMessageName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	rc.producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	return &rc
}

func (rpc *rpcChat) Run() {
	log2.Info("", "", "rpc get_token init...")

	address := utils.ServerIP + ":" + strconv.Itoa(rpc.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log2.Error("", "", "listen network failed, err = %s, address = %s", err.Error(), address)
		return
	}
	log2.Info("", "", "listen network success, address = %s", address)

	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()

	//service registers with etcd

	pbChat.RegisterChatServer(srv, rpc)
	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), utils.ServerIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log2.Error("", "", "register rpc get_token to etcd failed, err = %s", err.Error())
		return
	}

	err = srv.Serve(listener)
	if err != nil {
		log2.Info("", "", "rpc get_token fail, err = %s", err.Error())
		return
	}
	log2.Info("", "", "rpc get_token init success")
}
