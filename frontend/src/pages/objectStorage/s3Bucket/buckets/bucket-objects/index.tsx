import { useQuery } from '@tanstack/react-query'
import { useNavigate, useParams } from '@tanstack/react-router'
import { useDeferredValue, useState } from 'react'

import { fetchObjects } from '@/api/s3'
import { bucketObjectKeys } from '@/api/s3Keys'
import CustomPagination from '@/components/custom-pagination'
import ObjectTable from '@/components/object-table'
import { useTitle } from '@/components/providers/titleProvider'
import SearchField from '@/components/search-field'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '@/components/shadcn/breadcrumb'
import { Button } from '@/components/shadcn/button'
import { Card, CardContent, CardFooter } from '@/components/shadcn/card'
import ShareObject from '@/components/share-object'
import UploadObject from '@/components/upload-object'
import useEffectOnce from '@/hooks/useEffectOnce'

export default function BucketObjects() {
  const { setTitle } = useTitle()

  const { bucketName } = useParams({ strict: false })

  const [currentPage, setCurrentPage] = useState(1)

  useEffectOnce(() => {
    setTitle('SnappCloud - S3 Bucket Objects')
  })

  const navigate = useNavigate()

  const [openShareObject, setOpenShareObject] = useState(false)
  const [objectName, setObjectName] = useState<string | null>(null)

  const [searchValue, setSearchValue] = useState('')
  const deferredSearch = useDeferredValue(searchValue)

  const {
    data: objectList,
    isLoading,
    refetch,
    isError
  } = useQuery({
    queryKey: bucketObjectKeys.all(bucketName!, currentPage, deferredSearch),
    queryFn: () => fetchObjects(bucketName!, currentPage, deferredSearch)
  })

  return (
    <div>
      <h2 className="text-3xl">S3 Bucket</h2>
      <Breadcrumb className="mt-4">
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink>
              <Button
                size="sm"
                variant="link"
                onClick={() =>
                  navigate({
                    to: '/object-storage/s3-bucket/buckets'
                  })
                }
              >
                Buckets
              </Button>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Objects</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <div className="mt-6 flex items-center justify-between">
        <SearchField
          value={searchValue}
          onChange={value => {
            setSearchValue(value)
            setCurrentPage(1)
          }}
        />
        <UploadObject
          bucketName={bucketName!}
          refetchObjects={() => refetch()}
        />
      </div>
      <Card className="mt-8">
        <CardContent className="p-0">
          <ObjectTable
            isError={isError}
            isLoading={isLoading}
            objectList={objectList!}
            bucket={bucketName!}
            refetchObjects={() => refetch()}
            isSearch={!!searchValue}
            onShareObject={objectName => {
              setObjectName(objectName)
              setOpenShareObject(true)
            }}
          />
        </CardContent>
        <CardFooter className="relative">
          {objectList?.total_pages && objectList.total_pages > 1 ? (
            <CustomPagination
              currentPage={currentPage}
              totalPages={objectList?.total_pages}
              onPageChange={value => setCurrentPage(value)}
            />
          ) : null}
        </CardFooter>
      </Card>
      <ShareObject
        open={openShareObject}
        bucket={bucketName!}
        object={objectName!}
        closeHandler={() => setOpenShareObject(false)}
      />
    </div>
  )
}
