import type { UseFormReturn, FieldValues, Path } from 'react-hook-form'

export interface IItems {
  title: string
  value: string
}

export type SelectSearchProps<T extends FieldValues> = {
  fieldName: Path<T>
  label: string
  placeholder?: string
  disabled?: boolean
  defaultItem: string
  items: IItems[]
  form: UseFormReturn<T>
}
