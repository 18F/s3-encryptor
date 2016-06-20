## S3 Encryptor

Ensures that all your files in S3 are encrypted. Loops through all the buckets and objects and copies the one no encrypted to encrypted versions.

### How to run

```
export AWS_ACCESS_KEY_ID=YOUR_KEY_ID
export AWS_SECRET_ACCESS_KEY=YOUR_KEY
export AWS_REGION=YOUR_REGION

go run main.go
```