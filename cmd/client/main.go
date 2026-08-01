package main

import (
	pb "blob-service/gen/blob/v1"
	"context"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func printUsage() {
	fmt.Printf(`
		blob-service client\n
		usage: client <command> [flags]\n
		commands:
		  create   upload a file as a new blob
		  get      download a blob by id
		  list     list blobs with filters
		  update   replace a blob (full replace)
		  delete   remove a blob by id\n
		examples:
		  client create --file photo.jpg --tags demo,vacation
		  client get --id <uuid> --out downloaded.jpg
		  client list --name photo --tags demo --limit 10
		  client update --id <uuid> --file photo2.jpg --tags demo
		  client delete --id <uuid>
	`)
}

func runCreate(client pb.BlobServiceClient, args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	file := fs.String("file", "", "path to file to upload")
	tags := fs.String("tags", "", "comma-separated tags")

	fs.Parse(args)

	if *file == "" {
		fs.Usage()
		return
	}

	data, err := os.ReadFile(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read file:", err)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(*file))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var tagList []string
	if *tags != "" {
		tagList = strings.Split(*tags, ",")
	}

	resp, err := client.CreateBlob(context.Background(), &pb.CreateBlobRequest{
		Name:        filepath.Base(*file),
		ContentType: contentType,
		Data:        data,
		Tags:        tagList,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "creation failed:", err)
		return
	}

	fmt.Println(resp.Blob.Id)
}

func runGet(client pb.BlobServiceClient, args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	id := fs.String("id", "", "unique blob index")

	fs.Parse(args)

	if *id == "" {
		fs.Usage()
		return
	}

	resp, err := client.GetBlob(context.Background(), &pb.GetBlobRequest{
		Id: *id,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get blob:", err)
		return
	}

	fmt.Printf(resp.Blob.Id)
}

func runUpdate(client pb.BlobServiceClient, args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)

	id := fs.String("id", "", "unique blob index")
	file := fs.String("file", "", "path to new file (required)")
	tags := fs.String("tags", "", "comma-separated tags")

	fs.Parse(args)

	if *id == "" || *file == "" {
		fs.Usage()
		return
	}

	fs.Parse(args)

	var tagList []string
	if *tags != "" {
		tagList = strings.Split(*tags, ",")
	}

	data, err := os.ReadFile(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read file:", err)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(*file))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	resp, err := client.UpdateBlob(context.Background(), &pb.UpdateBlobRequest{
		Id:          *id,
		Name:        filepath.Base(*file),
		ContentType: contentType,
		Data:        data,
		Tags:        tagList,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "update failed:", err)
		return
	}
	fmt.Printf("updated: size=%d md5=%s\n", resp.Blob.Size, resp.Blob.Md5Hash)

}

func runDelete(client pb.BlobServiceClient, args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	id := fs.String("id", "", "unique blob index")

	fs.Parse(args)

	if *id == "" {
		fs.Usage()
		return
	}

	resp, err := client.DeleteBlob(context.Background(), &pb.DeleteBlobRequest{
		Id: *id,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to delete blob:", err)
		return
	}

	fmt.Printf(resp.Blob.Id)
}

func runList(client pb.BlobServiceClient, args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	name := fs.String("name", "", "blob name")
	contentType := fs.String("content_type", "", "content type of blob")
	tags := fs.String("tags", "", "comma-separated tags")
	limit := fs.Int("limit", 0, "max number of results")

	fs.Parse(args)

	if *limit <= 0 {
		fs.Usage()
		return
	}

	var tagFilters []string
	if *tags != "" {
		tagFilters = strings.Split(*tags, ",")
	}

	resp, err := client.ListBlobs(context.Background(), &pb.ListBlobsRequest{
		NameFilter:        *name,
		ContentTypeFilter: *contentType,
		TagFilters:        tagFilters,
		Limit:             int32(*limit),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "list failed:", err)
		return
	}

	for _, blob := range resp.Blobs {
		fmt.Printf("%s  %s  %d bytes  %v\n", blob.Id, blob.Name, blob.Size, blob.Tags)
	}
	fmt.Printf("total: %d\n", len(resp.Blobs))

}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewBlobServiceClient(conn)
	defer conn.Close()

	switch os.Args[1] {
	case "create":
		runCreate(client, os.Args[2:])
	case "get":
		runGet(client, os.Args[2:])
	case "update":
		runUpdate(client, os.Args[2:])
	case "delete":
		runDelete(client, os.Args[2:])
	case "list":
		runList(client, os.Args[2:])
	default:
		printUsage()
	}
}
