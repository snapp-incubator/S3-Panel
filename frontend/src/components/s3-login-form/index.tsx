import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

import { fetchBucketsList } from '@/api/s3'
import SelectField from '@/components/select-field'
import { Button } from '@/components/shadcn/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/shadcn/form'
import { Input } from '@/components/shadcn/input'
import { formItems } from '@/constants/s3/s3LoginForm'
import { useToast } from '@/hooks/use-toast'
import useS3Credentials from '@/hooks/useS3Credentials'
import { t } from '@/i18n'
import { updateRegion } from '@/services/http/centralClient'
import type { HTTPClientError } from '@/services/http/interceptorsConfig'
import { ProjectRegion } from '@/types/enums'
import type { TRegions } from '@/types/regions.types'

const formSchema = z.object({
  region: z.string({
    message: t('region_required')
  }),
  access_key: z.string().min(1, {
    message: t('access_key_required')
  }),
  secret_key: z.string().min(1, {
    message: t('secret_key_required')
  })
})

export default function S3LoginForm() {
  const navigate = useNavigate()
  const { toast } = useToast()

  const { fillCredentials } = useS3Credentials()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      region: ProjectRegion.Teh1,
      access_key: undefined,
      secret_key: undefined
    }
  })

  const { mutateAsync, isPending } = useMutation({
    mutationFn: async () => {
      const values = form.getValues()
      const { access_key, secret_key, region } = values

      // Set the credentials as separate headers
      updateRegion(region as ProjectRegion.Teh1 | ProjectRegion.Teh2, {
        access_key,
        secret_key
      })

      return await fetchBucketsList()
    },
    onError: (err: HTTPClientError<{ message?: string }>) => {
      toast({
        variant: 'destructive',
        title: err.response?.data.message || err.message
      })
    }
  })

  async function onSubmit(values: z.infer<typeof formSchema>) {
    const { secret_key, access_key, region } = values

    try {
      await mutateAsync()

      navigate({
        to: '/object-storage/s3-bucket/buckets'
      })

      fillCredentials(
        {
          secret_key,
          access_key
        },
        region as TRegions
      )
    } catch (err) {
      console.log(err)
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        {formItems.map(item => {
          return (
            <FormField
              control={form.control}
              name={item.key}
              key={item.key}
              render={({ field }) => (
                <FormItem className="flex items-baseline justify-between">
                  <FormLabel className="w-[200px]">{item.label}</FormLabel>
                  <div className="w-full">
                    <FormControl>
                      {item.type === 'select' ? (
                        <SelectField
                          items={item.selectItems!}
                          placeholder={item.placeholder}
                          value={field.value}
                          onChange={field.onChange}
                        />
                      ) : (
                        <Input
                          placeholder={item.placeholder}
                          {...field}
                          autoComplete="off"
                          data-test={`${item.key}_input`}
                        />
                      )}
                    </FormControl>
                    <FormMessage className="mt-3" />
                  </div>
                </FormItem>
              )}
            />
          )
        })}
        <div className="flex flex-row-reverse">
          <Button type="submit" disabled={isPending}>
            {isPending ? <Loader2 className="animate-spin" /> : t('submit')}
          </Button>
        </div>
      </form>
    </Form>
  )
}
