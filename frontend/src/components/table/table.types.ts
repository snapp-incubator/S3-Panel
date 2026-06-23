import { ColumnDef } from '@tanstack/react-table'

export type DataTableProps<TData> = {
  data: TData[]
  columns: ColumnDef<TData>[]
  filterColumn?: string
  columnVisibilityDropdown?: boolean
}
