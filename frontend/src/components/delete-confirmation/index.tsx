import { Button } from '@/components/shadcn/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/shadcn/dialog'
import { t } from '@/i18n'

import type { IDeleteConfirmationProps } from './deleteConfirmation.types'

export default function DeleteConfirmation({
  open,
  deleteItemName,
  isLoading,
  closeHandler,
  acceptDelete
}: IDeleteConfirmationProps) {
  return (
    <Dialog open={open} onOpenChange={closeHandler}>
      <DialogTrigger />
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t('delete_confirmation')}</DialogTitle>
        </DialogHeader>
        <h3 className=" text-lg">{t('delete_confirmation_title')}: </h3>
        <div className="truncate text-lg">
          <span className="font-bold">{deleteItemName}</span>
        </div>
        <DialogFooter>
          <Button
            disabled={isLoading}
            variant="destructive"
            data-test="confirm-delete-button"
            onClick={acceptDelete}
          >
            {t('delete')}
          </Button>
          <DialogClose asChild>
            <Button disabled={isLoading} type="button" variant="secondary" data-test="cancel-delete-button">
              {t('discard')}
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
