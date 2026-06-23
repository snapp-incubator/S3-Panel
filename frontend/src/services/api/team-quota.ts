import type { QuotaParams } from '@/types/quota/resources.type'
import {
  TeamAggregatedResponse,
  TQuotaTeamResponse,
  TNamespacesUsageResponse,
  THistoricalResponse,
  TUsageTotal
} from '@/types/quota/teams.type'

import centralClient from '../http/centralClient'
import { buildQueryString } from '../http/query-client'

export const teamQuotaService = {
  getAggregatedUsage: (params: { duration: string }) =>
    centralClient.get<TeamAggregatedResponse>(
      `/cost/api/teams/usage/aggregated?${buildQueryString({ resourceType: 'all', ...params })}`
    ),

  getQuota: (params: { team: string }) =>
    centralClient.get<TQuotaTeamResponse>(
      `/cost/api/teams/quota?${buildQueryString({ resourceType: 'all', teams: params.team })}`
    ),

  getNamespaceUsage: (params: QuotaParams) =>
    centralClient.get<TNamespacesUsageResponse>(
      `/cost/api/namespaces/usage?${buildQueryString(params)}`
    ),

  getHistorical: (params: QuotaParams) =>
    centralClient.get<THistoricalResponse>(
      `/cost/api/teams/historical?${buildQueryString(params)}`
    ),

  getNameSpacesHistorical: (params: QuotaParams) =>
    centralClient.get<THistoricalResponse>(
      `/cost/api/namespaces/historical?${buildQueryString(params)}`
    ),

  getnameSpacesQuota: (params: QuotaParams) =>
    centralClient.get<TUsageTotal>(
      `/cost/api/namespaces/quota?${buildQueryString(params)}`
    )
}
