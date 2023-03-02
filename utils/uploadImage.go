package utils

import (
	"audiophile/models"
	cloud "cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	client := models.App{}
	var err error
	client.Ctx = context.Background()
	credentialsFile := option.WithCredentialsJSON([]byte(os.Getenv("FIRE_KEY")))
	app, err := firebase.NewApp(client.Ctx, nil, credentialsFile)
	if err != nil {
		return "", err
	}

	client.Client, err = app.Firestore(client.Ctx)
	if err != nil {
		return "", err
	}

	client.Storage, err = cloud.NewClient(client.Ctx, credentialsFile)
	if err != nil {
		return "", err
	}

	imagePath := fileHeader.Filename + strconv.Itoa(int(time.Now().Unix()))
	bucket := "audiophile-6470f.appspot.com"
	bucketStorage := client.Storage.Bucket(bucket).Object(imagePath).NewWriter(client.Ctx)
	_, err = io.Copy(bucketStorage, file)
	if err != nil {
		return "", err
	}
	if err := bucketStorage.Close(); err != nil {
		return "", err
	}
	return imagePath, nil
}
