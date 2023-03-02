package utils

import (
	"audiophile/helpers"
	"audiophile/models"
	cloud "cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"os"
	"time"
)

func GetImageUrl(variantId uint, limit, page int64) ([]string, error) {
	imagePaths, tx := helpers.GetImagesPath(variantId, limit, page)
	if tx.Error != nil {
		return nil, tx.Error
	}

	client := models.App{}
	var err error
	client.Ctx = context.Background()
	credentialsFile := option.WithCredentialsJSON([]byte(os.Getenv("FIRE_KEY")))
	//fmt.Println(credentialsFile)
	app, err := firebase.NewApp(client.Ctx, nil, credentialsFile)
	if err != nil {
		return nil, err
	}

	client.Client, err = app.Firestore(client.Ctx)
	if err != nil {
		return nil, err
	}

	client.Storage, err = cloud.NewClient(client.Ctx, credentialsFile)
	if err != nil {
		return nil, err
	}
	signedUrl := &cloud.SignedURLOptions{
		Scheme:  cloud.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}
	bucket := "audiophile-6470f.appspot.com"
	bucketHandler := client.Storage.Bucket(bucket)
	var urls []string
	for i := 0; i < len(imagePaths); i++ {
		url, err := bucketHandler.SignedURL(imagePaths[i], signedUrl)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
