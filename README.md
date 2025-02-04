# S3 Panel Backend

## Getting started

S3 Panel Backend is a project to be used as backend for SnappCloud Unified Panel.

## APIs

Buckets
- [x] Bucket List
- [x] Bucket Quota
- [x] Bucket Create
- [x] Bucket Delete

Users
- [x] User Identification
- [x] User Quota

Objects
- [x] Object List
- [x] Object Delete
- [] Object Download
- [] Object Upload

## APIs Curls

### Bucket List

curl -XGET "127.0.0.1:8080/api/bucket/list" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Bucket Quota

curl -XGET "127.0.0.1:8080/api/bucket/quota" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Bucket Create

curl -XPOST "127.0.0.1:8080/api/bucket/create" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X" -d '{"bucket": "X"}'

### Bucket Delete

curl -XDELETE "127.0.0.1:8080/api/bucket/delete" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X" -d '{"bucket": "X"}'

### User Identification

curl -XGET "127.0.0.1:8080/api/user/id" -H "Content-Type: application/json" -H "access_key: X" -H "Authorization: Bearer X"

### User Quota

curl -XGET "127.0.0.1:8080/api/user/id" -H "Content-Type: application/json" -H "access_key: X" -H "Authorization: Bearer X"

### Object List

curl -XGET "127.0.0.1:8080/api/object/list?max_keys=10&page=1&bucket=X" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X"

### Objects Delete

curl -XDELETE "127.0.0.1:8080/api/object/delete" -H "Content-Type: application/json" -H "access_key: X" -H "secret_key: X" -H "Authorization: Bearer X" -d '{"bucket": "X", "objects":["object1.png", "object2.txt"]}'

## TODO

Object

[] Object Permanent Link
[] Object Temporary Link
[] Search Objects
[] Object Copy (source/destination bucket)
[] Object Set/Get Lock
[] Object Set/Get Retention

Bucket

[] Bucket Set/Get Retention

General

[] Support Directories
[] Bucket Set/Get Encryption/SSE
[] Bucket Set/Get Replication

## Release Plan

This list is prioritized from top to bottom. we will work on 2 items in a sprint.

- Object Download
- Object Upload