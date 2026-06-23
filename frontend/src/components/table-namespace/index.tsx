import { useNavigate } from '@tanstack/react-router'
import { ColumnDef } from '@tanstack/react-table'
import { Download, ArrowUpDown, TrendingUp } from 'lucide-react'

import { Button } from '@/components/shadcn/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle
} from '@/components/shadcn/card'
import { DataTable } from '@/components/table'
import { formatValue } from '@/lib/helper'

import { Namespace, TeamUsageCardProps } from './namespace.types'

const NamespaceTable = ({
  title,
  data,
  showDownloadButton = false
}: TeamUsageCardProps) => {
  const navigate = useNavigate()

  const handleClickNamespace = (value: string) => {
    navigate({
      to: '/usage/namespace',
      search: { ns: value }
    })
  }
  const columns: ColumnDef<Namespace>[] = [
    {
      accessorKey: 'namespace',
      header: 'Namespace',
      cell: ({ row }) => (
        <Button
          variant="link"
          onClick={() => handleClickNamespace(row.getValue('namespace'))}
          className="capitalize"
        >
          <TrendingUp className="mr-2 size-4" />
          <span className="truncate text-center">
            {row.getValue('namespace')}
          </span>
        </Button>
      )
    },
    {
      accessorKey: 'pod_counts',
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            className="w-[80px] text-center"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          >
            Pod Counts
            <ArrowUpDown className="ml-2 size-4" />
          </Button>
        )
      },
      cell: ({ row }) => (
        <div className="w-[80px] text-center">
          {row.getValue('pod_counts') ?? '0'}
        </div>
      )
    },
    {
      accessorKey: 'cpu_value',
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            className="w-[100px]"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          >
            CPU Allocated (Core)
            <ArrowUpDown className="ml-2 size-4" />
          </Button>
        )
      },
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
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            className="w-[120px] text-center"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          >
            Memory Allocated (GB)
            <ArrowUpDown className="ml-2 size-4" />
          </Button>
        )
      },
      cell: ({ row }) => {
        const memory = parseFloat(row.getValue('memory_value'))

        return (
          <div className="w-[120px] truncate text-center font-medium">
            {isNaN(memory) ? '0' : `${formatValue(memory)} GB`}
          </div>
        )
      },
      size: 20
    },
    {
      accessorKey: 'storage_value',
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            className="w-[90px]"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          >
            Storage Allocated (GB)
            <ArrowUpDown className="ml-2 size-4" />
          </Button>
        )
      },
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
      header: ({ column }) => {
        return (
          <Button
            variant="ghost"
            className="w-[150px]"
            onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
          >
            Ephemeral Storage Allocated (GB)
            <ArrowUpDown className="ml-2 size-4" />
          </Button>
        )
      },
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
  ]

  return (
    <Card className="col-span-4">
      <CardHeader>
        <div className="flex items-center justify-between space-y-2">
          <CardTitle>{title}</CardTitle>
          {showDownloadButton && (
            <div className="flex space-x-1">
              <Button variant="outline">
                <Download />
              </Button>
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className="pl-2">
        <DataTable data={data} columns={columns} filterColumn="namespace" />
      </CardContent>
    </Card>
  )
}

export default NamespaceTable
