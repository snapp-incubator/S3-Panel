import { buildQueryString } from '@/services/http/query-client'
import type { QuotaParams } from '@/types/quota/resources.type'
import type { THistoricalResponse, TProject } from '@/types/quota/teams.type'
import type {
  TVMCountResponse,
  TVMTotalResponse,
  TVMUsageResponse
} from '@/types/quota/vm.type'

import centralClient from '../http/centralClient'

export const VMQuotaService = {
  getHistoricalData: (params: QuotaParams) =>
    centralClient.get<THistoricalResponse>(
      `/cost/api/openstack/teams/historical?${buildQueryString(params)}`
    ),

  getResourceUsage: ({ team }: QuotaParams) =>
    centralClient.get<TVMTotalResponse>(
      `/cost/api/teams/vms/total?${buildQueryString({ resourceType: 'all', teams: team })}`
    ),

  getProjectCount: (team?: string) =>
    centralClient.get<TVMCountResponse>(
      `/cost/api/teams/vms/count?teams=${team}`
    ),

  getVMResourceUsage: (params: QuotaParams) =>
    centralClient.get<TVMUsageResponse>(
      `/cost/api/openstack/teams/quota/summery?${buildQueryString({ ...params })}`
    ),

  getProjectUsageSummery: (params: QuotaParams) =>
    centralClient.get<TVMUsageResponse>(
      `/cost/api/openstack/projects/quota/summery?${buildQueryString({ ...params })}`
    ),

  getProjectUsageDetails: (params: QuotaParams) =>
    centralClient.get<TVMUsageResponse[]>(
      `/cost/api/openstack/projects/usage/details?${buildQueryString({ ...params })}`
    ),

  getProjectQuotaDetails: (params: QuotaParams) =>
    centralClient.get<TVMUsageResponse[]>(
      `/cost/api/openstack/projects/quota/details?${buildQueryString({ ...params })}`
    ),

  getProjectUsage: ({ team, duration, resourceType }: QuotaParams) =>
    centralClient.get<TProject[]>(
      `/cost/api/openstack/teams/quota/details?${buildQueryString({ resourceType: resourceType, duration, team: team })}`
    ),

  getProjectHistorical: ({
    project,
    duration,
    resourceType,
    cluster
  }: QuotaParams) =>
    centralClient.get<THistoricalResponse>(
      `/cost/api/openstack/projects/historical?${buildQueryString({ resourceType: resourceType, duration, project, cluster })}`
    )
}
