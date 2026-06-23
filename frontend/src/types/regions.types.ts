// Regions are configuration-driven (advertised by the backend /regions
// endpoint), so the type is an open string rather than a fixed union.
export type TRegions = string

export interface IRegionsResponse {
  current: string
  regions: string[]
}
