import { LocalStorageKey, backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants";
import { displayCurrencyInUI } from "@/utils/currencyutil";
import { displayDateInUI } from "@/utils/datetime";
import { handleErrorInFrontend } from "@/utils/error";
import { LoggedInUserData, SearchedItemDetail } from "@/utils/models";
import Link from "next/link";
import { useRouter } from 'next/router';
import { useEffect, useState } from "react";
import useLocalStorage from "use-local-storage";

export default function SearchItemList() {
  const router = useRouter()
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  const [searchedItemList, setSearchedItemList] = useState<SearchedItemDetail[]>([])
  const [loading, setLoading] = useState(false)

  async function interactor_searchedItemList(){
    try {
      setLoading(true)
      if(loggedInUserData === undefined){
        throw new Error("Mohon sign in terlebih dahulu")
      }
      const searchedItemListRespJson = await fetch(`${backendBaseUrl}/v1/searchedItem.list`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(loggedInUserData.jwt),
        body: JSON.stringify({}),
      })
      .then(response => response.json())
      if(!searchedItemListRespJson.ok){
        throw new Error(searchedItemListRespJson.err?.code ?? "error searchedItemListRespJson")
      }
      setSearchedItemList(searchedItemListRespJson.data ?? [])
    } catch (error) {
      handleErrorInFrontend(error)
    } finally {
      setLoading(false)
    }
  }
  useEffect(() => {
    interactor_searchedItemList()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [true])

  return (
    <div className="space-y-3">
      {
        loggedInUserData === undefined
        ?
        <Link href="/settings" className="bg-gray-800 hover:bg-gray-700 outline-green-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto">
          <p className="text-sm">
          Jika anda sudah pernah melakukan sign up, silahkan klik text ini untuk sign in ke akun anda.
          <br/>
          Jika anda belum memiliki akun, pembuatan akun akan dilakukan setelah anda melakukan pembayaran.
          </p>
        </Link>
        :
        <></>
      }
      <Link href="/contactUs" className="bg-gray-800 hover:bg-gray-700 outline-red-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto">
        <p className="text-sm">
        Perlu bantuan? Punya pertanyaan, saran, atau kritik?
        <br/>
        Silahkan kontak kami dengan klik disini.
        </p>
      </Link>
      <div className="divide-y divide-gray-400 bg-gray-800 text-white block w-[100%] px-3 pb-3 rounded-xl mx-auto">
        <h2 className="text-center font-bold p-2">Daftar Barang yang Anda Cari</h2>
        <div className='flex justify-center flex-wrap text-center space-x-2'>
          <button
            type="button"
            className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-green-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-green-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
            onClick={() => {router.push(`/searchItem/add`);}}
            disabled={loading}
          >
            Tambah
          </button>
          <button
            type="button"
            className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-blue-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
            onClick={() => {router.push(`/searchItem/edit`);}}
            disabled={loading}
          >
            Edit
          </button>
        </div>
        {
          loading
          ?
          <div className="bg-gray-800 text-white block w-[100%] pb-3 pt-2 mx-auto">
            <p className="text-center">
              Mohon tunggu. Aplikasi sedang mengambil data...
            </p>
          </div>
          :
          <></>
        }
        {
          !loading && searchedItemList.length === 0
          ?
          <div className="bg-gray-800 text-white block w-[100%] pb-3 pt-2 mx-auto">
            <p className="text-center">
              Anda belum mencari barang apapun. Silahkan klik tombol Tambah untuk menambahkan link.
              <br/>
              <br/>
              Kami akan bantu mencarikan barang yang anda input.
              <br/>
              Jika kami berhasil menemukan barang tersebut sesuai dengan harga yang anda input,
              <br/>
              kami akan beritahu anda melalui WhatsApp
              <br/>
              Insert testimonial and example alert here ~kodok
            </p>
          </div>
          :
          <></>
        }
        {
          !loading && searchedItemList.length > 0
          ?
          <>
            {searchedItemList
              .sort((a,b)=>{
                return a.Name.localeCompare(b.Name)
              })
              .map(searchedItem => {
                return (
                  <div key={searchedItem.HubSearchedItemId} className="px-2 pt-2 pb-2 hover:bg-gray-700">
                    <p className="underline">
                      {searchedItem.Name}
                    </p>
                    <p>
                      Harga yang anda input: {displayCurrencyInUI(searchedItem.AlertPrice)}
                    </p>
                    <p>
                      Berhenti dicari pada: {displayDateInUI(searchedItem.TimeExpired)}
                    </p>
                    <p>
                      {searchedItem.Description}
                    </p>
                  </div>
                )
              })
            }
          </>
          :
          <></>
        }
      </div>
    </div>
  )
}
