import { useEffect, useMemo, useState } from 'react'

import ErrorState from '@/components/error-state'
import type { IItems } from '@/components/search-select/search.types'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle
} from '@/components/shadcn/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/shadcn/select'
import SkeletonGroup from '@/components/skeleton-group'
import { useBaremetalQuota } from '@/hooks/usage/useBaremetals'
import { useTenants } from '@/hooks/usage/useTenants'
import { t } from '@/i18n'

import BaremetalTable from './baremetal-table'

export default function UsageBaremetals() {
  const [selectedTenant, setSelectedTenant] = useState<string | null>(null)

  const { data: fetchedTenants, isLoading, isError } = useTenants()
  const {
    data: fetchedBaremetalQuota,
    isLoading: isLoadingBaremetal,
    isError: isErrorBaremetal
  } = useBaremetalQuota({
    venture: selectedTenant!
  })

  const tenants: IItems[] = useMemo(() => {
    if (!fetchedTenants?.data) return []

    return fetchedTenants.data.map(item => ({
      title: item,
      value: item
    }))
  }, [fetchedTenants])

  useEffect(() => {
    if (fetchedTenants?.data?.[0] && !selectedTenant) {
      setSelectedTenant(fetchedTenants.data[0])
    }
  }, [fetchedTenants?.data, selectedTenant])

  const renderContent = () => {
    if (isLoadingBaremetal) {
      return <SkeletonGroup count={1} height={300} orientation="horizontal" />
    }

    if (isErrorBaremetal) {
      return <ErrorState />
    }

    if (fetchedBaremetalQuota?.data) {
      return <BaremetalTable items={fetchedBaremetalQuota.data} />
    }

    return null
  }

  return (
    <Card className="mb-10">
      <CardHeader>
        <CardTitle className="flex justify-between">
          <span className="text-xl">{t('baremetals_quota')}</span>
          <div>
            <Select
              value={selectedTenant!}
              disabled={isLoading || isError}
              onValueChange={value => setSelectedTenant(value)}
            >
              <SelectTrigger className="h-9 w-[200px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {tenants.map(item => (
                  <SelectItem key={item.value} value={item.value}>
                    {item.title}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent>{renderContent()}</CardContent>
    </Card>
  )
}
