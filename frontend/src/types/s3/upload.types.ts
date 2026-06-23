export interface IUploadNames {
  name: string
  completed: boolean
  failed: boolean
  progress: number
  canceled?: boolean
  file?: File
}
