package repository

import (
	"github.com/txix-open/isp-kit/grpc/client"
)

func NewGrpcClientWithHost(host string) (*client.Client, error) {
	cli, err := client.Default()
	if err != nil {
		return nil, err
	}
	cli.Upgrade([]string{host})
	return cli, nil
}
