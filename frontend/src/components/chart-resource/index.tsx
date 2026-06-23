import ErrorState from '@/components/error-state'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardFooter
} from '@/components/shadcn/card'

import { ResourceChart } from './resource-chart'
import { ResourceChartProps } from './resource.types'

const ResourceCostChart = ({
  header,
  footer,
  data,
  team,
  actualUsage,
  requestHard,
  labelData,
  tooltips,
  resource,
  hasError
}: ResourceChartProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{header}</CardTitle>
      </CardHeader>
      <CardContent className="pl-2">
        {hasError ? (
          <ErrorState />
        ) : (
          <ResourceChart
            team={team}
            resource={resource}
            actualUsage={actualUsage}
            requestHard={requestHard}
            labelData={labelData}
            data={data}
            tooltips={tooltips}
            header={header}
          />
        )}
      </CardContent>
      {footer ? <CardFooter className="px-6 pt-0">{footer}</CardFooter> : ''}
    </Card>
  )
}

export default ResourceCostChart
