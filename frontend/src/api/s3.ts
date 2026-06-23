import centralClient from '@/services/http/centralClient'
import { buildQueryString } from '@/services/http/query-client'
import type { IRegionsResponse } from '@/types/regions.types'
import type {
  IBucketObjectResponse,
  IBucketResponse,
  IBucketsListResponse,
  IUserDetailsResponse,
  IUserQuotaResponse
} from '@/types/s3/buckets.types'

const fetchRegions = async () => {
  const res = await centralClient.get<IRegionsResponse>('/s3/api/regions')

  return res.data
}

const createBucket = async ({ bucket }: { bucket: string }) => {
  const res = await centralClient.post('/s3/api/bucket/create', {
    bucket
  })

  return res
}

const fetchUserQuota = async () => {
  const res = await centralClient.get<IUserQuotaResponse>('/s3/api/user/quota')

  return res.data
}

const fetchObjects = async (
  bucket: string,
  page: number,
  searchValue?: string
) => {
  const params = {
    bucket,
    max_keys: 20,
    page,
    search_string: searchValue || undefined
  }

  const res = await centralClient.get<IBucketObjectResponse>(
    `/s3/api/object/list?${buildQueryString(params)}`
  )

  return res.data
}

const fetchBucketsQuota = async (
  max_keys: number,
  page: number,
  searchValue?: string
) => {
  const params = {
    max_keys,
    page,
    search_string: searchValue || undefined
  }

  const res = await centralClient.get<IBucketResponse>(
    `/s3/api/bucket/quota?${buildQueryString(params)}`
  )

  return res.data
}

const fetchBucketsList = async () => {
  const res = await centralClient.get<IBucketsListResponse>(
    `/s3/api/bucket/list?${buildQueryString({
      max_keys: 10,
      page: 1
    })}`
  )

  return res
}

const fetchUserDetails = async () => {
  const res = await centralClient.get<IUserDetailsResponse>('/s3/api/user/id')

  return res.data
}

const deleteObject = async (bucket: string, objects: string[]) => {
  const res = await centralClient.delete<{
    deleted: boolean
  }>(
    `/s3/api/object/delete?${buildQueryString({
      bucket,
      objects: objects.join(',')
    })}`
  )

  return res
}

const deleteBucketApi = async (bucket: string) => {
  const res = await centralClient.delete(
    `/s3/api/bucket/delete?${buildQueryString({ bucket })}`
  )

  return res
}

const uploadObjects = async (
  formData: FormData,
  onUploadProgress?: (percent: number) => void
) => {
  const response = await centralClient.upload<{ created: boolean }>(
    '/s3/api/object/upload',
    formData,
    undefined,
    onUploadProgress
  )

  return response.data
}

const downloadObject = async (
  downloadLink: string,
  signal?: AbortSignal
): Promise<{ url: string }> => {
  const response = await centralClient.get<{ url: string }>(downloadLink, {
    signal
  })

  return response.data
}

interface IShareLinkResponse {
  url: string
}

const shareLink = async (
  bucket: string,
  object: string,
  expiration: string
): Promise<IShareLinkResponse> => {
  const params = {
    bucket,
    object,
    expiration
  }

  const response = await centralClient.get<IShareLinkResponse>(
    `/s3/api/object/share?${buildQueryString(params)}`
  )

  return response.data
}

export {
  createBucket,
  deleteBucketApi,
  deleteObject,
  downloadObject,
  fetchBucketsList,
  fetchBucketsQuota,
  fetchObjects,
  fetchRegions,
  fetchUserDetails,
  fetchUserQuota,
  shareLink,
  uploadObjects
}
