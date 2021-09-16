package main

import (
	"fmt"
	"log"
	"net"

	file "github.com/MauricioAntonioMartinez/grpc-file-upload/proto"
	"github.com/MauricioAntonioMartinez/grpc-file-upload/upload"
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
	// buffer := bytes.Buffer{}
	// imageSize := 0
	// for {
	// 	err := contextErr(stream.Context())
	// 	if err != nil {
	// 		return err
	// 	}
	// 	req, err := stream.Recv()

	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		return status.Errorf(codes.Internal, "Error reading input")
	// 	}
	// 	data := req.GetData()

	// 	size := len(data)
	// 	imageSize += size

	// 	if imageSize > maxIMageSize {
	// 		return logError(status.Errorf(codes.InvalidArgument, "image size too large"))
	// 	}

	// 	_, err = buffer.Write(data)

	// 	if err != nil {
	// 		return logError(status.Errorf(codes.Internal, "Error reading the bytes into buffer."))
	// 	}

	// }

	// name, err := saveImage(buffer)
	// if err != nil {
	// 	return status.Errorf(codes.Internal, "Error saving the image")
	// }

	store := upload.NewAwsStorage(upload.AwsConfig{})

	up := upload.NewUploader(upload.UploaderConfig{
		MessageType:   &file.File{},
		MessageNumber: 0,
		MaxSize:       1024 * 1024 * 100 * 10, // 1Gig
		ChuckSize:     1024,                   // UploadParts of 1Mb
	}, store)

	_, err := up.Upload(stream)
	fmt.Println(err)

	return stream.SendAndClose(&file.FileResponse{
		FileName: "test",
		Location: "mexico",
	})

}

// func saveImage(buffer bytes.Buffer) (string, error) {

// 	imageID, err := uuid.NewRandom()
// 	if err != nil {
// 		return "", err
// 	}

// 	err = os.WriteFile(fmt.Sprintf("%s.%s", imageID.String(), "png"), buffer.Bytes(), 0644)
// 	return imageID.String(), err

// }
