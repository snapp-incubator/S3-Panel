import { useQuery } from '@tanstack/react-query'

import { baremetalsService } from '@/services/api/usage/baremetals'
import { baremetalKeys } from '@/services/api/usage/baremetals/baremetalKeys'

interface IUseBaremetalQuotaProps {
  venture: string
}

export function useBaremetalQuota({ venture }: IUseBaremetalQuotaProps) {
  const { getBaremetalQuota } = baremetalsService

  return useQuery({
    queryFn: () => getBaremetalQuota({ venture }),
    queryKey: baremetalKeys.baremetalQuota(venture),
    enabled: !!venture
  })
}
