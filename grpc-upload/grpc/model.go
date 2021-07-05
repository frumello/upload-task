package grpc

import (
	"grpc-upload/pb"
)

func NewUploadRequest(fileName string) *pb.UploadRequest {
	return &pb.UploadRequest{
		Data: &pb.UploadRequest_FileName{
			FileName: fileName,
		},
	}
}
