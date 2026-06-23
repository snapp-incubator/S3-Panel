import type { ColumnDef } from '@tanstack/react-table'

import { t } from '@/i18n'
import type { IBareMetalQuota } from '@/services/api/usage/baremetals'

export const columns: ColumnDef<IBareMetalQuota>[] = [
  {
    id: 'index',
    header: '#',
    cell: ({ row }) => <span className="font-medium">{row.index + 1}</span>
  },

  {
    accessorKey: 'ServerName',
    header: t('server_name'),
    cell: ({ row }) => (
      <span className="font-medium">{row.getValue('ServerName')}</span>
    )
  },
  {
    accessorKey: 'RackMode',
    header: t('rack_mode'),
    cell: ({ row }) => (
      <span className="font-medium">{row.getValue('RackMode') || '-'}</span>
    )
  },
  {
    accessorKey: 'Generation',
    header: t('generation'),
    cell: ({ row }) => (
      <span className="font-medium">G{row.getValue('Generation')}</span>
    )
  },
  {
    accessorKey: 'Memory',
    header: () => t('memory_gb'),
    cell: ({ row }) => <div>{row.getValue('Memory')} GB</div>
  },
  {
    accessorKey: 'StorageSSD',
    header: () => t('ssd_gb'),
    cell: ({ row }) => <div>{row.getValue('StorageSSD')} GB</div>
  }
]
