package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Generate a random string to use with the container and package name.
func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

func main() {
	fmt.Printf("Azure Blob storage quick start sample\n")

	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic(errors.New("AZURE_STORAGE_ACCOUNT_NAME could not be found"))
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic(errors.New("AZURE_STORAGE_ACCOUNT_KEY could not be found"))
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	print(url, "\n")
	serviceClient, _ := azblob.NewServiceClientWithSharedKey(url, credential, nil)

	ctx := context.Background()

	// Create the container
	containerName := fmt.Sprintf("quickstart-%s", randomString())

	fmt.Printf("Creating a container named %s\n", containerName)

	containerClient, _ := serviceClient.NewContainerClient(containerName)
	_, err = containerClient.Create(ctx, nil)
	if err != nil {
		log.Fatal("Error Creating Container: " + err.Error())
	}

	fmt.Printf("Creating a dummy file to test the upload and download\n")

	data := []byte("\nhello world this is a blob\n")
	blobName := "quickstartblob" + "-" + randomString() + ".txt"

	blobClient, _ := containerClient.NewBlockBlobClient(url + containerName + "/" + blobName)

	// Upload to data to blob storage
	_, err = blobClient.UploadBuffer(ctx, data, azblob.UploadOption{})

	if err != nil {
		log.Fatalf("Failure to upload to blob: %+v", err)
	}

	// List the blobs in the container
	fmt.Println("Listing the blobs in the container:")

	pager := containerClient.ListBlobsFlat(nil)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, v := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			fmt.Println(*v.Name)
		}
	}

	if err = pager.Err(); err != nil {
		log.Fatalf("Failure to list blobs: %+v", err)
	}

	// Download the blob
	get, err := blobClient.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	downloadedData := &bytes.Buffer{}
	reader := get.Body(&azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		log.Fatal(err)
	}
	err = reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(downloadedData.String())
	fmt.Println("Blob container creation and file upload complete!")
	/*
		fmt.Printf("Press enter key to delete the blob, example container, and exit the application.\n")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Printf("Cleaning up.\n")

		// Delete the blob
		fmt.Printf("Deleting the blob " + blobName + "\n")

		_, err = blobClient.Delete(ctx, nil)
		if err != nil {
			log.Fatalf("Failure: %+v", err)
		}

		// Delete the container
		fmt.Printf("Deleting the blob " + containerName + "\n")
		_, err = containerClient.Delete(ctx, nil)

		if err != nil {
			log.Fatalf("Failure: %+v", err)
		}
	*/
}
