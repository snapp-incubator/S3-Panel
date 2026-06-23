import { useEffect, useState, ReactNode } from 'react'

import { Info } from 'lucide-react'
import {
  Line,
  LineChart,
  Tooltip,
  CartesianGrid,
  XAxis,
  YAxis,
  TooltipProps,
  ResponsiveContainer,
  Legend
} from 'recharts'
import type { Props as LegendContentProps } from 'recharts/types/component/DefaultLegendContent'

import { AlertMessage } from '@/components/alert-message'
import { Badge } from '@/components/shadcn/badge'
import { ChartConfig, ChartContainer } from '@/components/shadcn/chart'
import {
  Tooltip as TooltipLegend,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger
} from '@/components/shadcn/tooltip'
import { t } from '@/i18n'
import { formatValue, convertTimestamp } from '@/lib/helper'

import {
  ResourceChartProps,
  DataPoint,
  LineConfig,
  ResourceData,
  DataPointItem
} from './resource.types'

const chartConfig: ChartConfig = {
  actualUsage: { color: '#21aa57' },
  limitsQuota: { color: '#ff0000' },
  limitsUsed: { color: '#0000ff' },
  requestQuota: { color: '#0020ff' },
  requestAllocated: { color: '#ffa500' },
  requestHard: { color: '#8b4513' }
}

export const ResourceChart = ({
  data,
  team,
  resource,
  labelData = '',
  actualUsage = true,
  requestHard = false,
  tooltips
}: ResourceChartProps) => {
  const [transformedData, setTransformedData] = useState<DataPoint[]>([])
  const [error, setError] = useState<string | null>(null)
  const [activeLines, setActiveLines] = useState<Set<string>>(new Set())

  const lineConfigs: LineConfig[] = [
    {
      dataKey: 'actualUsage',
      name: 'Actual Usage',
      color: chartConfig.actualUsage.color,
      strokeWidth: 1
    },
    {
      dataKey: 'limitsQuota',
      name: 'Limits Quota',
      color: chartConfig.limitsQuota.color,
      strokeWidth: 1,
      tooltip: tooltips?.quota_limit_help
    },
    {
      dataKey: 'limitsUsed',
      name: 'Limits Used',
      color: chartConfig.limitsUsed.color,
      strokeWidth: 1,
      strokeDasharray: '5 5',
      tooltip: tooltips?.limits_used
    },
    {
      dataKey: 'requestAllocated',
      name: 'Request Allocated',
      color: chartConfig.requestAllocated.color,
      strokeWidth: 1,
      strokeDasharray: '3 3',
      tooltip: tooltips?.allocated_help
    },
    {
      dataKey: 'requestHard',
      name: 'Requests Quota',
      color: chartConfig.requestHard.color,
      strokeWidth: 1,
      strokeDasharray: '3 3',
      tooltip: tooltips?.request_hard
    },
    {
      dataKey: 'requestQuota',
      name: 'Request Quota',
      color: chartConfig.requestQuota.color,
      strokeWidth: 1,
      strokeDasharray: '3 3',
      tooltip: tooltips?.allocated_help
    }
  ]

  const tooltipsHelp: Record<string, string | undefined> = {
    actualUsage: tooltips?.actual_usage,
    limitsQuota: tooltips?.quota_limit_help,
    limitsUsed: tooltips?.limits_used,
    requestAllocated: tooltips?.allocated_help,
    requestHard: tooltips?.request_hard,
    requestQuota: tooltips?.request_quota
  }

  useEffect(() => {
    if (typeof data !== 'object' || data === null) {
      setError(t('data_empty_204'))

      return
    }

    try {
      const dataMap = new Map<string, DataPoint>()
      const newActiveLines = new Set<string>()

      const processData = (
        key: keyof ResourceData,
        dataKey: keyof DataPoint
      ) => {
        const dataArray = data[key] ?? []

        if (dataArray.length > 0) {
          newActiveLines.add(dataKey)
        }

        dataArray.forEach((item: DataPointItem) => {
          const timestamp = convertTimestamp(String(item.timestamp))

          const existingData = dataMap.get(timestamp) || { timestamp }

          dataMap.set(timestamp, {
            ...existingData,
            [dataKey]: item.value
          })
        })
      }

      actualUsage && processData('actual_usage', 'actualUsage')
      processData('limits_quota', 'limitsQuota')
      processData('limits_used', 'limitsUsed')
      processData('request_allocated', 'requestAllocated')
      processData('request_quota', 'requestQuota')
      requestHard && processData('request_hard', 'requestHard')

      const newTransformedData = Array.from(dataMap.values())

      setTransformedData(newTransformedData)
      setActiveLines(newActiveLines)

      setError(null)
    } catch (err) {
      setError(`${t('error_processing_data')}: ${(err as Error).message}`)
    }
  }, [data, actualUsage, requestHard])

  const customTooltip = ({
    active,
    payload,
    label
  }: TooltipProps<number, string>) => {
    if (active && payload && payload.length) {
      return (
        <div className="chart-tooltip rounded border border-gray-300 bg-white p-2 shadow">
          <p className="font-bold">{`Time: ${label}`}</p>
          {payload.map((entry, index) => (
            <p key={index} style={{ color: entry.color }}>
              {`${entry.name}: ${formatValue(entry.value as number)} ${labelData}`}
            </p>
          ))}
          <p className="mt-2">
            {team != null ? (
              <Badge>{`Team: ${team?.toUpperCase()}`}</Badge>
            ) : (
              <Badge>{`Namespace: ${resource?.toUpperCase()}`}</Badge>
            )}
          </p>
        </div>
      )
    }

    return null
  }

  const renderLegend = (props: LegendContentProps): ReactNode => {
    const { payload } = props

    if (!payload) return null

    return (
      <ul className="mt-4 flex flex-wrap justify-center gap-4">
        {payload
          .filter(entry => activeLines.has(String(entry.dataKey)))
          .map((entry, index: number) => (
            <li key={`item-${index}`} className="flex items-center">
              <span
                style={{ backgroundColor: entry.color }}
                className="mr-2 inline-block size-3 rounded"
              ></span>
              <span>{entry.value}</span>
              {tooltipsHelp[String(entry.dataKey)] ? (
                <TooltipProvider delayDuration={50}>
                  <TooltipLegend delayDuration={50}>
                    <TooltipTrigger asChild>
                      <Info className="ml-2 size-4 cursor-help text-muted-foreground" />
                    </TooltipTrigger>
                    <TooltipContent>
                      <p>{tooltipsHelp[String(entry.dataKey)]}</p>
                    </TooltipContent>
                  </TooltipLegend>
                </TooltipProvider>
              ) : null}
            </li>
          ))}
      </ul>
    )
  }

  if (error) {
    return <AlertMessage variant="warning" title="Warning" message={error} />
  }

  return (
    <div className="w-full">
      <ChartContainer config={chartConfig} className="h-[250px] w-full">
        <ResponsiveContainer width="100%" height="100%">
          {transformedData.length > 0 ? (
            <LineChart data={transformedData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="timestamp"
                angle={-45}
                textAnchor="end"
                height={70}
                interval="preserveStartEnd"
                tickFormatter={(value: string) => value}
              />
              <YAxis tickFormatter={formatValue} domain={['auto', 'auto']} />
              <Tooltip content={customTooltip} />
              <Legend content={renderLegend} />

              {lineConfigs.map(config => (
                <Line
                  key={config.dataKey}
                  type="monotone"
                  dataKey={config.dataKey}
                  stroke={config.color}
                  name={config.name}
                  dot={false}
                  hide={!activeLines.has(config.dataKey)}
                  strokeWidth={config.strokeWidth}
                  strokeDasharray={config.strokeDasharray}
                />
              ))}
            </LineChart>
          ) : (
            <AlertMessage
              variant="destructive"
              title="Error"
              message={t('no_data')}
            />
          )}
        </ResponsiveContainer>
      </ChartContainer>
    </div>
  )
}

export default ResourceChart
