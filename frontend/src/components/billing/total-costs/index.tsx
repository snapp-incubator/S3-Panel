import { useMemo } from 'react'

import { Boxes, Calculator, Server } from 'lucide-react'

import SkeletonGroup from '@/components/skeleton-group'
import { useUnified } from '@/hooks/billing/useBilling'
import { useFilterStore } from '@/hooks/useFilterStore'
import { t } from '@/i18n'

import FailedAlert from './failed-alert'
import TotalCostCard from './total-cost-card'

export default function TotalCosts() {
  const teamFilter = useFilterStore(state => state.filterValues.teamFilter)
  const timeRange = useFilterStore(state => state.filterValues.timeRange)
  const region = useFilterStore(state => state.filterValues.region)

  const [team, date] = [teamFilter, timeRange]

  const {
    data: totalCosts,
    isLoading,
    isError,
    refetch
  } = useUnified(team, date, region)

  const totalCostsItem = useMemo(() => {
    if (totalCosts && totalCosts.data) {
      const { total, total_okd, total_os, maintenance, raw_total } =
        totalCosts.data

      return [
        {
          title: t('total_cost'),
          total,
          rawTotal: raw_total,
          maintenance,
          default: true,
          icon: Calculator
        },
        {
          title: t('openstack_cost'),
          total: total_os,
          icon: Server
        },
        {
          title: t('okd_cost'),
          total: total_okd,
          icon: Boxes
        }
      ]
    }

    return []
  }, [totalCosts])

  if (isLoading) {
    return <SkeletonGroup count={3} orientation="horizontal" />
  } else if (isError) {
    return <FailedAlert refetch={refetch} />
  }

  return (
    <div className="grid grid-cols-1 gap-4 xl:grid-cols-3">
      {totalCostsItem.map(item => (
        <TotalCostCard key={item.title} item={item} />
      ))}
    </div>
  )
}
