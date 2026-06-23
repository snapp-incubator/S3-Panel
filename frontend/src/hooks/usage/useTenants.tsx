import { useQuery } from '@tanstack/react-query'

import { baremetalsService } from '@/services/api/usage/baremetals'
import { baremetalKeys } from '@/services/api/usage/baremetals/baremetalKeys'

export function useTenants() {
  const { getTenants } = baremetalsService

  return useQuery({
    queryFn: () => getTenants(),
    queryKey: baremetalKeys.tenants()
  })
}
