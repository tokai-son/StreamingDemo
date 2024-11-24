package service

import (
	"fmt"
	"log"
	"time"

	l "github.com/tokai-son/StreamingDemo/pkg/log"
	"gocv.io/x/gocv"

	pb "github.com/tokai-son/StreamingDemo/api/generated/github.com/tokai-son/StreamingDemo"
)

type VideoStreamServer struct {
	pb.UnimplementedVideoStreamServiceServer
	VideoCapture *gocv.VideoCapture
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

		// カメラからフレームを取得
		img := gocv.NewMat()
		defer img.Clone()

		// サンプルとして擬似的なチャンクデータを生成
		for {
			if ok := s.VideoCapture.Read(&img); !ok {
				logger.Error("Cannot read from video capture")
				return fmt.Errorf("cannot read from video capture")
			}

			// フレームをJPEG形式にエンコード
			buf, err := gocv.IMEncode(".jpg", img)
			if err != nil {
				logger.Error("Failed to encode image")
				return err
			}

			chunk := &pb.StreamResponse{
				ChunkData:   buf.GetBytes(),
				ChunkSize:   uint64(len(buf.GetBytes())),
				Sequence:    0, // シーケンス番号を適切に設定
				EndOfStream: false,
			}

			// クライアントに送信
			if err := stream.Send(chunk); err != nil {
				log.Printf("Error sending stream: %v\n", err)
				return err
			}

			// 1秒待機
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil
}
