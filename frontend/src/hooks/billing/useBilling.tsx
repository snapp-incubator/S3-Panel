import { useQuery } from '@tanstack/react-query'

import { billingService } from '@/services/api/billing'
import { billingKeys } from '@/services/api/billing/billingKeys'

export const useUnified = (team: string, date: string, region: string) => {
  const { getTotalCost } = billingService

  return useQuery({
    queryKey: billingKeys.totalCost(team, date, region),
    queryFn: () =>
      getTotalCost({
        team: team,
        date
      }),
    enabled: !!team && !!date
  })
}

export const useBillingTrends = (
  team: string,
  date: string,
  region: string
) => {
  const { getTrends } = billingService

  return useQuery({
    queryKey: billingKeys.trends(team, date, region),
    queryFn: () => getTrends({ team, month: date })
  })
}

export const useCostBreakdown = (
  team: string,
  date: string,
  region: string
) => {
  const { getCostBreakdown } = billingService

  return useQuery({
    queryKey: billingKeys.costBreakdown(team, date, region),
    queryFn: () => getCostBreakdown({ team, date })
  })
}
