import type { TProject } from '@/types/quota/teams.type'
import type { TVMUsageResponse } from '@/types/quota/vm.type'

export type Namespace = {
  id: string
  namespace: string
  cpu_value: number
  memory_value: number
  pod_counts: number
  storage_value: number
  ephemeral_storage_value: number
  memory_unit: string
  ephemeral_storage_unit: string
  storage_unit: string
}

export type ResourceTableProps = {
  title: string
  showDownloadButton?: boolean
  showLinkButton?: boolean
  data: Namespace[] | TProject[] | TVMUsageResponse[]
  isVM?: boolean
  isInstance?: boolean
}

export type { ResourceTableProps as TeamUsageCardProps }
