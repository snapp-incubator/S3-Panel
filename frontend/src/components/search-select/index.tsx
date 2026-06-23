import { useState } from 'react'

import { Check, ChevronsUpDown } from 'lucide-react'
import type { FieldValues } from 'react-hook-form'

import { Button } from '@/components/shadcn/button'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList
} from '@/components/shadcn/command'
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl
} from '@/components/shadcn/form'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/components/shadcn/popover'
import { cn } from '@/lib/utils'

import { SelectSearchProps } from './search.types'

const SelectSearch = <T extends FieldValues>({
  fieldName,
  label,
  placeholder,
  items,
  form,
  defaultItem,
  disabled
}: SelectSearchProps<T>) => {
  const [open, setOpen] = useState(false)

  return (
    <FormField
      control={form.control}
      name={fieldName}
      render={({ field }) => (
        <FormItem className="grow">
          <FormLabel>{label}</FormLabel>
          <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
              <FormControl className="w-full">
                <Button
                  disabled={disabled}
                  variant="outline"
                  role="combobox"
                  aria-expanded={open}
                  className="w-full text-start"
                >
                  <div className="w-full">
                    {field?.value
                      ? items?.find(item => item.value === field.value)?.title
                      : defaultItem}
                  </div>
                  <ChevronsUpDown className="ml-2 size-4 shrink-0 opacity-50" />
                </Button>
              </FormControl>
            </PopoverTrigger>
            <PopoverContent
              className="w-[--radix-popover-trigger-width] p-0"
              align="start"
              sideOffset={4}
            >
              <Command className="w-full">
                <CommandInput placeholder={placeholder} className="w-full" />
                <CommandList className="w-full">
                  <CommandEmpty>No team found.</CommandEmpty>
                  <CommandGroup>
                    {items?.map(item => (
                      <CommandItem
                        key={item.value}
                        value={item.value}
                        onSelect={currentValue => {
                          field.onChange(
                            currentValue === field.value ? '' : currentValue
                          )
                          setOpen(false)
                        }}
                      >
                        <Check
                          className={cn(
                            'mr-2 h-4 w-4',
                            field.value === item.value
                              ? 'opacity-100'
                              : 'opacity-0'
                          )}
                        />
                        {item.title}
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
        </FormItem>
      )}
    />
  )
}

export default SelectSearch
