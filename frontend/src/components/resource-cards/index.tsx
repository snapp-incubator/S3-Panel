import SkeletonGroup from '@/components/skeleton-group'
import TotalCostCard from '@/components/widget-cost'
import CPUCostCard from '@/components/widget-cpu'
import { t } from '@/i18n'
import { formatValue } from '@/lib/helper'
import type { TotalUsage } from '@/types/quota/resources.type'
import type { TUsageVirtualMachine } from '@/types/quota/teams.type'

import { ResourceCardsProps, CardType } from './types'

export const ResourceCards = ({
  resourceType,
  totalUsage,
  hasError,
  isLoading
}: ResourceCardsProps) => {
  let content: React.ReactNode

  if (isLoading) {
    return (
      <SkeletonGroup
        count={resourceType == CardType.PROJECT ? 3 : 5}
        height={130}
        orientation="horizontal"
      />
    )
  }

  const renderProjectCards = (
    usage: TUsageVirtualMachine,
    hasError: boolean
  ) => (
    <div className="mb-5 grid gap-4 md:grid-cols-2">
      <TotalCostCard
        title={t('cpu_quota')}
        unit="Core"
        amount={usage?.total_cpu}
        hasError={hasError}
      />
      <TotalCostCard
        title={t('memory_quota')}
        unit={usage?.total_memory_unit}
        amount={formatValue(usage?.total_memory)}
        hasError={hasError}
      />
    </div>
  )

  const renderDefaultCards = (usage: TotalUsage, hasError: boolean) => (
    <div className="mb-5 grid gap-4 md:grid-cols-5 ">
      <TotalCostCard
        title={t('total_pods')}
        amount={usage?.pod_counts}
        hasError={hasError}
      />
      <TotalCostCard
        title={t('memory_quota')}
        unit="GB"
        amount={formatValue(usage?.memory_value)}
        tooltipContent={t('memory_quota_team_help')}
        hasError={hasError}
      />
      <CPUCostCard
        title={t('cpu_quota')}
        tooltipContent={t('cpu_quota_team_help')}
        unit="Core"
        amount={{
          request: {
            title: t('cpu_request'),
            value: usage?.cpu_request,
            help: t('cpu_request_help')
          },
          limit: {
            title: t('cpu_limit'),
            value: usage?.cpu_limits,
            help: t('cpu_limit_help')
          }
        }}
        hasError={hasError}
      />
      <TotalCostCard
        title={t('storage_quota')}
        unit="GB"
        amount={formatValue(usage?.storage_value || 0)}
        tooltipContent={t('storage_quota_team_help')}
        hasError={hasError}
      />
      <TotalCostCard
        title={t('ephemeral_quota')}
        unit="GB"
        amount={formatValue(usage?.ephemeral_storage_value || 0)}
        tooltipContent={t('ephemeral_storage_quota_team_help')}
        hasError={hasError}
      />
    </div>
  )

  const renderVMCards = (usage: TUsageVirtualMachine, hasError: boolean) => (
    <div className="mb-5 grid gap-4 md:grid-cols-4">
      <TotalCostCard
        title={t('total_project')}
        amount={usage?.project_count}
        hasError={hasError}
      />
      <TotalCostCard
        title={t('cpu_quota')}
        amount={formatValue(usage?.total_cpu)}
        unit="Core"
        hasError={hasError}
      />
      <TotalCostCard
        title={t('memory_quota')}
        unit={usage.total_memory_unit}
        amount={formatValue(usage?.total_memory)}
        hasError={hasError}
      />
      <TotalCostCard
        title={t('storage_quota')}
        unit={usage.storage_unit}
        amount={formatValue(usage?.storage)}
        hasError={hasError}
      />
    </div>
  )

  switch (resourceType) {
    case CardType.NAMESPACE:
      content = renderDefaultCards(totalUsage as TotalUsage, hasError)
      break
    case CardType.OKD:
      content = renderDefaultCards(totalUsage as TotalUsage, hasError)
      break
    case CardType.VM:
      content = renderVMCards(totalUsage as TUsageVirtualMachine, hasError)
      break
    case CardType.PROJECT:
      content = renderProjectCards(totalUsage as TUsageVirtualMachine, hasError)
      break
  }

  return <>{content}</>
}
