# readwise-s3

Performs a backup of all [Readwise.io](https://readwise.io) highlights to AWS S3.

## Prerequisites

- Have a working, configured [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html)
- Get a Readwise access token [here](https://readwise.io/access_token)

## Usage

```shell
make build
READWISE_API_TOKEN='<your_token>' S3_BUCKET_PREFIX='<whatever>' bin/readwise-s3
```

This will generate a bucket named `<whatever>-readwiseio-<current_timestamp>`, where each object is an array of Highlight JSON objects defined as described [here](https://readwise.io/api_deets).

Its content will look like this:

```shell
aws s3 ls s3://<bucket_name>
books-001.json
...
books-00N.json
highlights-001.json
...
highlights-00N.json
```