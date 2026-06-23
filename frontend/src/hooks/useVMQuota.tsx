import { useQuery } from '@tanstack/react-query'

import { quotaVMKeys } from '@/api/quotaKeys'
import { VMQuotaService } from '@/services/api/vm-quota'
import { QuotaParams } from '@/types/quota/resources.type'

export function useVMUsage(params: QuotaParams) {
  return useQuery({
    queryKey: quotaVMKeys.usage(params),
    queryFn: () =>
      VMQuotaService.getVMResourceUsage(params).then(res => res.data),
    enabled: !!params.duration && !!params.team
  })
}

export function useVMHistorical(params: QuotaParams) {
  return useQuery({
    queryKey: quotaVMKeys.historical(params),
    queryFn: () =>
      VMQuotaService.getHistoricalData(params).then(res => res.data)
  })
}

export function useVMCount(team: string, region: string) {
  return useQuery({
    queryKey: quotaVMKeys.count(team, region),
    queryFn: () => VMQuotaService.getProjectCount(team).then(res => res.data)
  })
}

export function useVMQuota(params: QuotaParams) {
  return useQuery({
    queryKey: quotaVMKeys.quota(params),
    queryFn: () =>
      VMQuotaService.getResourceUsage(params).then(res => res.data),
    enabled: !!params.team
  })
}

export function useVMProject(params: QuotaParams) {
  return useQuery({
    queryKey: quotaVMKeys.namespaces(params),
    queryFn: () => VMQuotaService.getProjectUsage(params).then(res => res.data),
    enabled: !!params.team && !!params.duration
  })
}

export function useProjectUsageSummery(params: QuotaParams) {
  return useQuery({
    queryKey: quotaVMKeys.projectQuota(params),
    queryFn: () =>
      VMQuotaService.getProjectUsageSummery(params).then(res => res.data),
    enabled: !!params.project
  })
}

export function useProjectUsageDetails(params: QuotaParams) {
  return useQuery({
    queryKey: quotaVMKeys.projectUsageDetails(params),
    queryFn: () =>
      VMQuotaService.getProjectUsageDetails(params).then(res => res.data),
    enabled: !!params.project
  })
}

export function useProjectQuotaDetails(params: QuotaParams) {
  return useQuery({
    queryKey: quotaVMKeys.projectQuotaDetails(params),
    queryFn: () =>
      VMQuotaService.getProjectQuotaDetails(params).then(res => res.data),
    enabled: !!params.project
  })
}

export function useProjectHistorical(params: QuotaParams) {
  return useQuery({
    queryKey: quotaVMKeys.projectHistorical(params),
    queryFn: () =>
      VMQuotaService.getProjectHistorical(params).then(res => res.data),
    enabled: !!params.project
  })
}
