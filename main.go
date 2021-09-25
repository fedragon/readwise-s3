package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fedragon/readwise-client/aws"
	"github.com/fedragon/readwise-client/readwise"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	ReadwiseApiToken string `envconfig:"READWISE_API_TOKEN" required:"true"`
	S3BucketPrefix   string `envconfig:"S3_BUCKET_PREFIX" required:"true"`
}

func main() {
	var appCfg AppConfig
	if err := envconfig.Process("", &appCfg); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	client := s3.NewFromConfig(cfg)

	bucket := fmt.Sprintf("%s-readwiseio-%v", appCfg.S3BucketPrefix, time.Now().UTC().Format("20060102150405"))
	if err := aws.CreateS3Bucket(ctx, client, &cfg.Region, &bucket); err != nil {
		log.Fatal(err)
	}

	if err := backup(ctx, "books", appCfg.ReadwiseApiToken, client, &bucket); err != nil {
		log.Fatal(err)
	}

	if err := backup(ctx, "highlights", appCfg.ReadwiseApiToken, client, &bucket); err != nil {
		log.Fatal(err)
	}
}

func backup(ctx context.Context, resource string, token string, client *s3.Client, bucket *string) error {
	count := 1
	u, _ := url.Parse(fmt.Sprintf("https://readwise.io/api/v2/%s/", resource))
	reqCtx := readwise.NewListRequestContext(&http.Client{Timeout: time.Second * 30}, u, token, 1000)

	for {
		fmt.Printf("GET %v\n", u)
		res, err := readwise.Get(reqCtx)
		if err != nil {
			log.Fatal(err)
		}

		objectName := fmt.Sprintf("%s-%03d.json", resource, count)
		if err := aws.UploadToS3(ctx, client, bucket, &objectName, bytes.NewReader(res.Results)); err != nil {
			log.Fatal(err)
		}

		if res.Next == nil {
			break
		}

		count++
		reqCtx.SetURL(&res.Next.Url)
	}

	return nil
}
