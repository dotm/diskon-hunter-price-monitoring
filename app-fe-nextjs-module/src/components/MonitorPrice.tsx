import { LocalStorageKey } from "@/utils/constants";
import { displayPriceInUI } from "@/utils/currencyutil";
import { displayDateInUI, displayDateTimeInUI } from "@/utils/datetime";
import { LoggedInUserData, UserLinkDetail } from "@/utils/models";
import Link from "next/link";
import { useRouter } from 'next/router';
import { useEffect, useState } from "react";
import useLocalStorage from "use-local-storage";

export default function MonitorPrice() {
  const router = useRouter()
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  const [userLinkList, setUserLinkList] = useState<UserLinkDetail[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const mockData: UserLinkDetail[] = [
      {
        HubUserId: "abc1234567890",
        HubMonitoredLinkUrl: "https://mock.com/product/1",
        AlertPrice: {amount: 10000, currency: "IDR"},
        ActiveAlertMethodList: ["Email"],
        PaidAlertMethodList: ["Email", "WhatsApp"],
        TimeExpired: new Date(11111111111111),
        LatestPrice: {amount: 9000, currency: "IDR"},
        TimeLatestScrapped: new Date(),
      },
      {
        HubUserId: "abc1234567890",
        HubMonitoredLinkUrl: "https://mock.com/product/2",
        AlertPrice: {amount: 20000, currency: "IDR"},
        ActiveAlertMethodList: [],
        PaidAlertMethodList: [],
        TimeExpired: new Date(11111111111111),
        LatestPrice: {amount: 30000, currency: "IDR"},
        TimeLatestScrapped: new Date(),
      },
    ]
    setUserLinkList(mockData)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [true])
  
  return (
    <div className="space-y-3">
      {
        loggedInUserData !== undefined
        ?
        <></>
        :
        <Link href="/settings" className="bg-gray-800 hover:bg-gray-700 outline-green-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto">
          <p className="text-sm">
          Jika anda sudah pernah melakukan sign up, silahkan klik text ini untuk sign in ke akun anda.
          <br/>
          Jika anda belum memiliki akun, pembuatan akun akan dilakukan setelah anda melakukan pembayaran.
          </p>
        </Link>
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
          userLinkList.length === 0
          ?
          <div className="bg-gray-800 text-white block w-[100%] pb-3 pt-2 mx-auto">
            <p className="text-center">
              Anda belum memonitor link apapun. Silahkan klik tombol Tambah untuk menambahkan link.
              <br/>
              Insert testimonial and example alert here ~kodok
            </p>
          </div>
          :
          <>
            {userLinkList.map(userLink => {
              return (
                <a key={userLink.HubMonitoredLinkUrl} href={userLink.HubMonitoredLinkUrl} target="_blank" className="block px-2 pt-2 pb-2 hover:bg-gray-700">
                  <p className="underline">
                    {userLink.HubMonitoredLinkUrl}
                  </p>
                  <p>
                    Harga terakhir: <strong>{
                    displayPriceInUI(userLink.LatestPrice)
                    }</strong> (pada <strong>{
                    displayDateTimeInUI(userLink.TimeLatestScrapped)
                    }</strong>)
                  </p>
                  <p>
                    Harga yang anda input: {displayPriceInUI(userLink.AlertPrice)}
                  </p>
                  <p>
                    Berhenti dimonitor pada: {displayDateInUI(userLink.TimeExpired)}
                  </p>
                </a>
              )
            })}
          </>
        }
      </div>
      {
        <>
        </>
      }
    </div>
  )
}
