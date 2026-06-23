import { Button } from '@/components/shadcn/button'
import { t } from '@/i18n'

type TDeleteBucketProps = {
  disable: boolean
  clickHandler: () => void
}

export default function DeleteBucket({
  disable,
  clickHandler
}: TDeleteBucketProps) {
  return (
    <Button
      disabled={disable}
      variant="ghost"
      className="text-red-500"
      data-test="delete-bucket-button"
      onClick={clickHandler}
    >
      {t('delete_bucket')}
    </Button>
  )
}
