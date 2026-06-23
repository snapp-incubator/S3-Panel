import { ProjectRegion } from '@/types/enums'

interface IFormItems {
  label: string
  type: 'select' | 'text'
  selectItems?: Record<'title' | 'value', string>[]
  placeholder: string
  key: 'region' | 'access_key' | 'secret_key'
}

const regionItems = [
  {
    title: 'Teh 1',
    value: ProjectRegion.Teh1
  },
  {
    title: 'Teh 2',
    value: ProjectRegion.Teh2
  }
]

export const formItems: IFormItems[] = [
  {
    label: 'Region Selection',
    type: 'select',
    selectItems: regionItems,
    placeholder: 'Select a region',
    key: 'region'
  },
  {
    label: 'Access Key',
    type: 'text',
    placeholder: 'Insert Access Key of Bucket',
    key: 'access_key'
  },
  {
    label: 'Secret Key',
    type: 'text',
    placeholder: 'Insert Secret Key',
    key: 'secret_key'
  }
]
