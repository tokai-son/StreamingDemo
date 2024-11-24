package main

import (
	"bytes"
	"context"
	"image/jpeg"
	"log"
	"runtime"

	pb "github.com/tokai-son/StreamingDemo/api/generated/github.com/tokai-son/StreamingDemo"
	"gocv.io/x/gocv"
	"google.golang.org/grpc"
)

func main() {
	// メインスレッドにロック
	runtime.LockOSThread()

	// gRPCクライアントのセットアップ
	conn, err := grpc.Dial("localhost:5432", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewVideoStreamServiceClient(conn)
	stream, err := client.StreamVideo(context.Background(), &grpc.EmptyCallOption{})
	if err != nil {
		log.Fatalf("Failed to start stream: %v", err)
	}

	// ウィンドウを作成
	window := gocv.NewWindow("Video Stream")
	defer window.Close()

	// ストリームにリクエストを送信
	stream.Send(&pb.StreamRequest{
		Quality:   pb.VideoQuality_HIGH,
		StartTime: 0,
		VideoID:   "video_id_test",
	})

	// ストリームからデータを受信し続ける
	for {
		chunk, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive chunk: %v", err)
		}

		// JPEGデータをデコード
		img, err := jpeg.Decode(bytes.NewReader(chunk.ChunkData))
		if err != nil {
			log.Printf("Failed to decode image: %v", err)
			continue
		}

		// gocv.Matに変換
		mat, err := gocv.ImageToMatRGB(img)
		if err != nil {
			log.Printf("Failed to convert image to Mat: %v", err)
			continue
		}

		// 画像をウィンドウに表示
		window.IMShow(mat)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
