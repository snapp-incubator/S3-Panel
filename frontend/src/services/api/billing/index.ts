import centralClient from '@/services/http/centralClient'
import { buildQueryString } from '@/services/http/query-client'

interface ITotalCostsResponse {
  okd_compute: number
  okd_storage: number
  openstack_compute: number
  openstack_storage: number
  total_okd: number
  total_os: number
  raw_total: number
  maintenance: number
  total: number
}

interface ITrendsResponse {
  team: string
  month: string
  trends: ITrendItems[]
}

export interface ITrendItems {
  date: string
  okd_compute: number
  okd_storage: number
  openstack_compute: number
  openstack_storage: number
  maintenance: number
  total: number
}

export interface ICostBreakdownResponse {
  team: string
  date: string
  breakdown: IBreakdownItem[]
}

export interface IBreakdownItem {
  service_description: string
  quantity: number
  unit: string
  unit_price: number
  service_price: number
  maintenance: number
  total: number
}

export const billingService = {
  getTotalCost: (params: { team: string; date: string }) =>
    centralClient.get<ITotalCostsResponse>(
      `/billing/api/billing/unified?${buildQueryString(params)}`
    ),
  getAllTeams: () => centralClient.get<string[]>('/billing/api/billing/teams'),
  getTrends: (params: { team: string; month: string }) =>
    centralClient.get<ITrendsResponse>(
      `/billing/api/billing/trends?${buildQueryString(params)}`
    ),
  getCostBreakdown: (params: { team: string; date: string }) =>
    centralClient.get<ICostBreakdownResponse>(
      `/billing/api/billing/cost-breakdown?${buildQueryString(params)}`
    )
}
