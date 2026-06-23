export type Namespace = {
  id: string
  namespace: string
  cpu_value: number
  memory_value: number
  pod_counts: number
  storage_value: number
  ephemeral_storage_value: number
  memory_unit: string
  ephemeral_storage_unit: string
  storage_unit: string
}

export type TeamUsageCardProps = {
  title: string
  showDownloadButton?: boolean
  data: Namespace[]
}
