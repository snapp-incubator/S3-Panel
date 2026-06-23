export type TotalCostCardProps = {
  title: string
  unit?: string
  amount: {
    request: {
      title: string
      value: number
      help: string
    }
    limit: {
      title: string
      value: number
      help: string
    }
  }
  tooltipContent?: string
  hasError?: boolean
}
