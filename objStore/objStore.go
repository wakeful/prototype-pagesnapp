package objStore

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go"
)

type ObjStore struct {
	client     *minio.Client
	bucketName string
}

// create "unique" sha256 string
func (o ObjStore) createObjName(name string) string {
	hash := sha256.New()
	hash.Write([]byte(time.Now().String() + name))
	objName := hex.EncodeToString(hash.Sum(nil))

	return objName
}

// upload PNG file to object store
func (o ObjStore) SavePngFile(path string) (name string, size int64, err error) {

	name = o.createObjName(path) + ".png"

	size, err = o.client.FPutObject(o.bucketName, name, path, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return name, 0, err
	}

	return name, size, nil
}

// generate download link for object
func (o ObjStore) GenerateAccessLink(objectName string) (link *url.URL, err error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"pic.png\"")

	link, err = o.client.PresignedGetObject(o.bucketName, objectName, 2*time.Hour, reqParams)
	if err != nil {
		return link, err
	}

	return link, nil
}

// configure new object store
func NewObjStore(objStoreUrl, accessKey, secretKey, bucketName string) (*ObjStore, error) {

	// configure new Minio client
	mClient, err := minio.New(objStoreUrl, accessKey, secretKey, false)
	if err != nil {
		return nil, err
	}

	// create bucket
	if err = mClient.MakeBucket(bucketName, bucketName); err != nil {
		// Check if bucket already exists
		exists, err := mClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("bucket already exists: %s\n", bucketName)
		} else {
			return nil, err
		}
	}

	return &ObjStore{
		client:     mClient,
		bucketName: bucketName,
	}, nil
}
