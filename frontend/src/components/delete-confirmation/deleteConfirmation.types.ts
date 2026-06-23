export interface IDeleteConfirmationProps {
  open: boolean
  deleteItemName: string
  isLoading?: boolean
  closeHandler: () => void
  acceptDelete: () => void
}
