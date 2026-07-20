import { createKeyStore } from './query-keys'

export const s3Keys = {
  user: 'user',
  objects: 'objects',
  buckets: 'buckets'
}

export const userKeys = createKeyStore(s3Keys.user, {
  allQuota: () => [],
  details: () => []
})

export const bucketObjectKeys = createKeyStore(s3Keys.objects, {
  all: (
    bucketName: string,
    page: number,
    searchValue?: string,
    prefix?: string
  ) => [bucketName, page, searchValue, prefix]
})

export const bucketsKeys = createKeyStore(s3Keys.buckets, {
  all: (maxItems: number = 10, page: number = 1, searchValue?: string) => [
    maxItems,
    page,
    searchValue
  ]
})
