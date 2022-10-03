package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func setupConfig() {
	cli := flag.Bool("cli", false, "use the current `awscli` context, e.g. `./superkey_cleaner -cli`")
	access := flag.String("access", "", "which access key to use")
	secret := flag.String("secret", "", "which secret key to use")
	flag.Parse()

	if *access != "" && *secret != "" {
		fmt.Println("Loading from cli args...")
		cfg = try(config.LoadDefaultConfig(ctx,
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(*access, *secret, "cleaner"),
			),
		))
	} else if *cli {
		fmt.Println("Loading from awscli config...")
		cfg = try(config.LoadDefaultConfig(ctx))
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
