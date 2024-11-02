package firebaseconnection

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	fscon "veles/db/FirestoreConnection"
	problem "veles/db/Problem"
	sconn "veles/db/StorageConnection"
)

var (
	ctx = context.Background()
	app *firebase.App
)

type FirebaseConnection struct {
	StorageConnection   sconn.StorageConnection
	FirestoreConnection fscon.FirestoreConnection
}

// Constructor of FirebaseConnection with path to a key.json
func New(path_to_key string, bucketName string) *FirebaseConnection {
	opt := option.WithCredentialsFile(path_to_key)
	var err error
	app, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
		return nil
	}
	firesoreClient, err_fs := app.Firestore(ctx)
	if err_fs != nil {
		fmt.Printf("error initializing app: %v", err_fs)
		return nil
	}
	firesoreConnection := fscon.New(firesoreClient)

	storageClient, err_st := app.Storage(ctx)
	if err_st != nil {
		fmt.Printf("error initializing app: %v", err_st)
		return nil
	}

	storageConnection := sconn.New(storageClient, bucketName)

	c := FirebaseConnection{storageConnection, firesoreConnection}

	return &c
}

func (fc FirebaseConnection) GetProblemsList() []problem.Problem {
	return fc.FirestoreConnection.GetProblemsList()
}

func (fc FirebaseConnection) SaveProblemSubmission(number string,
	text string, answer string, ext []string, comment string) []string {
	return fc.FirestoreConnection.SaveProblemSubmission(number, text, answer, ext, comment)
}

func (fc FirebaseConnection) SaveFiles(files []multipart.File, number []string) error {
	return fc.StorageConnection.SaveFiles(files, number)
}

func (fc FirebaseConnection) GetReadyMdFile(number string) io.Reader {
	return fc.StorageConnection.GetReadyMdFile(number)
}

func (fc FirebaseConnection) GetNumsOfSubmissions(number string) int {
	return fc.FirestoreConnection.GetNumsOfSubmissions(number)
}
