export type TVMUsage = {
  projects: number
  cpu: number
  memory: number
  memory_unit?: string
  storage: number
  storage_unit: string
  vms: number
}

export type TVMUsageResponse = TVMUsage

export type TVMTotalResponse = {
  team: string
  total_cpu: number
  total_memory: number
  total_memory_unit: string
}[]

export type TVMCount = {
  count: number
  team: string
}

export type TVMCountResponse = {
  responses: TVMCount[]
}
