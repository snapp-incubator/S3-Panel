import type { THistoricalResponse } from '@/types/quota/teams.type'

export type ResourceData = {
  cpu?: THistoricalResponse
  memory?: THistoricalResponse
  storage?: THistoricalResponse
}

export type ResourceChartsProps = {
  isVM: boolean
  teamFilter: string
  resourceData: ResourceData
  hasError: boolean
  isLoading: boolean
}
