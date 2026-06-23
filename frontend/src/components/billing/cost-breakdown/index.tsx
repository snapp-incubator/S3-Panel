import { Download } from 'lucide-react'

import ErrorState from '@/components/error-state'
import { Button } from '@/components/shadcn/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle
} from '@/components/shadcn/card'
import SkeletonGroup from '@/components/skeleton-group'
import { useCostBreakdown } from '@/hooks/billing/useBilling'
import { useFilterStore } from '@/hooks/useFilterStore'
import { t } from '@/i18n'

import BreakdownTable from './breakdown-table'

export default function CostBreakdown() {
  const teamFilter = useFilterStore(state => state.filterValues.teamFilter)
  const timeRange = useFilterStore(state => state.filterValues.timeRange)
  const region = useFilterStore(state => state.filterValues.region)

  const {
    data: costBreakdown,
    isError,
    isLoading
  } = useCostBreakdown(teamFilter, timeRange, region)

  let content = null

  if (isLoading) {
    content = <SkeletonGroup count={1} height={300} orientation="horizontal" />
  } else if (isError) {
    content = <ErrorState />
  } else if (costBreakdown && costBreakdown.data?.breakdown) {
    content = <BreakdownTable data={costBreakdown.data.breakdown} />
  }

  return (
    <Card className="shadow-sm">
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-xl font-semibold">
          {t('cost_breakdown')}
        </CardTitle>
        <Button variant="outline" disabled size="sm">
          <Download className="mr-2 size-4" /> {t('download_bill')}
        </Button>
      </CardHeader>

      <CardContent>{content}</CardContent>
    </Card>
  )
}
