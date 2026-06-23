import { useEffect } from 'react'

import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { useForm } from 'react-hook-form'
import * as z from 'zod'

import SelectSearch from '@/components/search-select'
import { Button } from '@/components/shadcn/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel
} from '@/components/shadcn/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/shadcn/select'
import useEffectOnce from '@/hooks/useEffectOnce'
import { useFilterStore } from '@/hooks/useFilterStore'
import { t } from '@/i18n'
import {
  type TRegions,
  updateRegion as updateRegionClient
} from '@/services/http/centralClient'
import { ProjectRegion, ProjectVariant } from '@/types/enums'

type TimeRangeOption = {
  title: string
  value: string | number
  default?: boolean
}

export interface ICustomFilter {
  label: string
  key: string
  items: {
    title: string
    value: string
    default?: boolean
  }[]
}

export interface GlobalFilterProps {
  teams: string[]
  timeRanges: TimeRangeOption[]
  disabled?: boolean
  customFilters?: ICustomFilter[]
  updateBaseUrlFunc?: (region: TRegions) => void
  onSubmit?: (values: Record<string, string>) => void
}

const formSchema = z.object({
  region: z.string(),
  teamFilter: z.string(),
  timeRange: z.string()
})

export function GlobalFilter({
  teams,
  timeRanges,
  disabled,
  customFilters,
  updateBaseUrlFunc = updateRegionClient,
  onSubmit
}: GlobalFilterProps) {
  const { filterValues, setFilterValues, updateRegion, setTeam } =
    useFilterStore()
  const navigate = useNavigate()
  const search = useSearch({ strict: false })
  const nsParam = search?.team as string

  const extendedSchema = formSchema.extend(
    Object.fromEntries(
      customFilters?.map(f => [f.key, z.string().optional()]) || []
    )
  )

  type BaseFilterValues = z.infer<typeof formSchema>
  type FilterFormValues = BaseFilterValues & {
    [key: string]: string | undefined
  }

  const defaultValues = {
    ...filterValues,
    ...Object.fromEntries(
      customFilters?.map(f => [
        f.key,
        f.items.find(item => item.default)?.value || ''
      ]) || []
    ),
    timeRange: String(timeRanges.find(item => item.default)?.value),
    teamFilter: filterValues.teamFilter
  }

  useEffectOnce(() => {
    setFilterValues(defaultValues)
  })

  const form = useForm<FilterFormValues>({
    resolver: zodResolver(extendedSchema),
    defaultValues
  })

  useEffect(() => {
    if (nsParam && nsParam !== form.getValues('teamFilter')) {
      setTeam(nsParam)
      form.setValue('teamFilter', nsParam)
    }
  }, [nsParam, form, setTeam])

  const projectVariant = import.meta.env.VITE_VARIANT as ProjectVariant

  const handleSubmit = async (values: FilterFormValues) => {
    setFilterValues(values)
    await navigate({
      search: prev => ({ ...prev, team: values.teamFilter })
    })
    onSubmit?.(values as Record<string, string>)
  }

  const handleRegionChange = (value: string) => {
    form.setValue('region', value)
    updateRegion(value, updateBaseUrlFunc)
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(handleSubmit)}
        className="grid grid-cols-4 gap-4"
      >
        <FormField
          control={form.control}
          name="region"
          render={({ field }) => (
            <FormItem className="grow">
              <FormLabel>{t('region')}</FormLabel>
              <Select
                disabled={disabled}
                onValueChange={handleRegionChange}
                value={field.value}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={t('select_region')} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {projectVariant === ProjectVariant.Cab ? (
                    <>
                      <SelectItem value={ProjectRegion.Teh1}>
                        {t('teh1')}
                      </SelectItem>
                      <SelectItem value={ProjectRegion.Teh2}>
                        {t('teh2')}
                      </SelectItem>
                    </>
                  ) : (
                    <SelectItem value={ProjectRegion.Box}>
                      {t('box')}
                    </SelectItem>
                  )}
                </SelectContent>
              </Select>
            </FormItem>
          )}
        />

        {customFilters?.map(customFilter => {
          return (
            <FormField
              control={form.control}
              name={customFilter.key}
              key={customFilter.key}
              render={() => (
                <FormItem>
                  <SelectSearch<FilterFormValues>
                    label={customFilter.label}
                    disabled={disabled}
                    fieldName={customFilter.key}
                    defaultItem={
                      customFilter.items.find(item => item.default)?.value || ''
                    }
                    placeholder={t('search_team')}
                    form={form}
                    items={customFilter.items}
                  />
                </FormItem>
              )}
            />
          )
        })}

        <SelectSearch<FilterFormValues>
          label={t('team')}
          disabled={disabled}
          fieldName="teamFilter"
          defaultItem={filterValues.teamFilter}
          placeholder={t('search_team')}
          form={form}
          items={teams.map(team => ({
            title: team,
            value: team
          }))}
        />

        <FormField
          control={form.control}
          name="timeRange"
          render={({ field }) => (
            <FormItem className="grow">
              <FormLabel>{t('time_range')}</FormLabel>
              <Select
                disabled={disabled}
                onValueChange={field.onChange}
                value={field.value}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={t('select_time_range')} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {timeRanges.map(item => (
                    <SelectItem key={item.value} value={String(item.value)}>
                      {item.title}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </FormItem>
          )}
        />
        <Button
          className={customFilters ? 'col-end-5 ml-auto mt-8' : 'mt-8'}
          type="submit"
          disabled={disabled}
        >
          {t('apply')}
        </Button>
      </form>
    </Form>
  )
}
