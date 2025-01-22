# S3 Panel Backend

## Getting started

S3 Panel Backend is a project to be used as backend for SnappCloud Unified Panel.

## APIs

Buckets
[x] Bucket List
[x] Bucket Quota
[x] Bucket Create
[] Bucket Delete

Users
[x] User Identification
[x] User Quota

Objects
[x] Object List
[] Object Delete
[] Object Download
[] Object Upload

## APIs Curls

### Bucket List

curl -XGET "127.0.0.1:8080/api/bucket/list" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Bucket Quota

curl -XGET "127.0.0.1:8080/api/bucket/quota" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### User Identification

curl -XGET "127.0.0.1:8080/api/user/id" -H "Content-Type: application/json" -H "access_key: X" -H "Authorization: Bearer X"

### User Quota

curl -XGET "127.0.0.1:8080/api/user/id" -H "Content-Type: application/json" -H "access_key: X" -H "Authorization: Bearer X"

### Object List



### Bucket Create