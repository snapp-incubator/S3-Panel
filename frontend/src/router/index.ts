import { createRoute, createRouter, redirect } from '@tanstack/react-router'
import { lazy } from 'react'

import { handleAuthRedirect } from '@/lib/utils'

import { rootRoute } from './root'

const S3BucketPage = lazy(() => import('@/pages/objectStorage/s3Bucket'))
const BucketsPage = lazy(() => import('@/pages/objectStorage/s3Bucket/buckets'))
const BucketObjectsPage = lazy(
  () => import('@/pages/objectStorage/s3Bucket/buckets/bucket-objects')
)

export const HomeRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: () => {
    throw redirect({ to: '/object-storage/s3-bucket' })
  }
})

export const s3BucketRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/object-storage/s3-bucket',
  component: S3BucketPage,
  beforeLoad: () =>
    handleAuthRedirect('/object-storage/s3-bucket/buckets', false)
})

export const BucketsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/object-storage/s3-bucket/buckets',
  component: BucketsPage,
  beforeLoad: () => handleAuthRedirect('/object-storage/s3-bucket')
})

export const BucketObjectsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/object-storage/s3-bucket/buckets/$bucketName',
  component: BucketObjectsPage,
  beforeLoad: () => handleAuthRedirect('/object-storage/s3-bucket')
})

const routeTree = rootRoute.addChildren([
  HomeRoute,
  s3BucketRoute,
  BucketsRoute,
  BucketObjectsRoute
])

export const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}
