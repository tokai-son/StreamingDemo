package service

import (
	"fmt"
	"log"
	"time"

	l "github.com/tokai-son/StreamingDemo/pkg/log"

	pb "github.com/tokai-son/StreamingDemo/api/generated/github.com/tokai-son/StreamingDemo"
)

type VideoStreamServer struct {
	pb.UnimplementedVideoStreamServiceServer
}

func (s *VideoStreamServer) StreamVideo(stream pb.VideoStreamService_StreamVideoServer) error {
	// logger
	logger := l.NewLogger()

	for {
		req, err := stream.Recv()

		if err != nil {
			if err.Error() == "EOF" {
				logger.Info("EOF")
				break
			}
			logger.Error(err)
			return err
		}

		logger.Info(req.Quality, req.StartTime, req.VideoID)

		// サンプルとして擬似的なチャンクデータを生成
		for i := 0; i < 10; i++ {
			chunk := &pb.StreamResponse{
				ChunkData:   []byte(fmt.Sprintf("Chunk %d of video %s", i, req.VideoID)),
				ChunkSize:   uint64(len([]byte(fmt.Sprintf("Chunk %d", i)))),
				Sequence:    uint64(i),
				EndOfStream: i == 9, // 最後のチャンクにフラグを立てる
			}

			// クライアントに送信
			if err := stream.Send(chunk); err != nil {
				log.Printf("Error sending stream: %v\n", err)
				return err
			}

			// 擬似的な遅延
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil
}
