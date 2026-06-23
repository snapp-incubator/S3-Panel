import ResourceCostChart from '@/components/chart-resource'
import SkeletonGroup from '@/components/skeleton-group'
import { t } from '@/i18n'

import type { ResourceChartsProps } from './resource.types'

export const QuotaCharts = ({
  isVM,
  teamFilter,
  resourceData,
  hasError,
  isLoading
}: ResourceChartsProps) => {
  if (isLoading) {
    return <SkeletonGroup count={1} height={130} orientation="horizontal" />
  }

  return (
    <div className="flex flex-col gap-4">
      <ResourceCostChart
        team={teamFilter}
        data={resourceData.cpu}
        labelData="Core"
        requestHard={true}
        actualUsage={true}
        header={
          <div className="flex w-full items-center justify-between">
            <span className="text-base font-semibold">{t('cpu_usage')}</span>
          </div>
        }
        hasError={hasError}
      />
      <ResourceCostChart
        team={teamFilter}
        data={resourceData.memory}
        labelData="GB"
        requestHard={true}
        actualUsage={true}
        header={
          <div className="flex w-full items-center justify-between">
            <span className="text-base font-semibold">{t('memory_usage')}</span>
          </div>
        }
        hasError={hasError}
      />
      {isVM == true && resourceData.storage ? (
        <ResourceCostChart
          team={teamFilter}
          data={resourceData.storage}
          labelData="GB"
          requestHard={true}
          actualUsage={true}
          header={
            <div className="flex w-full items-center justify-between">
              <span className="text-base font-semibold">
                {t('storage_usage')}
              </span>
            </div>
          }
          hasError={hasError}
        />
      ) : null}
    </div>
  )
}
