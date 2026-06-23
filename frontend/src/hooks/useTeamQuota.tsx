import { useQuery } from '@tanstack/react-query'

import { quotaTeamKeys } from '@/api/quotaKeys'
import { teamQuotaService } from '@/services/api/team-quota'
import type { QuotaParams, NamespaceParams } from '@/types/quota/resources.type'

export function useTeamUsage(region: string, duration: string) {
  return useQuery({
    queryKey: quotaTeamKeys.usage(region, duration),
    queryFn: () =>
      teamQuotaService.getAggregatedUsage({ duration }).then(res => res.data),
    enabled: !!region && !!duration
  })
}

export function useTeamQuota(team: string) {
  return useQuery({
    queryKey: quotaTeamKeys.quota(team),
    queryFn: () => teamQuotaService.getQuota({ team }).then(res => res.data),
    enabled: !!team
  })
}

export function useTeamHistorical(params: QuotaParams) {
  return useQuery({
    queryKey: quotaTeamKeys.historical(params),
    queryFn: () => teamQuotaService.getHistorical(params).then(res => res.data)
  })
}

export function useTeamNamespaceUsage(params: QuotaParams) {
  return useQuery({
    queryKey: quotaTeamKeys.namespaces(params),
    queryFn: () =>
      teamQuotaService
        .getNamespaceUsage(params)
        .then(res => res.data.namespaces),
    enabled: !!params.team && !!params.duration
  })
}

export function useNamespaceHistorical(
  params: NamespaceParams,
  enabled = true
) {
  return useQuery({
    queryKey: quotaTeamKeys.namespaceHistorical(params),
    queryFn: () =>
      teamQuotaService.getNameSpacesHistorical(params).then(res => res.data),
    enabled: enabled
  })
}

export function useNamespaceQuota(params: NamespaceParams, enabled = true) {
  return useQuery({
    queryKey: quotaTeamKeys.namespaceQuota(params),
    queryFn: () =>
      teamQuotaService.getnameSpacesQuota(params).then(res => res.data),
    enabled: enabled
  })
}
