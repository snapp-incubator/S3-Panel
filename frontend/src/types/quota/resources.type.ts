export type TResource = 'cpu' | 'memory' | 'storage'

export type TVirtualMachineResource = 'cpu' | 'memory'

export interface TotalUsage {
  memory_value: number
  cpu_request: number
  cpu_limits: number
  pod_counts: number
  storage_value?: number
  ephemeral_storage_value?: number
}

export type QuotaParams = {
  resourceType: TResource | 'all'
  region?: string
  project?: string
  namespace?: string
  duration?: string
  team?: string
  cluster?: string
}

export type NamespaceParams = {
  resourceType: TResource | 'all'
  region?: string
  namespace?: string
  namespaces?: string
  duration?: string
  team?: string
}
