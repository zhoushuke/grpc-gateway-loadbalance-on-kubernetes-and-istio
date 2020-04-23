package main

//client.go

import (
	"fmt"
	"log"
	"math/rand"
	"time"
        "net"

	// "github.com/sercand/kuberesolver"
	// "google.golang.org/grpc/resolver"
	"google.golang.org/grpc/balancer/roundrobin"
	pb "gitlab.bj.sensetime.com/SenseGo/grpc-gateway-demo/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var conns []*grpc.ClientConn

func resolverconn() *grpc.ClientConn {
        // client, err := kuberesolver.NewInClusterK8sClient()
	// if err != nil {
	//	fmt.Println("errrr", err)
	//	return nil
	//}
	// resolver.Register(kuberesolver.NewBuilder(client, "kubernetes"))

	//fmt.Println("=============")
	// cc, err := grpc.Dial("kubernetes:///progress-service-svc:50051", grpc.WithInsecure())
	cc, err := grpc.Dial("progress-service-svc:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v \n", err)
		return nil
	}
	return cc
}

func pool() *grpc.ClientConn {
	if len(conns) == 0 {
		for i := 0; i < 3; i++ {
			conn, err := grpc.Dial("progress-service-svc:50051", grpc.WithInsecure())
			if err != nil {
				fmt.Printf("did not connect: %v \n", err)
				return nil
			}
			conns = append(conns, conn)
		}
		return conns[0]
	}
	rand.Seed(time.Now().UnixNano())
	return conns[rand.Intn(3)]
}

func lbconn() *grpc.ClientConn {
	//conn, err := grpc.Dial("progress-service-svc:50051",
	conn, err := grpc.Dial("dns:///progress-service-svc.just-4-test.svc.cluster.local:50051",
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
		//grpc.WithDisableServiceConfig(),
		//grpc.WithDefaultServiceConfig(`{"LB":"round_robin"}`),
	)
	if err != nil {
		fmt.Printf("did not connect: %v \n", err)
		return nil
	}
	return conn
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP
}


func main() {
	ip_v4 := GetOutboundIP().String()
	fmt.Println("begin client, ipv4: ", ip_v4)
	conn := resolverconn()
	c := pb.NewHelloClient(conn)
	for i := 0; i < 1000000; i++ {
		call(c, i, ip_v4)
		time.Sleep(time.Second * 5)
	}

}

func call(c pb.HelloClient, index int, ip_v4 string) {
	r, err := c.Hi(context.Background(), &pb.HiRequest{ToWho: fmt.Sprintf("index: %v", index), Message: "helloWorld", Clientip: ip_v4})
	if err != nil {
		fmt.Printf("could not SendMessage: %v", err)
		return
	}
	log.Printf("%s, %s", r.Message, r.Serverip)
}
