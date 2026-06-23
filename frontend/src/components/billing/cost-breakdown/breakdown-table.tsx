import {
  useReactTable,
  createColumnHelper,
  flexRender,
  getCoreRowModel
} from '@tanstack/react-table'
import { Info } from 'lucide-react'

import { Button } from '@/components/shadcn/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/shadcn/table'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/shadcn/tooltip'
import { t } from '@/i18n'
import { formatNumber } from '@/lib/utils'

type CostItem = {
  service_description: string
  quantity: number
  unit: string
  unit_price: number
  service_price: number
  maintenance: number
  total: number
  tooltip?: string
}

type BreakdownTableProps = {
  data: CostItem[]
}

const isTotalRow = (desc: string) => desc.toLowerCase().includes('total')

const columnHelper = createColumnHelper<CostItem>()

const columns = [
  columnHelper.accessor('service_description', {
    header: t('service_description'),
    cell: info => {
      const tooltip = info.row.original?.tooltip || null

      const value = info.getValue()

      return isTotalRow(value) ? (
        <span className="font-semibold">{t('total')}</span>
      ) : (
        <div className="flex items-center gap-2">
          {value}
          {tooltip ? (
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger>
                  <Button variant="link" size="icon" className="text-blue-500">
                    <Info size={20} />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  <p>{tooltip}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          ) : null}
        </div>
      )
    }
  }),
  columnHelper.accessor('quantity', {
    header: t('quantity'),
    cell: info => {
      const row = info.row.original

      return isTotalRow(row.service_description)
        ? ''
        : formatNumber(Number(row.quantity))
    }
  }),
  columnHelper.accessor('unit', {
    header: t('unit'),
    cell: info => {
      const row = info.row.original

      return isTotalRow(row.service_description) ? '' : row.unit
    }
  }),
  columnHelper.accessor('unit_price', {
    header: t('unit_price'),
    cell: info => {
      const row = info.row.original

      return isTotalRow(row.service_description)
        ? ''
        : formatNumber(Number(row.unit_price), 'T')
    }
  }),
  columnHelper.accessor('service_price', {
    header: t('service_price'),
    cell: info => formatNumber(info.getValue(), 'T')
  }),
  columnHelper.accessor('maintenance', {
    header: t('maintenance'),
    cell: info => formatNumber(info.getValue(), 'T')
  }),
  columnHelper.accessor('total', {
    header: t('total'),
    cell: info => formatNumber(info.getValue(), 'T')
  })
]

export default function BreakdownTable({ data }: BreakdownTableProps) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel()
  })

  return (
    <Table>
      <TableHeader>
        {table.getHeaderGroups().map(headerGroup => (
          <TableRow key={headerGroup.id}>
            {headerGroup.headers.map(header => (
              <TableHead key={header.id}>
                {flexRender(
                  header.column.columnDef.header,
                  header.getContext()
                )}
              </TableHead>
            ))}
          </TableRow>
        ))}
      </TableHeader>

      <TableBody>
        {table.getRowModel().rows.map(row => {
          const isTotal = isTotalRow(row.original.service_description)

          return (
            <TableRow
              key={row.id}
              className={
                isTotal
                  ? 'bg-muted/30 font-semibold'
                  : 'transition hover:bg-muted/10'
              }
            >
              {row.getVisibleCells().map(cell => (
                <TableCell key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
              ))}
            </TableRow>
          )
        })}
      </TableBody>
    </Table>
  )
}
