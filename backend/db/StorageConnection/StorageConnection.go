package storageconnection

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	altstorage "cloud.google.com/go/storage"
	storage "firebase.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var (
	ctx = context.Background()
)

type StorageConnection struct {
	bucketName string
	client     *storage.Client
}

func New(client *storage.Client, bucketName string) StorageConnection {
	sc := StorageConnection{client: client, bucketName: bucketName}
	return sc
}

func (sc StorageConnection) SaveFiles(files []multipart.File, filenames []string) error {
	bucket, _ := sc.client.Bucket(sc.bucketName)

	for indx, file := range files {
		object := bucket.Object(strings.Split(filenames[indx], "(")[0] + "/" + filenames[indx])

		wc := object.NewWriter(ctx)

		if _, err := io.Copy(wc, file); err != nil {
			return err
		}
		if err := wc.Close(); err != nil {
			return err
		}

	}
	return nil
}

func (sc StorageConnection) GetImages(number string) []io.Reader {

	bucket, _ := sc.client.Bucket(sc.bucketName)

	var names []string
	var filesreaders []io.Reader
	q := &altstorage.Query{Prefix: number}
	it := bucket.Objects(ctx, q)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		if strings.Contains(attrs.Name, "final") {
			names = append(names, attrs.Name)
		}
	}

	fmt.Println("names", names)

	for _, name := range names {
		obj_handler := bucket.Object(name)
		file_reader, _ := obj_handler.NewReader(ctx)
		file_reader.Close()
		filesreaders = append(filesreaders, file_reader)
	}

	return filesreaders
}

func (sc StorageConnection) GetReadyMdFile(number string) io.Reader {
	bucket, _ := sc.client.Bucket(sc.bucketName)
	obj_handler := bucket.Object("Final/" +
		strings.ReplaceAll(number, ".", "-") + "/" +
		strings.ReplaceAll(number, ".", "-") + "urled.md")
	file_reader, _ := obj_handler.NewReader(ctx)
	file_reader.Close()
	return file_reader
}
