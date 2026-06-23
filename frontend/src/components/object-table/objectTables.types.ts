import type { IBucketObjectResponse } from '@/types/s3/buckets.types'

export type TObjectTablesProps = {
  isError: boolean
  isLoading: boolean
  bucket: string
  objectList: IBucketObjectResponse
  isSearch: boolean
  refetchObjects: () => void
  onShareObject: (objectName: string) => void
}
