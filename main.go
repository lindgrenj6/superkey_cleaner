package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	cost "github.com/aws/aws-sdk-go-v2/service/costandusagereportservice"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var ctx = context.Background()
var cfg aws.Config

func main() {
	setupConfig()

	costClient := cost.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)

	// --------------------
	// Deleting Cost and Usage Reports (to get around ReportLimitExceeded Exception)
	// --------------------

	reports, err := costClient.DescribeReportDefinitions(ctx, nil)
	panicOn(err)

	// pulling out all of the reports that have the 'koku-' prefix
	reportNames := make([]string, 0, len(reports.ReportDefinitions))
	for _, report := range reports.ReportDefinitions {
		if strings.HasPrefix(*report.ReportName, "koku-") {
			reportNames = append(reportNames, *report.ReportName)
		}
	}

	fmt.Printf("Deleting Reports: [%v]\n", strings.Join(reportNames, ","))
	fmt.Println("Hit Ctrl-C to abort, sleeping 5s...")
	time.Sleep(5 * time.Second)

	for _, report := range reportNames {
		fmt.Printf("\tDeleting report %v\n", report)

		_, err := costClient.DeleteReportDefinition(ctx, &cost.DeleteReportDefinitionInput{
			ReportName: &report,
		})
		panicOn(err)
	}

	// --------------------
	// Deleting related s3 buckets
	// --------------------

	buckets, err := s3Client.ListBuckets(ctx, nil)
	panicOn(err)

	bucketNames := make([]string, len(buckets.Buckets))
	for _, bucket := range buckets.Buckets {
		for _, report := range reportNames {
			pieces := strings.Split(report, "-")
			guid := pieces[len(pieces)-1]

			if strings.HasSuffix(*bucket.Name, guid) {
				bucketNames = append(bucketNames, *bucket.Name)
			}
		}
	}

	fmt.Printf("Deleting Buckets: [%v]\n", strings.Join(bucketNames, ","))
	fmt.Println("Hit Ctrl-C to abort, sleeping 5s...")
	time.Sleep(5 * time.Second)

	for _, bucket := range bucketNames {
		fmt.Printf("\tDeleting bucket %v\n", bucket)

		objects, err := s3Client.ListObjects(ctx, &s3.ListObjectsInput{
			Bucket: &bucket,
		})
		panicOn(err)

		for _, object := range objects.Contents {
			_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    object.Key,
			})
			panicOn(err)
		}

		_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			Bucket: &bucket,
		})
		panicOn(err)
	}

	fmt.Println("Deleted all Reports + Attached buckets successsfully!")
}

func setupConfig() {
	cli := flag.Bool("cli", false, "use the current `awscli` context, e.g. `./superkey_cleaner -cli`")
	access := flag.String("access", "", "which access key to use")
	secret := flag.String("secret", "", "which secret key to use")
	flag.Parse()

	var err error

	if *access != "" && *secret != "" {
		fmt.Println("Loading from cli args...")
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(*access, *secret, "cleaner"),
			),
		)
		if err != nil {
			panic("configuration error, " + err.Error())
		}
	} else if *cli {
		fmt.Println("Loading from awscli config...")
		cfg, err = config.LoadDefaultConfig(ctx)
		if err != nil {
			panic("configuration error, " + err.Error())
		}
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
