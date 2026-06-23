import { create } from 'zustand'

import type { TRegions } from '@/services/http/centralClient'
import { ProjectRegion, ProjectVariant } from '@/types/enums'

type FilterValues = {
  region: string
  teamFilter: string
  timeRange: string
  [key: string]: string | undefined
}

type FilterStore = {
  filterValues: FilterValues
  setTeam: (teamFilter: string) => void
  setFilterValues: (newValues: FilterValues) => void
  updateRegion: (
    region: string,
    updateBaseURL: (region: TRegions) => void
  ) => void
  init: (defaultTimeRange?: string) => void
}

export const useFilterStore = create<FilterStore>(set => ({
  filterValues: {
    region:
      import.meta.env.VITE_VARIANT === ProjectVariant.Cab
        ? ProjectRegion.Teh1
        : ProjectRegion.Box,
    teamFilter:
      import.meta.env.VITE_VARIANT === ProjectVariant.Cab
        ? 'smapp'
        : 'box-production',
    timeRange: '1h'
  },

  setTeam: teamFilter =>
    set(state => ({
      filterValues: { ...state.filterValues, teamFilter }
    })),

  setFilterValues: newValues => set({ filterValues: newValues }),

  updateRegion: (region, updateBaseURL) => {
    set(state => ({
      filterValues: { ...state.filterValues, region }
    }))
    updateBaseURL(region as TRegions)
  },

  // 👇 optional override
  init: defaultTimeRange =>
    set(state => ({
      filterValues: {
        ...state.filterValues,
        timeRange: defaultTimeRange ?? state.filterValues.timeRange
      }
    }))
}))
