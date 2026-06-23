export const formatValue = (value: number): string => {
  if (typeof value !== 'number' || isNaN(value) || value < 0) {
    return '0'
  }

  // Convert bytes to GB
  const BYTES_PER_GB = 1024 ** 3

  if (value > 10000) {
    return (value / BYTES_PER_GB).toFixed(2)
  }

  // For values <= 10000, only use decimals if needed
  return value % 1 === 0 ? value.toString() : value.toFixed(1)
}

export const convertTimestamp = (timestamp: string) => {
  const date = new Date(timestamp)
  const day = date.getDate().toString().padStart(2, '0')
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')

  return `${month}-${day} [${hours}:${minutes}]`
}
