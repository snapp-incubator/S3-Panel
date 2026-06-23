import { useState } from 'react'

import { useNavigate } from '@tanstack/react-router'
import { Copy, Ellipsis } from 'lucide-react'

import { Button } from '@/components/shadcn/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '@/components/shadcn/dropdown-menu'
import { useToast } from '@/hooks/use-toast'
import useDeleteBucket from '@/hooks/useDeleteBucket'
import { t } from '@/i18n'
import { copyString, dateFormat } from '@/lib/utils'
import type { TBucketResponse } from '@/types/s3/buckets.types'

import DeleteConfirmation from '../delete-confirmation'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '../shadcn/tooltip'

import { TUsageInfoProps } from './bucketCard.types'
import DeleteBucket from './delete-bucket'
import ShowUsage from './show-usage'

const BucketCard = ({
  bucket,
  modify_time_stamp,
  access,
  hard_bytes,
  used_bytes,
  hard_objects,
  used_objects,
  quota_enabled,
  hard_bytes_unit,
  used_bytes_unit,
  used_bytes_raw,
  hard_bytes_raw
}: TBucketResponse) => {
  return (
    <div
      className="flex h-[210px] w-full min-w-[365px]"
      data-test="bucket-card"
      data-test-bucket-name={bucket}
    >
      <BucketHeader bucket={bucket} date={modify_time_stamp} access={access} />
      <div className="flex w-[55%] flex-col gap-3 p-2">
        <UsageInfo
          quotaEnabled={quota_enabled}
          hardData={hard_bytes}
          usedData={used_bytes}
          objectQuota={hard_objects}
          usedObjects={used_objects}
          hardUnit={hard_bytes_unit}
          usedUnit={used_bytes_unit}
          usedRaw={used_bytes_raw}
          hardRaw={hard_bytes_raw}
          bucket={bucket}
        />
        <ActionButtons bucket={bucket} />
      </div>
    </div>
  )
}

export default BucketCard

const BucketHeader = ({
  bucket,
  date,
  access
}: {
  bucket: string
  date: string
  access: string
}) => {
  const { toast } = useToast()

  const formatedDate = dateFormat(date)

  const copyName = (name: string) => {
    copyString(name).then(key => {
      toast({
        title: `${key} copied`,
        variant: 'info'
      })
    })
  }

  return (
    <div className="flex w-[45%] flex-col justify-end rounded-lg bg-gradient-to-b from-green-400 to-black px-4 py-6">
      <div className="flex flex-col text-white">
        <div className="flex items-start justify-between gap-2">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <h3 className="line-clamp-3 text-xl underline">{bucket}</h3>
              </TooltipTrigger>
              <TooltipContent>
                <p>{bucket}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <Button
            variant="link"
            size="icon"
            className="shrink-0"
            onClick={() => copyName(bucket)}
          >
            <Copy className="size-5 text-white" />
          </Button>
        </div>

        <span className="mt-2 text-sm">{`${t('modification_date')}: ${formatedDate}`}</span>
        <span className="text-sm">{`${t('access')}: ${access}`}</span>
      </div>
    </div>
  )
}

const BucketMoreDetails = ({ bucket }: { bucket: string }) => {
  const [openDeleteConfirmation, setOpenDeleteConfirmation] = useState(false)
  const [openMenu, setOpenMenu] = useState(false)

  const { mutate: deleteBucket, isPending } = useDeleteBucket()

  return (
    <div data-test="bucket-more-details">
      <DropdownMenu
        open={openMenu && !openDeleteConfirmation}
        onOpenChange={open => setOpenMenu(open)}
      >
        <DropdownMenuTrigger asChild>
          <Button
            variant="ghost"
            size="sm"
            data-test="bucket-more-details-trigger"
          >
            <Ellipsis />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onSelect={e => e.preventDefault()} className="p-0">
            <DeleteBucket
              disable={isPending}
              clickHandler={() => {
                setOpenDeleteConfirmation(true)
                setOpenMenu(false)
              }}
            />
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <DeleteConfirmation
        open={openDeleteConfirmation}
        closeHandler={() => setOpenDeleteConfirmation(false)}
        acceptDelete={() => deleteBucket(bucket)}
        deleteItemName={bucket}
      ></DeleteConfirmation>
    </div>
  )
}

const UsageInfo = ({
  quotaEnabled,
  hardData,
  hardUnit,
  usedUnit,
  usedData,
  objectQuota,
  usedObjects,
  bucket,
  usedRaw,
  hardRaw
}: TUsageInfoProps) => {
  return (
    <>
      <div className="flex flex-col gap-1">
        <div className="flex items-center justify-between">
          <span className="text-base font-bold">{`${t('bucket_size')}:`}</span>
          <BucketMoreDetails bucket={bucket} />
        </div>
        <ShowUsage
          quotaEnabled={quotaEnabled}
          hardData={hardData}
          usedData={usedData}
          hardUnit={hardUnit}
          usedUnit={usedUnit}
          usedRaw={usedRaw}
          hardRaw={hardRaw}
          withUnit={true}
        />
      </div>
      <div className="flex flex-col gap-1">
        <span className="text-base font-bold">{`${t('object_count')}:`}</span>
        <ShowUsage
          quotaEnabled={quotaEnabled}
          hardData={objectQuota}
          usedData={usedObjects}
        />
      </div>
    </>
  )
}

const ActionButtons = ({ bucket }: { bucket: string }) => {
  const navigate = useNavigate()

  return (
    <div className="mt-auto flex items-center gap-2">
      <Button
        size="sm"
        className="w-full"
        data-test="browse-bucket-button"
        onClick={() =>
          navigate({
            to: `/object-storage/s3-bucket/buckets/${bucket}`
          })
        }
      >
        Browse
      </Button>
    </div>
  )
}
