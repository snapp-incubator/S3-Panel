import { useQuery } from '@tanstack/react-query'

import { billingService } from '@/services/api/billing'
import { billingKeys } from '@/services/api/billing/billingKeys'

export function useAllTeams() {
  return useQuery({
    queryFn: () => billingService.getAllTeams(),
    queryKey: billingKeys.teams()
  })
}
