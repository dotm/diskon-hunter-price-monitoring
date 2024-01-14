import SearchItemAddItems from '@/components/SearchItemAddItems'
import Sidebar from '@/components/Sidebar'
import { LocalStorageKey } from '@/utils/constants'
import useLocalStorage from 'use-local-storage'

export default function SearchItemAddItemsPage() {
  const [appVersion, setAppVersion] = useLocalStorage<string>(LocalStorageKey.appVersion, "")

  //refreshIfNewAppVersionAvailable useEffect is done inside the component

  return (
    <Sidebar pageTitle='Tambah Barang yang Dicari' navTitle='Cari Barang'>
      <SearchItemAddItems appVersion={appVersion} setAppVersion={setAppVersion} />
    </Sidebar>
  )
}
