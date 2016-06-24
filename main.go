package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	for {
		Start()
		time.Sleep(time.Hour * 1)
	}
}

func Start() {

	// Connect to S3
	svc := s3.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})

	// List buckets
	var params *s3.ListBucketsInput
	resp, err := svc.ListBuckets(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	count := 0
	var wg sync.WaitGroup
	// List objects
	for _, b := range resp.Buckets {
		lastKey := ""
		perPage := int64(1000)
		for {
			params := &s3.ListObjectsInput{
				Bucket:  b.Name,
				MaxKeys: &perPage,
				Marker:  &lastKey,
			}
			objects, err := svc.ListObjects(params)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// Is each object encrypted?
			wg.Add(1)
			go func(bucket s3.Bucket, obj []*s3.Object) {
				defer wg.Done()
				for _, o := range obj {
					if strings.HasSuffix(*o.Key, "/") || strings.HasPrefix(*o.Key, ".") {
						continue
					}
					if !CheckEncrypted(svc, bucket.Name, o.Key) {
						EncryptObject(svc, bucket.Name, o.Key)
						fmt.Println(*o.Key)
						count++
					}
				}

			}(*b, objects.Contents)

			// Paginate
			if *objects.IsTruncated {
				lastKey = *objects.Contents[len(objects.Contents)-1].Key
			} else {
				break
			}

		}

	}
	fmt.Println("Waiting")
	wg.Wait()
	fmt.Println("Encrypted", count, "objects")
}

// Returns the state of the object encryption
func CheckEncrypted(svc *s3.S3, bucket *string, key *string) bool {
	params := &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	}
	resp, _ := svc.HeadObject(params)
	return resp.ServerSideEncryption != nil
}

func EncryptObject(svc *s3.S3, bucket *string, key *string) error {
	params := &s3.CopyObjectInput{
		Bucket:               bucket,
		CopySource:           aws.String(*bucket + "/" + *key),
		Key:                  key,
		ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
	}
	_, err := svc.CopyObject(params)
	if err != nil {
		fmt.Println(err)
	}

	return err
}
