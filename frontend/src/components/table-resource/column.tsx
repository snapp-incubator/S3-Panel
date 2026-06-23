import { useNavigate } from '@tanstack/react-router'
import { ColumnDef } from '@tanstack/react-table'
import { ArrowUpDown, TrendingUp } from 'lucide-react'

import { Badge } from '@/components/shadcn/badge'
import { Button } from '@/components/shadcn/button'
import { t } from '@/i18n'
import { formatValue } from '@/lib/helper'
import type { TProject } from '@/types/quota/teams.type'

import { Namespace } from './resource.types'

type TColumn = {
  type: 'namespace' | 'project' | 'instance'
  showLinkButton?: boolean
}

export const useColumns = ({ type, showLinkButton = true }: TColumn) => {
  const navigate = useNavigate()

  const handleClickNamespace = (value: string) => {
    navigate({
      to: '/usage/namespace',
      search: { ns: value }
    })
  }

  const handleClickProject = (value: string, cluster: string) => {
    navigate({
      to: '/usage/project',
      search: { pr: value, cluster: cluster }
    })
  }

  switch (type) {
    case 'project':
      return [
        {
          accessorKey: 'name',
          header: t('vms_projects'),
          cell: ({ row }) =>
            showLinkButton ? (
              <Button
                variant="link"
                onClick={() =>
                  handleClickProject(
                    row.getValue('name'),
                    row.getValue('cluster')
                  )
                }
                className="h-auto p-0 capitalize"
              >
                <TrendingUp className="mr-2 size-4" />
                <span className="truncate text-center text-blue-500 underline">
                  {row.getValue('name')}
                </span>
              </Button>
            ) : (
              <div className="truncate">{row.getValue('name')}</div>
            )
        },
        {
          accessorKey: 'cluster',
          header: () => <div className="text-center">{t('cluster')}</div>,
          cell: ({ row }) => (
            <div className="text-center">
              <Badge variant="secondary">
                {String(row.getValue('cluster')).toUpperCase()}
              </Badge>
            </div>
          )
        },
        {
          accessorKey: 'cpu',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[100px]"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('cpu_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const cpu = parseFloat(row.getValue('cpu'))

            return (
              <div className="w-[100px] text-center font-medium">
                {isNaN(cpu) ? '0' : cpu.toFixed(2)}
              </div>
            )
          }
        },
        {
          accessorKey: 'memory',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[120px] text-center"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('memory_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const memory = parseFloat(row.getValue('memory'))
            const unit = row.original.memory_unit || 'GB'

            return (
              <div className="w-[120px] truncate text-center font-medium">
                {isNaN(memory) ? '0' : `${formatValue(memory)} ${unit}`}
              </div>
            )
          }
        },
        {
          accessorKey: 'storage',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[120px] text-center"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('storage_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const memory = parseFloat(row.getValue('storage'))
            const unit = row.original.storage_unit || 'GB'

            return (
              <div className="w-[120px] truncate text-center font-medium">
                {isNaN(memory) ? '0' : `${formatValue(memory)} ${unit}`}
              </div>
            )
          }
        }
      ] as ColumnDef<TProject>[]
    case 'instance':
      return [
        {
          accessorKey: 'name',
          header: t('vms_projects'),
          cell: ({ row }) =>
            showLinkButton ? (
              <Button
                variant="link"
                onClick={() =>
                  handleClickProject(
                    row.getValue('name'),
                    row.getValue('cluster')
                  )
                }
                className="h-auto p-0 capitalize"
              >
                <TrendingUp className="mr-2 size-4" />
                <span className="truncate text-center text-blue-500 underline">
                  {row.getValue('name')}
                </span>
              </Button>
            ) : (
              <div className="truncate">{row.getValue('name')}</div>
            )
        },
        {
          accessorKey: 'flavor',
          header: () => <div className="text-center">{t('flavor')}</div>,
          cell: ({ row }) => (
            <div className="text-center">
              <Badge variant="secondary">
                {String(row.getValue('flavor'))}
              </Badge>
            </div>
          )
        },
        {
          accessorKey: 'cluster',
          header: () => <div className="text-center">{t('cluster')}</div>,
          cell: ({ row }) => (
            <div className="text-center">
              <Badge variant="secondary">
                {String(row.getValue('cluster')).toUpperCase()}
              </Badge>
            </div>
          )
        },
        {
          accessorKey: 'cpu',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[100px]"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('cpu_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const cpu = parseFloat(row.getValue('cpu'))

            return (
              <div className="w-[100px] text-center font-medium">
                {isNaN(cpu) ? '0' : cpu.toFixed(2)}
              </div>
            )
          }
        },
        {
          accessorKey: 'memory',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[120px] text-center"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('memory_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const memory = parseFloat(row.getValue('memory'))
            const unit = row.original.memory_unit || 'GB'

            return (
              <div className="w-[120px] truncate text-center font-medium">
                {isNaN(memory) ? '0' : `${formatValue(memory)} ${unit}`}
              </div>
            )
          }
        },
        {
          accessorKey: 'storage',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[120px] text-center"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('storage_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const memory = parseFloat(row.getValue('storage'))
            const unit = row.original.storage_unit || 'GB'

            return (
              <div className="w-[120px] truncate text-center font-medium">
                {isNaN(memory) ? '0' : `${formatValue(memory)} ${unit}`}
              </div>
            )
          }
        }
      ] as ColumnDef<TProject>[]
    case 'namespace':
      return [
        {
          accessorKey: 'namespace',
          header: t('namespaces'),
          cell: ({ row }) => (
            <Button
              variant="link"
              onClick={() => handleClickNamespace(row.getValue('namespace'))}
              className="capitalize"
            >
              <TrendingUp className="mr-2 size-4" />
              <span className="truncate text-center text-blue-500 underline">
                {row.getValue('namespace')}
              </span>
            </Button>
          )
        },
        {
          accessorKey: 'pod_counts',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[80px] text-center text-xs"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('pod_count')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => (
            <div className="w-[80px] text-center">
              {row.getValue('pod_counts') ?? '0'}
            </div>
          )
        },
        {
          accessorKey: 'cpu_value',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[100px] text-xs"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('cpu_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const cpu = parseFloat(row.getValue('cpu_value'))

            return (
              <div className="w-[100px] text-center font-medium">
                {isNaN(cpu) ? '0' : cpu}
              </div>
            )
          }
        },
        {
          accessorKey: 'memory_value',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[120px] text-center text-xs"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('memory_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const memory = parseFloat(row.getValue('memory_value'))

            return (
              <div className="w-[120px] truncate text-center font-medium">
                {isNaN(memory) ? '0' : `${formatValue(memory)} GB`}
              </div>
            )
          }
        },
        {
          accessorKey: 'storage_value',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[90px] text-xs"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('storage_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const storage = parseFloat(row.getValue('storage_value'))

            return (
              <div className="w-[100px] text-center font-medium">
                {isNaN(storage) ? '0' : `${formatValue(storage)} GB`}
              </div>
            )
          }
        },
        {
          accessorKey: 'ephemeral_storage_value',
          header: ({ column }) => (
            <Button
              variant="ghost"
              className="w-[150px] text-xs"
              onClick={() =>
                column.toggleSorting(column.getIsSorted() === 'asc')
              }
            >
              {t('ephemeral_allocated')}
              <ArrowUpDown className="ml-2 size-4" />
            </Button>
          ),
          cell: ({ row }) => {
            const storageEphemeral = parseFloat(
              row.getValue('ephemeral_storage_value')
            )

            return (
              <div className="w-[100px] text-center font-medium">
                {isNaN(storageEphemeral)
                  ? '0'
                  : `${formatValue(storageEphemeral)} GB`}
              </div>
            )
          }
        }
      ] as ColumnDef<Namespace>[]
  }
}
