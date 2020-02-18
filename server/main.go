package main

import (
	"fmt"
	"log"
	"net"

	gg "gitlab.bj.sensetime.com/SenseGo/grpc-gateway-demo/grpc_gateway"
	pb "gitlab.bj.sensetime.com/SenseGo/grpc-gateway-demo/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)
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

type grpcServer struct{}

func (s *grpcServer) Hi(ctx context.Context, in *pb.HiRequest) (*pb.HiResponse, error) {
	fmt.Println("service got request!", in)
        serverIP := GetOutboundIP().String()
	return &pb.HiResponse{FromWho: in.ToWho, Message: "message: " + in.Message, Serverip: " from serverIP: " + serverIP}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterHelloServer(s, &grpcServer{})

	go func() {
		cf := &gg.GatewayConfig{
			Fn:       pb.RegisterHelloHandlerFromEndpoint,
			GrpcPort: ":50051",
		}
		if err := gg.RegisterGRPCGateway(cf); err != nil {
			fmt.Println("grpc gateway closed:", err)
		}
	}()

	fmt.Println("start...")

	s.Serve(lis)
}
