package grpc

import (
	"fmt"
	"io"
	"os"
	"sync"

	"upload-server/pb"
	"upload-server/util"
)

var validFileExtensions = map[string]bool{
	"image/jpeg": true,
	"image/gif":  true,
	"image/png":  true,
}

type Handler struct {
	pb.UnsafeImageServiceServer
	lock sync.RWMutex

	httpTarget string
}

func NewHandler(httpTarget string) *Handler {
	return &Handler{
		httpTarget: httpTarget,
	}
}

func (h *Handler) UploadImage(stream pb.ImageService_UploadImageServer) error {
	fileName, err := h.receiveStream(stream)
	if err != nil {
		return err
	}
	fileType, err := util.GetFileContentType(fileName)
	if err != nil {
		return ErrUnexpectedly(err)
	}
	if !validFileExtensions[fileType] {
		return ErrBadFileType()
	}
	return stream.SendAndClose(&pb.UploadResponse{Url: h.buildResponseUrl(fileName)})
}

func (h *Handler) receiveStream(stream pb.ImageService_UploadImageServer) (string, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	var file *os.File
	defer file.Close()

	fileName := ""

	for {
		err := ContextError(stream.Context())
		if err != nil {
			return "", err
		}
		// Start receiving stream messages from the client
		req, err := stream.Recv()
		// Check if the stream has finished
		if err != nil {
			if err == io.EOF {
				break
			}
			// Close the connection and return the response to the client
			return "", ErrUnexpectedly(err)
		}

		if fileName == "" {
			fileName = req.GetFileName()
			if fileName == "" {
				return "", ErrBadFileName()
			}
			if util.FileExists(fileName) {
				return "", ErrFileNameExist(fileName)
			}
			file, err = os.Create(fileName)
			if err != nil {
				return "", ErrUnexpectedly(err)
			}
		}
		if err = util.WriteChunkToFile(file, req.GetChunkData()); err != nil {
			return "", ErrUnexpectedly(err)
		}
	}
	_ = file.Close()
	return fileName, nil
}

func (h *Handler) buildResponseUrl(fileName string) string {
	return fmt.Sprintf("http://localhost%s/%s", h.httpTarget, fileName)
}
