package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	file "github.com/MauricioAntonioMartinez/grpc-file-upload/example/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return
	}
	client := file.NewUploadServiceClient(lis)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	req, err := client.UploadFile(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	f, err := os.Open("./vid.mkv")
	reader := bufio.NewReader(f)
	checkErr(err)
	defer f.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			fmt.Println("Done!")
			break
		}
		if err != nil {
			log.Fatal(err)
			break
		}

		req.Send(&file.File{
			Data: buffer[:n],
		})
	}
	res, err := req.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}
	log.Printf("image uploaded with name: %s, location: %s", res.GetFileName(), res.GetLocation())

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic("Something went wrong my son.")
	}
}

func getFile() *bufio.Reader {
	f, err := os.Open("./gopher.png")
	checkErr(err)
	reader := bufio.NewReader(f)
	return reader

}
