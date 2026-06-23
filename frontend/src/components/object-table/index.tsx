import { AlertCircle, Share2 } from 'lucide-react'

import { Alert, AlertDescription, AlertTitle } from '@/components/shadcn/alert'
import { Button } from '@/components/shadcn/button'
import { Skeleton } from '@/components/shadcn/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/shadcn/table'
import { t } from '@/i18n'
import { dateFormat } from '@/lib/utils'

import DeleteObject from './delete-object'
import DownloadObject from './download-object'
import type { TObjectTablesProps } from './objectTables.types'

const ObjectSkeleton = () => {
  return new Array(6).fill('').map((_, index) => {
    return (
      <TableRow key={`item_${index}`}>
        <TableCell>
          <Skeleton className="h-4 w-full" />
        </TableCell>
        <TableCell>
          <Skeleton className="h-4 w-full" />
        </TableCell>
        <TableCell>
          <Skeleton className="h-4 w-full" />
        </TableCell>
        <TableCell>
          <Skeleton className="h-4 w-full" />
        </TableCell>
      </TableRow>
    )
  })
}

export default function ObjectTable({
  isError,
  isLoading,
  bucket,
  objectList,
  isSearch,
  refetchObjects,
  onShareObject
}: TObjectTablesProps) {
  return (
    <div>
      {isError || (!isLoading && !objectList?.items) ? (
        <div className="p-4">
          <Alert>
            <AlertCircle className="size-5 !text-blue-400" />
            <AlertTitle className="!text-blue-400">{t('no_data')}</AlertTitle>
            {isSearch ? null : (
              <AlertDescription>{t('upload_first_object')}</AlertDescription>
            )}
          </Alert>
        </div>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t('name')}</TableHead>
              <TableHead className="text-center">
                {t('modification_date')}
              </TableHead>
              <TableHead className="text-center">{t('size')}</TableHead>
              <TableHead className="text-center">{t('actions')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <ObjectSkeleton />
            ) : (
              objectList?.items.map(item => (
                <TableRow key={item.name}>
                  <TableCell className="font-medium">{item.name}</TableCell>
                  <TableCell className="text-center">
                    {dateFormat(item.last_modified_timestamp)}
                  </TableCell>
                  <TableCell className="text-center">
                    {`${item.size_value} ${item.size_unit}`}
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center justify-center gap-2">
                      <Button
                        size="icon"
                        variant="ghost"
                        onClick={() => onShareObject(item.name)}
                      >
                        <Share2 />
                      </Button>
                      <DeleteObject
                        object={item.name}
                        bucket={bucket}
                        refetchObjects={refetchObjects}
                      />
                      <DownloadObject bucket={bucket} object={item.name} />
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      )}
    </div>
  )
}
