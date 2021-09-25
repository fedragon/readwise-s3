package aws

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func CreateS3Bucket(ctx context.Context, client *s3.Client, region *string, bucketName *string) error {
	if _, err := client.CreateBucket(
		ctx,
		&s3.CreateBucketInput{
			Bucket: bucketName,
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(*region),
			},
		}); err != nil {
		return fmt.Errorf("could not create bucket %v: %v", *bucketName, err)
	}

	if _, err := client.PutBucketEncryption(
		ctx,
		&s3.PutBucketEncryptionInput{
			Bucket: bucketName,
			ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
				Rules: []types.ServerSideEncryptionRule{{
					ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
						SSEAlgorithm: types.ServerSideEncryptionAes256,
					}},
				},
			},
		}); err != nil {
		return fmt.Errorf("could not encrypt bucket %v: %v", *bucketName, err)
	}

	return nil
}

func UploadToS3(ctx context.Context, client *s3.Client, bucketName *string, objectName *string, body io.Reader) error {
	uploader := manager.NewUploader(client)
	if _, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: bucketName,
		Key:    objectName,
		Body:   body,
	}); err != nil {
		return fmt.Errorf("got error uploading object: %v", err)
	}

	return nil
}
