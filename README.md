# readwise-s3

Performs a backup of all [Readwise.io](https://readwise.io) highlights to AWS S3.

## Prerequisites

- Have a working, configured [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html)
- Get a Readwise access token [here](https://readwise.io/access_token)

## Usage

Grab one of the binaries from the `Releases` page, then

```shell
READWISE_API_TOKEN='<your_token>' S3_BUCKET_PREFIX='<whatever>' /path/to/readwise-s3
```

This will generate a bucket named `<whatever>-readwiseio-<current_timestamp>`, whose contents will look like this:

```shell
aws s3 ls s3://<bucket_name>
books-001.json
...
books-00N.json
highlights-001.json
...
highlights-00N.json
```

Each object contains an array of `Book`/`Highlight` JSON objects defined as described [here](https://readwise.io/api_deets).
