import { useMemo, useState } from 'react'

import { Info } from 'lucide-react'

import ErrorState from '@/components/error-state'
import { Alert, AlertDescription, AlertTitle } from '@/components/shadcn/alert'
import { Button } from '@/components/shadcn/button'
import {
  Card,
  CardTitle,
  CardContent,
  CardHeader
} from '@/components/shadcn/card'
import SkeletonGroup from '@/components/skeleton-group'
import { useBillingTrends } from '@/hooks/billing/useBilling'
import { useFilterStore } from '@/hooks/useFilterStore'
import { t } from '@/i18n'

import CostBarChart from './cost-bar-chart'

type TFilterValues = 'all' | 'okd' | 'openstack'

interface FilterOption {
  label: string
  value: TFilterValues
}

const FILTER_OPTIONS: FilterOption[] = [
  { label: 'All', value: 'all' },
  { label: 'OKD', value: 'okd' },
  { label: 'Openstack', value: 'openstack' }
]

function FilterButtons({
  options,
  activeFilter,
  onChange
}: {
  options: FilterOption[]
  activeFilter: TFilterValues
  onChange: (value: TFilterValues) => void
}) {
  return (
    <div className="flex items-center gap-4">
      {options.map(({ label, value }) => (
        <Button
          key={value}
          variant={activeFilter === value ? 'default' : 'secondary'}
          onClick={() => onChange(value)}
        >
          {label}
        </Button>
      ))}
    </div>
  )
}

export default function CostTrends() {
  const teamFilter = useFilterStore(state => state.filterValues.teamFilter)
  const timeRange = useFilterStore(state => state.filterValues.timeRange)
  const region = useFilterStore(state => state.filterValues.region)

  const { data, isError, isLoading } = useBillingTrends(
    teamFilter,
    timeRange,
    region
  )

  const trends = useMemo(() => data?.data.trends ?? [], [data])

  const [filter, setFilter] = useState<TFilterValues>('all')

  let content = null

  if (isLoading) {
    content = <SkeletonGroup count={1} height={300} orientation="horizontal" />
  } else if (isError) {
    content = <ErrorState />
  } else if (trends.length === 0) {
    content = (
      <Alert className="flex items-start gap-3">
        <Info className="mt-0.5 size-5 text-blue-500" />
        <div>
          <AlertTitle>{t('data_empty_204')}</AlertTitle>
          <AlertDescription>{t('cost_trend_no_data')}</AlertDescription>
        </div>
      </Alert>
    )
  } else {
    content = <CostBarChart trends={trends} filter={filter} />
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <h2 className="text-xl font-semibold">{t('cost_trends')}</h2>
          {!isLoading && !isError && (
            <FilterButtons
              options={FILTER_OPTIONS}
              activeFilter={filter}
              onChange={setFilter}
            />
          )}
        </CardTitle>
      </CardHeader>
      <CardContent className="relative w-full">{content}</CardContent>
    </Card>
  )
}
