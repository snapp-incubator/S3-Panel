// import { useState } from 'react'
// import { type ChangeEvent } from 'react'

import {
  Select,
  SelectContent,
  //   SelectGroup,
  SelectItem,
  //   SelectLabel,
  SelectTrigger,
  SelectValue
} from '@/components/shadcn/select'

import type { selectFieldProps } from './selectField.types'

export default function SelectField({
  items,
  placeholder,
  value,
  onChange
}: selectFieldProps) {
  //   const [selectValue, setSelectValue] = useState<null | string>(null)

  //   const changeValue = (newVal: string) => {
  //     setSelectValue(newVal)
  //   }

  return (
    <Select value={value} onValueChange={onChange}>
      <SelectTrigger>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {items.map(item => {
          return (
            <SelectItem value={item.value} key={item.value}>
              {item.title}
            </SelectItem>
          )
        })}
      </SelectContent>
    </Select>
  )
}
