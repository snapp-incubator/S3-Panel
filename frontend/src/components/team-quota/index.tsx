import { QuotaCharts } from '@/components/quota-charts'
import { ResourceCards } from '@/components/resource-cards'
import { CardType } from '@/components/resource-cards/types'
import SkeletonGroup from '@/components/skeleton-group'
import ResourceTable from '@/components/table-resource'
import {
  useTeamQuota,
  useTeamHistorical,
  useTeamNamespaceUsage
} from '@/hooks/useTeamQuota'
import { t } from '@/i18n'
import { TotalUsage } from '@/types/quota/resources.type'

type TeamQuotaProps = {
  region: string
  timeRange: string
  teamFilter: string
}

export const TeamQuota = ({
  region,
  timeRange,
  teamFilter
}: TeamQuotaProps) => {
  const params = {
    region: region,
    duration: timeRange,
    team: teamFilter
  }

  const {
    data: totalUsage,
    isLoading: isTotalLoading,
    isError: hasTotalError
  } = useTeamQuota(teamFilter)

  const { data: cpuData, isLoading: isCpuLoading } = useTeamHistorical({
    ...params,
    resourceType: 'cpu'
  })

  const { data: memoryData, isLoading: isMemoryLoading } = useTeamHistorical({
    ...params,
    resourceType: 'memory'
  })

  const { data: storageData, isLoading: isStorageLoading } = useTeamHistorical({
    ...params,
    resourceType: 'storage'
  })

  const { data: resources = [], isLoading: isNamespaceLoading } =
    useTeamNamespaceUsage({
      ...params,
      resourceType: 'all'
    })

  const isLoading =
    isTotalLoading ||
    isCpuLoading ||
    isMemoryLoading ||
    isStorageLoading ||
    isNamespaceLoading
  const hasError = hasTotalError

  const resourceData = {
    cpu: cpuData,
    memory: memoryData,
    storage: storageData
  }

  if (isLoading) {
    return <SkeletonGroup count={3} height={130} orientation="horizontal" />
  }

  return (
    <>
      <ResourceCards
        resourceType={CardType.NAMESPACE}
        totalUsage={totalUsage as TotalUsage}
        hasError={hasError}
        isLoading={isLoading}
      />

      <QuotaCharts
        isVM={false}
        teamFilter={teamFilter}
        resourceData={resourceData}
        hasError={hasError}
        isLoading={isLoading}
      />

      <ResourceTable
        title={t('namespace')}
        data={resources.map(ns => ({ ...ns, id: ns.namespace }))}
        isVM={false}
      />
    </>
  )
}
