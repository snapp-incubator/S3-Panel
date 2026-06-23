import centralClient from '@/services/http/centralClient'
import { buildQueryString } from '@/services/http/query-client'

export interface IBareMetalQuota {
  ID: number
  TeamName: string
  Venture: string
  ServerName: string
  Date: string
  CPU: number
  Memory: number
  StorageSSD: number
  StorageHDD: number
  ServerModel: string
  CreatedAt: string
}

export const baremetalsService = {
  getTenants: () => centralClient.get<string[]>('/cost/api/tenants'),
  getBaremetalQuota: (params: { venture: string }) =>
    centralClient.get<IBareMetalQuota[]>(
      `/cost/api/baremetal/quota?${buildQueryString(params)}`
    )
}
