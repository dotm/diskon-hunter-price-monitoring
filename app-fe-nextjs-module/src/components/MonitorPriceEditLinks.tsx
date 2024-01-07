import { refreshIfNewAppVersionAvailable } from "@/utils/appversionutil";
import { LocalStorageKey, backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants";
import { convertCurrencyToIntegerString, convertCurrencyToNumber, convertNumberStringToCurrency, displayCurrencyInUI } from "@/utils/currencyutil";
import { handleErrorInFrontend } from "@/utils/error";
import { AvailableAlertMethodList, LoggedInUserData, UserLinkDetail, UserMonitorsLinkDetailEditRequestDTO, UserMonitorsLinkListEditRequestDTO, emptyUserMonitorsLinkListEditRequestDTO } from "@/utils/models";
import Link from "next/link";
import { useRouter } from 'next/router';
import { MouseEvent, useEffect, useState } from "react";
import useLocalStorage from "use-local-storage";
import AlertMethodChip from "./AlertMethodChip";

export default function MonitorPriceEditLinks({
  appVersion,
  setAppVersion,
}:{
  appVersion: string,
  setAppVersion: (value: string)=>void,
}) {
  const router = useRouter()
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  const [requestDTO, setRequestDTO] = useLocalStorage<UserMonitorsLinkListEditRequestDTO>(
    LocalStorageKey.UserMonitorsLinkListEditRequestDTO,
    emptyUserMonitorsLinkListEditRequestDTO(),
  )
  const [loading, setLoading] = useState(false)

  function updateMonitoredLinkList(monitoredLink: UserMonitorsLinkDetailEditRequestDTO){
    setRequestDTO({
      MonitoredLinkList: requestDTO.MonitoredLinkList.map(o=>{
        if(o.HubMonitoredLinkUrl !== monitoredLink.HubMonitoredLinkUrl){
          return o
        }
        return monitoredLink
      })
    })
  }
  function resetRequestDTO(){
    setRequestDTO(emptyUserMonitorsLinkListEditRequestDTO())
  }
  useEffect(() => {
    refreshIfNewAppVersionAvailable(appVersion, setAppVersion, router, function(){resetRequestDTO()})
  })

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
      const initialValue: UserMonitorsLinkListEditRequestDTO = {
        MonitoredLinkList: monitoredLinkListRespJson.data.map((o: UserLinkDetail):UserMonitorsLinkDetailEditRequestDTO=>{
          return {
            ...o,
            AlertPriceString: convertCurrencyToIntegerString(o.AlertPrice)
          }
        })
      }
      setRequestDTO(initialValue)
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

  async function interactor_monitoredLinkEditMultiple(event: MouseEvent<HTMLButtonElement, globalThis.MouseEvent>){
    event.preventDefault()
    try {
      setLoading(true)
      if(loggedInUserData === undefined){
        throw new Error("Mohon sign in terlebih dahulu")
      }
      
      if(requestDTO.MonitoredLinkList.find(o=>!o.HubMonitoredLinkUrl) !== undefined){
        throw new Error("Mohon isi semua URL")
      }
      const processedMonitoredLinkList = requestDTO.MonitoredLinkList.map(o=>{
        return {
          ...o,
          AlertPrice: convertNumberStringToCurrency(o.AlertPriceString, "IDR")
        }
      })
      if(processedMonitoredLinkList.find(o=>convertCurrencyToNumber(o.AlertPrice) <= 0) !== undefined){
        throw new Error("Mohon isi semua Harga")
      }
      const validatedRequestDTO = {
        MonitoredLinkList: processedMonitoredLinkList
      }
      console.log("kodok req", validatedRequestDTO)
      const monitoredLinkEditRespJson = await fetch(`${backendBaseUrl}/v1/monitoredLink.editMultiple`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(loggedInUserData.jwt),
        body: JSON.stringify(validatedRequestDTO),
      })
      .then(response => response.json())
      if(!monitoredLinkEditRespJson.ok || !monitoredLinkEditRespJson.data){
        throw new Error(monitoredLinkEditRespJson.err?.code ?? "error monitoredLinkEditRespJson")
      }
      //do something with monitoredLinkEditRespJson.data
      resetRequestDTO()
      console.log("kodok res", monitoredLinkEditRespJson)
      alert(`Berhasil mengubah data`)
      router.replace('/')
    } catch (error) {
      handleErrorInFrontend(error)
    } finally {
      setLoading(false)
    }
  }
  
  return (
    <div className="space-y-3">
      {
        loggedInUserData === undefined
        ?
        <Link href="/settings" className="bg-gray-800 hover:bg-gray-700 outline-green-500 text-white block w-[100%] outline pb-2 pt-1 px-3 rounded-xl mx-auto">
          <p className="text-sm">
          Mohon sign in ke akun anda terlebih dahulu.
          </p>
        </Link>
        :
        <>
          <div className="divide-y divide-gray-400 bg-gray-800 text-white block w-[100%] px-3 pb-3 rounded-xl mx-auto">
            <h2 className="text-center font-bold p-2">Daftar Link yang Anda Monitor</h2>
            <p className="text-center text-sm pt-1 pb-2">
              Anda tidak dapat mengubah link yang anda telah monitor.
              <br />
              Silahkan menambahkan link baru jika anda ingin memonitor link diluar yang telah anda monitor.
              <br />
              Mohon klik Simpan untuk mengubah data yang anda telah edit.
            </p>
            <div className='flex justify-center flex-wrap text-center space-x-2'>
              <button
                type="button"
                className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-blue-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                onClick={interactor_monitoredLinkEditMultiple}
                disabled={loading}
              >
                Simpan
              </button>
              <button
                type="button"
                className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-red-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
                onClick={() => {
                  resetRequestDTO()
                  router.back()
                }}
                disabled={loading}
              >
                Kembali
              </button>
            </div>
            <div className="bg-gray-800 text-white block w-[100%] pb-3 px-3 mx-auto space-y-4 divide-y divide-gray-400">
              {
                requestDTO.MonitoredLinkList.map((linkDetail)=>{
                  return (
                    <div
                      key={linkDetail.HubMonitoredLinkUrl} //edit form should use HubMonitoredLinkUrl
                      className="pt-4"
                    >
                      <p className="text-center">{linkDetail.HubMonitoredLinkUrl}</p>
                      <div>
                        <label
                          htmlFor={`monitoredLinkEdit-url-${linkDetail.HubMonitoredLinkUrl}`}
                          className="block text-sm font-medium leading-6 text-gray-100"
                        >
                          Ingatkan saya pada harga: {displayCurrencyInUI(
                            convertNumberStringToCurrency(linkDetail.AlertPriceString,"IDR"))
                          }
                        </label>
                        <div>
                          <input
                            id={`monitoredLinkEdit-url-${linkDetail.HubMonitoredLinkUrl}`}
                            name={`monitoredLinkEdit-url-${linkDetail.HubMonitoredLinkUrl}`}
                            type="number"
                            placeholder="Masukkan harga (angka tanpa titik atau koma)"
                            value={linkDetail.AlertPriceString}
                            onKeyDown={event=>{
                              if(["e",".",",","-"].includes(event.key)){ //disable keys
                                event.preventDefault()
                              }
                            }}
                            onChange={event=>updateMonitoredLinkList({
                              ...linkDetail,
                              AlertPriceString: event.target.value,
                            })}
                            className="placeholder:text-gray-400 block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                          />
                        </div>
                      </div>
                      <div className="flex flex-row mt-2 flex-wrap">
                        {
                          AvailableAlertMethodList
                          .filter(alertMethod=>linkDetail.PaidAlertMethodList?.includes(alertMethod.backendValue))
                          .map(alertMethod=>{
                            return (
                              <AlertMethodChip
                                key={alertMethod.backendValue}
                                name={alertMethod.frontendValue}
                                active={linkDetail.ActiveAlertMethodList?.includes(alertMethod.backendValue) ?? false}
                                onClick={(active)=>{
                                  let newData = (linkDetail.ActiveAlertMethodList ?? [])
                                    .filter(s=>s!==alertMethod.backendValue)
                                  if(active){
                                    newData.push(alertMethod.backendValue)
                                  }
                                  updateMonitoredLinkList({
                                    ...linkDetail,
                                    ActiveAlertMethodList: newData,
                                  })
                                }}
                              />
                            )
                          })
                        }
                      </div>
                    </div>
                  )
                })
              }
            </div>
          </div>
        </>
      }
    </div>
  )
}
