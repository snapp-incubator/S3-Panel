import { useEffectOnce } from 'react-use'

import { useTitle } from '@/components/providers/titleProvider'
import S3LoginForm from '@/components/s3-login-form'
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle
} from '@/components/shadcn/card'
import { t } from '@/i18n'

export default function S3BucketPage() {
  const { setTitle } = useTitle()

  useEffectOnce(() => {
    setTitle(t('s3_pages_title'))
  })

  return (
    <div className="flex min-h-full min-w-full flex-col">
      <h2 className="text-3xl">{t('s3_bucket')}</h2>
      <Card className="m-auto min-h-[363px] w-[500px] p-6 pb-0">
        <CardTitle>{t('s3_bucket')}</CardTitle>
        <CardDescription className="mt-4">
          {t('s3_bucket_login_description')}
        </CardDescription>
        <CardContent className="pt-6">
          <S3LoginForm />
        </CardContent>
      </Card>
    </div>
  )
}
