import { RouterProvider } from '@tanstack/react-router'

import { UploadProgressProvider } from './components/providers/uploadProgressContext'
import GlobalUploadProgress from './components/upload-progress'
import useS3Credentials from './hooks/useS3Credentials'
import { router } from './router'
import { updateRegion } from './services/http/centralClient'

function App() {
  const { access_key, secret_key, region } = useS3Credentials()

  if (access_key && secret_key && region) {
    updateRegion(region, {
      access_key,
      secret_key
    })
  }

  return (
    <UploadProgressProvider>
      <GlobalUploadProgress />
      <RouterProvider router={router} />
    </UploadProgressProvider>
  )
}

export default App
