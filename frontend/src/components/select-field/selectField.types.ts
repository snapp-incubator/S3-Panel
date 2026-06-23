export interface selectFieldProps {
  items: {
    title: string
    value: string
  }[]
  placeholder: string
  value: string
  onChange: (value: string) => void
}
