import { Fragment } from 'react'

import { QuotaCharts } from '@/components/quota-charts'
import { ResourceCards } from '@/components/resource-cards'
import { CardType } from '@/components/resource-cards/types'
import SkeletonGroup from '@/components/skeleton-group'
import ResourceTable from '@/components/table-resource'
import { useVMUsage, useVMHistorical, useVMProject } from '@/hooks/useVMQuota'
import { t } from '@/i18n'
import { TUsageVirtualMachine } from '@/types/quota/teams.type'

type VMQuotaProps = {
  region: string
  timeRange: string
  teamFilter: string
}

export const VMQuota = ({ region, timeRange, teamFilter }: VMQuotaProps) => {
  const params = {
    region: region,
    resourceType: 'all' as const,
    duration: timeRange,
    team: teamFilter
  }

  const { data: vmQuotaData, isLoading: isQuotaLoading } = useVMUsage(params)
  const { data: vmProject, isLoading: isProjectLoading } = useVMProject(params)

  const { data: cpuData, isLoading: isCpuLoading } = useVMHistorical({
    ...params,
    resourceType: 'cpu'
  })
  const { data: memoryData, isLoading: isMemoryLoading } = useVMHistorical({
    ...params,
    resourceType: 'memory'
  })

  const isLoading =
    isProjectLoading ||
    isQuotaLoading ||
    isProjectLoading ||
    isCpuLoading ||
    isMemoryLoading
  const hasError = !vmQuotaData

  const totalUsage: TUsageVirtualMachine = {
    total_cpu: vmQuotaData?.cpu || 0,
    total_memory: vmQuotaData?.memory || 0,
    project_count: vmQuotaData?.projects || 0,
    total_memory_unit: vmQuotaData?.memory_unit || '',
    storage: vmQuotaData?.storage || 0,
    storage_unit: vmQuotaData?.storage_unit || ''
  }

  const resourceData = {
    cpu: cpuData,
    memory: memoryData
  }

  if (isLoading) {
    return <SkeletonGroup count={3} height={130} orientation="horizontal" />
  }

  return (
    <Fragment>
      <ResourceCards
        resourceType={CardType.VM}
        totalUsage={totalUsage}
        hasError={hasError}
        isLoading={isLoading}
      />

      {resourceData ? (
        <QuotaCharts
          isVM={true}
          teamFilter={teamFilter}
          resourceData={resourceData}
          hasError={hasError}
          isLoading={isLoading}
        />
      ) : null}
      {vmProject ? (
        <ResourceTable title={t('projects')} data={vmProject} isVM={true} />
      ) : null}
    </Fragment>
  )
}
