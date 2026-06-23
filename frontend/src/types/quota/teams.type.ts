import type { TResource } from '@/types/quota/resources.type'

export type TUsageTotal = {
  memory_value: number
  cpu_request: number
  cpu_limits: number
  pod_counts: number
  storage_value?: number
  ephemeral_storage_value?: number
}

export type TUsageVirtualMachine = {
  total_cpu: number
  total_memory: number
  project_count: number
  total_memory_unit: string
  storage: number
  storage_unit: string
}

export type TQuotaParams = {
  resourceType: TResource | 'all'
  duration: string
  team: string
}

export type THistoricalResponse = {
  actual_usage: {
    timestamp: string
    value: number
  }[]
  request_allocated: {
    timestamp: string
    value: number
  }[]
}

export type TProject = {
  cpu: number
  memory_unit: string
  memory: number
  name: string
  flavor?: string
  storage: number
  storage_unit: string
  cluster: string
}

export type TAggregatedTeam = {
  cpu_limit_value: number
  cpu_request_value: number
  ephemeral_storage_limit_unit: string
  ephemeral_storage_limit_value: number
  ephemeral_storage_request_unit: string
  ephemeral_storage_request_value: number
  memory_limit_unit: string
  memory_limit_value: number
  memory_request_unit: string
  memory_request_value: number
  pod_counts: number
  storage_limit_unit: string
  storage_limit_value: number
  storage_request_unit: string
  storage_request_value: number
  team: string
}

export type TQuotaTeam = {
  cpu_limit_value: number
  cpu_request_value: number
  ephemeral_storage_limit_unit: string
  ephemeral_storage_limit_value: number
  ephemeral_storage_request_unit: string
  ephemeral_storage_request_value: number
  memory_limit_unit: string
  memory_limit_value: number
  memory_request_unit: string
  memory_request_value: number
  pod_counts: number
  storage_limit_unit: string
  storage_limit_value: number
  storage_request_unit: string
  storage_request_value: number
  team: string
}

export type TeamAggregatedResponse = {
  responses: TAggregatedTeam[]
}

export type TQuotaTeamResponse = {
  team: string
  cpu_request: number
  cpu_limits: number
  memory_value: number
  pod_counts: number
  storage_value: number
  ephemeral_storage_value: number
}

export type THistoricalTeam = {
  timestamp: string
  value: number
}

export type TNamespaces = {
  cpu_value: number
  ephemeral_storage_unit: string
  ephemeral_storage_value: number
  memory_unit: string
  memory_value: number
  namespace: string
  pod_counts: number
  storage_unit: string
  storage_value: number
}

export type TNamespacesUsageResponse = {
  namespaces: TNamespaces[]
  team: string
  total_cpu: number
  total_ephemeral_storage: number
  total_memory: number
  total_storage: number
}
