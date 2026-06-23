/**
 * Expiration options for S3 object sharing
 *
 * These values represent the available time periods for which
 * a shared link will remain valid. The 'value' property is sent
 * to the API, while 'title' is displayed to the user.
 *
 */

export const expirationItems = [
  {
    title: '1 Hour',
    value: '1h'
  },
  {
    title: '6 Hour',
    value: '6h'
  },
  {
    title: '12 Hour',
    value: '12h'
  },
  {
    title: '1 Day',
    value: '1d'
  },
  {
    title: '3 Day',
    value: '3d'
  },
  {
    title: '1 Week',
    value: '1w'
  },
  {
    title: '2 Week',
    value: '2w'
  }
]
