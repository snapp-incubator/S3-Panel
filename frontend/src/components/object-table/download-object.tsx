import { Download } from 'lucide-react'

import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/shadcn/tooltip'
import useDownloadObject from '@/hooks/useDownloadObject'
import { t } from '@/i18n'
import { buildQueryString } from '@/services/http/query-client'

import { Button } from '../shadcn/button'

import { TDownloadObjectProps } from './downloadObject.types'

export default function DownloadObject({
  bucket,
  object
}: TDownloadObjectProps) {
  const downloadLink = `/s3/api/object/download?${buildQueryString({
    bucket,
    object
  })}`

  const { mutate: downloadObject } = useDownloadObject(downloadLink)

  return (
    <TooltipProvider delayDuration={50}>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button size="icon" variant="ghost" onClick={() => downloadObject()}>
            <Download />
          </Button>
        </TooltipTrigger>
        <TooltipContent>
          <p>{t('download_object')}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
