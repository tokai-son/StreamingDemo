package main

import (
	"log"
	"net"

	pb "github.com/tokai-son/StreamingDemo/api/generated/github.com/tokai-son/StreamingDemo"
	service "github.com/tokai-son/StreamingDemo/internal/service"

	"gocv.io/x/gocv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// ポート5432でリスニングする
	lis, err := net.Listen("tcp", ":5432")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// カメラを開く
	webcam, err := gocv.OpenVideoCapture(0) // 0はデフォルトカメラ
	if err != nil {
		log.Fatalf("Error opening video capture device: %v\n", err)
	}
	defer webcam.Close()

	// gRPCサーバを作成
	s := grpc.NewServer()

	// gRPCサーバにリフレクションを登録 (gRPCurlなどで利用するため)
	reflection.Register(s)

	// サービスを登録
	pb.RegisterVideoStreamServiceServer(s, &service.VideoStreamServer{
		VideoCapture: webcam,
	})

	log.Println("Server is running on port 5432")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
