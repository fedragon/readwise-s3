package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
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

	req := readwise.GetHighlightsRequest{
		PageSize: 1000,
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

	u, _ := url.Parse("https://readwise.io/api/v2/highlights/")
	count := 1
	for {
		fmt.Printf("GET %v\n", u)
		res, err := readwise.Get(u, appCfg.ReadwiseApiToken, &req)
		if err != nil {
			log.Fatal(err)
		}

		objectName := fmt.Sprintf("page-%03d.json", count)
		if err := aws.UploadToS3(ctx, client, &bucket, &objectName, bytes.NewReader(res.Results)); err != nil {
			log.Fatal(err)
		}

		if res.Next == nil {
			break
		}

		count++
		u = &res.Next.Url
	}
}
