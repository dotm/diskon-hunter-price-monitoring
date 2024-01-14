import MonitorPriceEditLinks from '@/components/MonitorPriceEditLinks'
import Sidebar from '@/components/Sidebar'
import { LocalStorageKey } from '@/utils/constants'
import useLocalStorage from 'use-local-storage'

export default function MonitorPriceEditLinksPage() {
  const [appVersion, setAppVersion] = useLocalStorage<string>(LocalStorageKey.appVersion, "")

  //refreshIfNewAppVersionAvailable useEffect is done inside the component

  return (
    <Sidebar pageTitle='Ubah Link Monitor Harga' navTitle='Monitor Harga'>
      <MonitorPriceEditLinks appVersion={appVersion} setAppVersion={setAppVersion} />
    </Sidebar>
  )
}
