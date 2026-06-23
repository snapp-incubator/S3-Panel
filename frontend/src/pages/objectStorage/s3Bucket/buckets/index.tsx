import { useState, useEffect, useDeferredValue } from 'react'

import { useQuery } from '@tanstack/react-query'
import { useEffectOnce } from 'react-use'

import { fetchBucketsQuota } from '@/api/s3'
import { bucketsKeys } from '@/api/s3Keys'
import { AlertMessage } from '@/components/alert-message'
import BucketCard from '@/components/bucket-card'
import CreateBucket from '@/components/create-bucket'
import CustomPagination from '@/components/custom-pagination'
import ErrorState from '@/components/error-state'
import { useTitle } from '@/components/providers/titleProvider'
import SearchField from '@/components/search-field'
import { Button } from '@/components/shadcn/button'
import BucketCardSkeleton from '@/components/skeletons/BucketCardSkeleton'
import UserQuota from '@/components/user-quota'
import { t } from '@/i18n'

export default function Buckets() {
  const { setTitle } = useTitle()
  const [openCreate, setOpenCreate] = useState(false)
  const [page, setPage] = useState(1)
  const [initialTotalPages, setInitialTotalPages] = useState<number | null>(
    null
  )
  const [searchValue, setSearchValue] = useState('')
  const deferredSearch = useDeferredValue(searchValue)
  const maxItems: number = 10

  const {
    data: buckets,
    isFetching,
    isError,
    refetch
  } = useQuery({
    queryFn: () => fetchBucketsQuota(maxItems, page, deferredSearch),
    queryKey: bucketsKeys.all(maxItems, page, deferredSearch)
  })

  useEffectOnce(() => {
    setTitle('SnappCloud - s3 Buckets')
  })

  useEffect(() => {
    if (buckets && initialTotalPages === null) {
      setInitialTotalPages(buckets.total_pages)
    }
  }, [buckets, initialTotalPages])

  const returnBuckets = () => {
    if (buckets && buckets.items) {
      return buckets.items.map(bucket => (
        <BucketCard {...bucket} key={bucket.bucket} />
      ))
    }

    return (
      <AlertMessage
        title={t('empty_buckets')}
        message={t('create_first_bucket')}
      />
    )
  }

  return (
    <div>
      <div className="flex flex-col justify-between gap-4 md:flex-row md:gap-0">
        <h2 className="text-3xl">{t('s3_bucket')}</h2>
      </div>
      <UserQuota />
      <span className="mt-2 block text-xl font-semibold">{t('buckets')}</span>
      <div className="mt-4 flex items-center justify-between">
        <SearchField
          value={searchValue}
          onChange={value => {
            setSearchValue(value)
            setPage(1)
          }}
        />
        <Button size="sm" onClick={() => setOpenCreate(true)}>
          {t('create_bucket')}
        </Button>
      </div>

      {isError ? (
        <div className="mt-20">
          <ErrorState />
        </div>
      ) : (
        <>
          <div className="mt-12 grid grid-cols-1 gap-6 lg:grid-cols-2 xl:grid-cols-3">
            {isFetching ? <BucketCardSkeleton count={6} /> : returnBuckets()}
          </div>
          <div className="relative mt-5">
            {initialTotalPages && initialTotalPages > 1 && !searchValue ? (
              <CustomPagination
                currentPage={page}
                totalPages={initialTotalPages}
                onPageChange={value => setPage(value)}
              />
            ) : null}
          </div>
        </>
      )}
      <CreateBucket
        open={openCreate}
        closeHandler={() => setOpenCreate(false)}
        updateBuckets={refetch}
      />
    </div>
  )
}
