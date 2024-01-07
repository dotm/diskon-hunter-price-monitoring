import { LocalStorageKey, backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants";
import { displayCurrencyInUI } from "@/utils/currencyutil";
import { displayDateInUI, displayDateTimeInUI } from "@/utils/datetime";
import { handleErrorInFrontend } from "@/utils/error";
import { AvailableAlertMethodList, LoggedInUserData, UserLinkDetail } from "@/utils/models";
import Link from "next/link";
import { useRouter } from 'next/router';
import { useEffect, useState } from "react";
import useLocalStorage from "use-local-storage";
import AlertMethodChip from "./AlertMethodChip";

export default function MonitorPriceList() {
  const router = useRouter()
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  const [userLinkList, setUserLinkList] = useState<UserLinkDetail[]>([])
  const [loading, setLoading] = useState(false)

  async function interactor_monitoredLinkList(){
    try {
      setLoading(true)
      if(loggedInUserData === undefined){
        throw new Error("Mohon sign in terlebih dahulu")
      }
      const monitoredLinkListRespJson = await fetch(`${backendBaseUrl}/v1/monitoredLink.list`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(loggedInUserData.jwt),
        body: JSON.stringify({}),
      })
      .then(response => response.json())
      if(!monitoredLinkListRespJson.ok || !monitoredLinkListRespJson.data){
        throw new Error(monitoredLinkListRespJson.err?.code ?? "error monitoredLinkListRespJson")
      }
      console.log("kodok",monitoredLinkListRespJson.data)
      setUserLinkList(monitoredLinkListRespJson.data)
    } catch (error) {
      handleErrorInFrontend(error)
    } finally {
      setLoading(false)
    }
  }
  useEffect(() => {
    interactor_monitoredLinkList()
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
      <Link href="/searchItem" className="bg-gray-800 hover:bg-gray-700 outline-yellow-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto">
        <p className="text-sm">
        Ragu dengan harga Rp.1000/link? Atau ingin mendapatkan informasi harga selain dari link e-commerce?
        <br/>
        Silahkan coba fitur eksperimental kami Cari Barang hanya Rp.500/5 barang dengan meng-klik disini!
        </p>
      </Link>
      <Link href="/contactUs" className="bg-gray-800 hover:bg-gray-700 outline-red-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto">
        <p className="text-sm">
        Perlu bantuan? Punya pertanyaan, saran, atau kritik?
        <br/>
        Silahkan kontak kami dengan klik disini
        </p>
      </Link>
      <div className="divide-y divide-gray-400 bg-gray-800 text-white block w-[100%] px-3 pb-3 rounded-xl mx-auto">
        <h2 className="text-center font-bold p-2">Daftar Link yang Anda Monitor</h2>
        <div className='flex justify-center flex-wrap text-center space-x-2'>
          <button
            type="button"
            className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-green-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-green-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
            onClick={() => {router.push(`/monitorPrice/add`);}}
            disabled={loading}
          >
            Tambah
          </button>
          <button
            type="button"
            className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-blue-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
            onClick={() => {router.push(`/monitorPrice/edit`);}}
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
          !loading && userLinkList.length === 0
          ?
          <div className="bg-gray-800 text-white block w-[100%] pb-3 pt-2 mx-auto">
            <p className="text-center">
              Anda belum memonitor link apapun. Silahkan klik tombol Tambah untuk menambahkan link.
              <br/>
              Insert testimonial and example alert here ~kodok
            </p>
          </div>
          :
          <></>
        }
        {
          !loading && userLinkList.length > 0
          ?
          <>
            {userLinkList.map(userLink => {
              console.log("kodok 1", userLink)
              return (
                <a key={userLink.HubMonitoredLinkUrl} href={userLink.HubMonitoredLinkUrl} target="_blank" className="block px-2 pt-2 pb-2 hover:bg-gray-700">
                  <p className="underline">
                    {userLink.HubMonitoredLinkUrl}
                  </p>
                  <p>
                    Harga terakhir: <strong>{
                    displayCurrencyInUI(userLink.LatestPrice)
                    }</strong> (pada <strong>{
                    displayDateTimeInUI(userLink.TimeLatestScrapped)
                    }</strong>)
                  </p>
                  <p>
                    Harga yang anda input: {displayCurrencyInUI(userLink.AlertPrice)}
                  </p>
                  <p>
                    Berhenti dimonitor pada: {displayDateInUI(userLink.TimeExpired)}
                  </p>
                  <div className="flex flex-row flex-wrap">
                    {
                      AvailableAlertMethodList
                      .filter(alertMethod=>userLink.PaidAlertMethodList?.includes(alertMethod.backendValue))
                      .map(alertMethod=>{
                        return (
                          <AlertMethodChip
                            key={alertMethod.backendValue}
                            name={alertMethod.frontendValue}
                            active={userLink.ActiveAlertMethodList?.includes(alertMethod.backendValue) ?? false}
                          />
                        )
                      })
                    }
                  </div>
                </a>
              )
            })}
          </>
          :
          <></>
        }
      </div>
    </div>
  )
}
