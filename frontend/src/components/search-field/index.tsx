import { type FC } from 'react'

import { Search } from 'lucide-react'

import { Input } from '../shadcn/input'

interface SearchFieldProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
}

const SearchField: FC<SearchFieldProps> = ({
  value,
  onChange,
  placeholder = 'Search...'
}) => {
  return (
    <div className="relative flex w-[250px] items-center">
      <Search className="absolute right-2 top-1/2 size-4 -translate-y-1/2" />
      <Input
        placeholder={placeholder}
        value={value}
        onChange={e => onChange(e.target.value)}
        className="w-full pr-8"
      />
    </div>
  )
}

export default SearchField
