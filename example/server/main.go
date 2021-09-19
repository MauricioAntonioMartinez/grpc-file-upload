package main

import (
	"fmt"
	"log"
	"net"

	upload "github.com/MauricioAntonioMartinez/grpc-file-upload"
	file "github.com/MauricioAntonioMartinez/grpc-file-upload/example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct{}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
		return
	}

	s := grpc.NewServer()

	file.RegisterUploadServiceServer(s, &Server{})
	reflection.Register(s)

	fmt.Println("Staring gRPC server on port 8080")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
		return
	}

}

func (s *Server) UploadFile(stream file.UploadService_UploadFileServer) error {

	_ = upload.NewAwsStorage(upload.AwsConfig{
		Bucket: "obok-test",
		Key:    "video",
	})
	fsStorage := upload.NewFileSystemStorage(upload.FileSystemStorageConfig{
		Path: "./files/myimage.png",
	})

	up := upload.NewUploader(upload.UploaderConfig{
		MessageType:   &file.File{},
		MessageNumber: 0,
		MaxSize:       1024 * 1024 * 100 * 10, // 1Gig
		ChuckSize:     1024 * 1024 * 10,       // UploadPa r ts  o f 1Mb
	}, fsStorage)

	_, err := up.Upload(stream)
	fmt.Println(err)

	return stream.SendAndClose(&file.FileResponse{
		FileName: "test",
		Location: "mexico",
	})

}
