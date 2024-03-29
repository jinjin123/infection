package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"infection/util/lib"
	"os"
	"time"
)

func put() {
	ip := os.Args[1]
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{lib.MIDETCD},

		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = cli.Put(ctx, "/url/ip/", ip)
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/url/ip/")
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
