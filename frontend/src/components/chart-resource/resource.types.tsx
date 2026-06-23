type TChartTooltips = {
  actual_usage?: string
  allocated_help?: string
  limits_used?: string
  request_hard?: string
  quota_limit_help?: string
  request_quota?: string
}

export type DataPointItem = {
  timestamp: string
  value: number
}

export type ResourceData = {
  actual_usage?: DataPointItem[]
  limits_quota?: DataPointItem[]
  limits_used?: DataPointItem[]
  request_allocated?: DataPointItem[]
  request_quota?: DataPointItem[]
  request_hard?: DataPointItem[]
}

export type ResourceChartProps = {
  data: ResourceData | undefined
  team?: string
  resource?: string
  actualUsage?: boolean
  requestHard?: boolean
  labelData?: string
  tooltips?: TChartTooltips
  header: string | React.ReactNode
  footer?: React.ReactNode
  title?: string
  hasError?: boolean
}

export type DataPoint = {
  timestamp: string
  actualUsage?: number
  limitsQuota?: number
  limitsUsed?: number
  requestAllocated?: number
  requestQuota?: number
  requestHard?: number
}

export type LineConfig = {
  dataKey: keyof DataPoint
  strokeWidth: number
  name: string
  color?: string
  strokeDasharray?: string
  orientation?: 'left' | 'right'
  yAxisId?: string
  tooltip?: string
}
