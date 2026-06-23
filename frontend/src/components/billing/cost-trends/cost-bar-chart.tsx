import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts'

import { ChartContainer } from '@/components/shadcn/chart'
import type { ITrendItems } from '@/services/api/billing'

import CustomTooltipBarChart from './custom-bar-chart-tooltip'

interface ICostBarChartProps {
  trends: ITrendItems[]
  filter: 'all' | 'okd' | 'openstack'
}

const CostBarChart = ({ trends, filter }: ICostBarChartProps) => {
  return (
    <ChartContainer config={{}} className="h-[500px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart
          width={500}
          height={300}
          data={trends}
          margin={{
            top: 20,
            right: 30,
            left: 20,
            bottom: 5
          }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis
            dataKey="date"
            angle={300}
            height={80}
            tickMargin={30}
            tick={{ dx: -20 }}
          />
          <YAxis
            width={100}
            tick={{ dx: -10 }}
            tickFormatter={(value: number) => `${value.toLocaleString()} T`}
          />
          <Tooltip content={<CustomTooltipBarChart />} />
          <Legend />
          {(filter === 'all' || filter === 'okd') && (
            <>
              <Bar
                dataKey="okd_storage"
                name="OKD Storage"
                stackId="a"
                fill="#5A6ACF"
                barSize={20}
              />
              <Bar
                dataKey="okd_compute"
                name="OKD Compute"
                stackId="a"
                fill="#394867"
                barSize={20}
                radius={filter === 'okd' ? [7, 7, 0, 0] : [0, 0, 0, 0]}
              />
            </>
          )}
          {(filter === 'all' || filter === 'openstack') && (
            <>
              <Bar
                dataKey="openstack_storage"
                name="Openstack Storage"
                stackId="a"
                fill="#2E7D32"
                barSize={20}
              />
              <Bar
                dataKey="openstack_compute"
                name="Openstack Compute"
                stackId="a"
                fill="#66BB6A"
                barSize={20}
                radius={filter === 'openstack' ? [7, 7, 0, 0] : [0, 0, 0, 0]}
              />
            </>
          )}
          {filter === 'all' ? (
            <Bar
              dataKey="maintenance"
              name="Maintenance"
              stackId="a"
              fill="#F7B7C3"
              barSize={20}
              radius={[7, 7, 0, 0]}
            />
          ) : null}
        </BarChart>
      </ResponsiveContainer>
    </ChartContainer>
  )
}

export default CostBarChart
