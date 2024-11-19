package rpc

import (
	"fmt"
	"images-api/server"
	"net"
	"net/rpc"

	"cloud.google.com/go/storage"
	"github.com/kastoras/go-utilities"
)

type RPCServer struct {
	Cache         *server.Cache
	StorageClient *storage.Client
}

func NewServer() {

	rpcServer := &RPCServer{}

	err := rpcServer.Init()
	if err != nil {
		panic(err)
	}

	err = rpc.Register(rpcServer)
	if err != nil {
		panic(err)
	}

	rppPort, err := utilities.GetEnvParam("RPC_PORT", "5050")
	if err != nil {
		fmt.Printf("RPC Port could not be found")
		panic(err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", rppPort))
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	server.CreateCacheClient()
	defer server.CloseCacheClient()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go rpc.ServeConn(rpcConn)
	}

}

func (r *RPCServer) Init() error {
	client, err := server.InitStorageClient()
	if err != nil {
		return err
	}

	r.StorageClient = client

	return nil
}
