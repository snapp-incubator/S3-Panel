import { TooltipProps } from 'recharts'

import { t } from '@/i18n'

const CustomTooltipBarChart = ({
  active,
  payload,
  label
}: TooltipProps<number, string>) => {
  if (active && payload && payload.length) {
    const total = payload.reduce((sum, entry) => sum + (entry.value || 0), 0)

    return (
      <div className="rounded border bg-white p-2 shadow dark:bg-slate-800">
        <p className="font-semibold">{label}</p>
        {payload.map((entry, index) => (
          <p key={`item-${index}`} style={{ color: entry.fill }}>
            {entry.name}: {entry.value?.toLocaleString()} T
          </p>
        ))}
        <p className="font-bold">
          {`${t('total')}: `} {total.toLocaleString()} T
        </p>
      </div>
    )
  }

  return null
}

export default CustomTooltipBarChart
