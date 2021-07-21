# Superkey Cleaner!

This tool will delete all superkey related reports (and their respective s3 buckets) from the aws account currently set up for use with the `awscli` tool. 

### Pre-requisites
1. aws CLI set up locally pointing at the account you want to clean.
2. Go toolchain, (`dnf install golang` or `brew install golang` or https://golang.org/doc/install)

### Running
There is a Makefile present for easy operation. 

To build, type `make build`. To run type `make run` which will build + run the cleaner. Run as much as you want but it'll do nothing if the account is empty!
