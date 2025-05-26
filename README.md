# S3 Panel Backend

## Getting started

S3 Panel Backend is a project to be used as backend for SnappCloud Unified Panel.

## APIs

Buckets
- [x] Bucket List (with Pagination)
- [x] Bucket Quota (with Pagination)
- [x] Bucket Create
- [x] Bucket Delete

Users
- [x] User Identification
- [x] User Quota

Objects
- [x] Objects Delete
- [x] Object List (with Pagination)
- [x] Object Download
- [x] Object Upload (Multiple Files)
- [x] Object Head
- [x] Object Temporary Link (with Custom Expiration)

## APIs Curls

### Bucket List

curl -XGET "127.0.0.1:8080/api/bucket/list?page=1&max_keys=10" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Bucket Quota

curl -XGET "127.0.0.1:8080/api/bucket/quota?page=1&max_keys=10" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Bucket Create

curl -XPOST "127.0.0.1:8080/api/bucket/create" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X" -d '{"bucket": "X"}'

### Bucket Delete

curl -XDELETE "127.0.0.1:8080/api/bucket/delete" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X" -d '{"bucket": "X"}'

### User Identification

curl -XGET "127.0.0.1:8080/api/user/id" -H "Content-Type: application/json" -H "access_key: X" -H "Authorization: Bearer X"

### User Quota

curl -XGET "127.0.0.1:8080/api/user/quota" -H "Content-Type: application/json" -H "access_key: X" -H "Authorization: Bearer X"

### Object List

curl -XGET "127.0.0.1:8080/api/object/list?max_keys=10&page=1&bucket=X" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Objects Delete

curl -XDELETE "127.0.0.1:8080/api/object/delete?bucket=X&objects=X1&objects=X2" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Object Upload

curl -XPOST "127.0.0.1:8080/api/object/upload" -H "Content-Type: multipart/form-data" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X" -F 'files=@/path/to/file' -F 'files=@/path/to/file2' -F 'bucket="X"'

### Object Download

curl -XGET "127.0.0.1:8080/api/object/download?bucket=X&object=X" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X" | jq

### Object Head

curl -XGET "127.0.0.1:8080/api/object/head?bucket=X&object=X" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Object Share

curl -XGET "127.0.0.1:8080/api/object/share?bucket=X&object=X&expiration=1d" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X" | jq

## TODO

Object

- [] Object Copy (source/destination bucket)
- [] Object Set/Get Lock
- [] Object Set/Get Retention

Bucket

- [] Bucket Set/Get Retention

General

- [] Support Directories
- [] Bucket Set/Get Encryption/SSE
- [] Bucket Set/Get Replication
