import type { QuotaParams, NamespaceParams } from '@/types/quota/resources.type'

import { createKeyStore } from './query-keys'

export const quotaTeamKeys = createKeyStore('teamQuota', {
  usage: (region: string, duration: string) => ['usage', region, duration],
  quota: (team: string) => ['quota', team],
  historical: (params: QuotaParams) => [
    'historical',
    params.team,
    params.region,
    params.duration,
    params.resourceType
  ],
  namespaces: (params: QuotaParams) => [
    'namespaces',
    params.team,
    params.region,
    params.duration,
    params.resourceType
  ],
  namespaceHistorical: (params: NamespaceParams) => [
    'namespaceHistorical',
    params.namespace
  ],
  namespaceQuota: (params: NamespaceParams) => [
    'namespaceQuota',
    params.namespace,
    params.namespaces
  ]
})

export const quotaVMKeys = createKeyStore('teamVMQuota', {
  usage: (params: QuotaParams) => [
    'usage',
    params.region,
    params.team,
    params.duration,
    params.resourceType
  ],
  quota: (params: QuotaParams) => [
    'quota',
    params.team,
    params.region,
    params.duration,
    params.resourceType
  ],
  count: (team: string, region: string) => ['count', team, region],
  historical: (params: QuotaParams) => [
    'historical',
    params.region,
    params.team,
    params.duration,
    params.resourceType
  ],
  namespaces: (params: QuotaParams) => [
    'namespaces',
    params.team,
    params.region,
    params.duration,
    params.resourceType
  ],
  namespaceHistorical: (params: NamespaceParams) => [
    'namespaceHistorical',
    params.team,
    params.region,
    params.duration,
    params.resourceType
  ],
  projectHistorical: (params: QuotaParams) => [
    'projectHistorical',
    params.project
  ],
  namespaceQuota: (params: QuotaParams) => [
    'namespaceQuota',
    params.project,
    params.duration,
    params.resourceType
  ],
  projectQuota: (params: QuotaParams) => [
    'projectQuota',
    params.project,
    params.duration,
    params.cluster || '',
    params.resourceType
  ],
  projectUsageDetails: (params: QuotaParams) => [
    'projectUsageDetails',
    params.project,
    params.cluster
  ],
  projectQuotaDetails: (params: QuotaParams) => [
    'projectQuotaDetails',
    params.project,
    params.cluster
  ]
})
