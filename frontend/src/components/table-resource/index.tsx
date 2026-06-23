import { ColumnDef } from '@tanstack/react-table'
import { Download } from 'lucide-react'

import { Button } from '@/components/shadcn/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle
} from '@/components/shadcn/card'
import { DataTable } from '@/components/table'
import type { TProject } from '@/types/quota/teams.type'

import { useColumns } from './column'
import { Namespace, ResourceTableProps } from './resource.types'

const ResourceTable = ({
  title,
  data,
  showDownloadButton = false,
  showLinkButton = true,
  isVM = false,
  isInstance = false
}: ResourceTableProps) => {
  const projectColumns = useColumns({ type: 'project', showLinkButton })
  const instanceColumns = useColumns({ type: 'instance', showLinkButton })
  const namespaceColumns = useColumns({ type: 'namespace', showLinkButton })

  return (
    <Card className="col-span-4 mt-5">
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
        {isVM ? (
          <DataTable<TProject>
            data={data as TProject[]}
            columns={
              isInstance
                ? (instanceColumns as ColumnDef<TProject>[])
                : (projectColumns as ColumnDef<TProject>[])
            }
            filterColumn="name"
          />
        ) : (
          <DataTable<Namespace>
            data={data as Namespace[]}
            columns={namespaceColumns as ColumnDef<Namespace>[]}
            filterColumn="namespace"
          />
        )}
      </CardContent>
    </Card>
  )
}

export default ResourceTable
