import type { TotalUsage } from '@/types/quota/resources.type'
import type { TUsageVirtualMachine } from '@/types/quota/teams.type'

export type totalUsageType = {
  memory_value: number
  cpu_request: number
  pod_counts: number
  cpu_limits: number
  storage_value: number
  ephemeral_storage_value: number
}

export const enum CardType {
  OKD = 'okd',
  VM = 'vm',
  NAMESPACE = 'namespace',
  PROJECT = 'project'
}

type ResourceType = `${CardType}`

export type ResourceCardsProps =
  | {
      resourceType: ResourceType
      totalUsage: TUsageVirtualMachine
      hasError: boolean
      isLoading: boolean
    }
  | {
      resourceType: ResourceType
      totalUsage: TotalUsage
      hasError: boolean
      isLoading: boolean
    }
