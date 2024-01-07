import { refreshIfNewAppVersionAvailable } from "@/utils/appversionutil";
import { LocalStorageKey, backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants";
import { convertCurrencyToNumber, convertNumberStringToCurrency, createZeroCurrency, displayCurrencyInUI } from "@/utils/currencyutil";
import { handleErrorInFrontend } from "@/utils/error";
import { AlertMethod, AvailableAlertMethodList, LoggedInUserData, UserMonitorsLinkDetailAddRequestDTO, UserMonitorsLinkListAddRequestDTO, emptyUserMonitorsLinkListAddRequestDTO } from "@/utils/models";
import Link from "next/link";
import { useRouter } from 'next/router';
import { MouseEvent, useEffect, useState } from "react";
import useLocalStorage from "use-local-storage";
import { v4 } from 'uuid';
import AlertMethodChip from "./AlertMethodChip";

export default function MonitorPriceAddLinks({
  appVersion,
  setAppVersion,
}:{
  appVersion: string,
  setAppVersion: (value: string)=>void,
}) {
  const router = useRouter()
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  const [requestDTO, setRequestDTO] = useLocalStorage<UserMonitorsLinkListAddRequestDTO>(
    LocalStorageKey.UserMonitorsLinkListAddRequestDTO,
    emptyUserMonitorsLinkListAddRequestDTO(),
  )
  const [loading, setLoading] = useState(false)

  function updateMonitoredLinkList(monitoredLink: UserMonitorsLinkDetailAddRequestDTO){
    setRequestDTO({
      MonitoredLinkList: requestDTO.MonitoredLinkList.map(o=>{
        if(o.FrontendID !== monitoredLink.FrontendID){
          return o
        }
        return monitoredLink
      })
    })
  }
  function resetRequestDTO(){
    setRequestDTO(emptyUserMonitorsLinkListAddRequestDTO())
  }
  useEffect(() => {
    refreshIfNewAppVersionAvailable(appVersion, setAppVersion, router, function(){resetRequestDTO()})
  })
  async function interactor_monitoredLinkAddMultiple(event: MouseEvent<HTMLButtonElement, globalThis.MouseEvent>){
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
      const monitoredLinkAddRespJson = await fetch(`${backendBaseUrl}/v1/monitoredLink.addMultiple`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(loggedInUserData.jwt),
        body: JSON.stringify(validatedRequestDTO),
      })
      .then(response => response.json())
      if(!monitoredLinkAddRespJson.ok || !monitoredLinkAddRespJson.data){
        throw new Error(monitoredLinkAddRespJson.err?.code ?? "error monitoredLinkAddRespJson")
      }
      //do something with monitoredLinkAddRespJson.data
      resetRequestDTO()
      console.log("kodok res", monitoredLinkAddRespJson)
      alert(`Berhasil memonitor link`)
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
              Klik Tambah Link untuk mulai menambahkan link yang anda ingin monitor.
              <br />
              Mohon klik Simpan untuk mulai memonitor link yang anda telah Tambahkan.
              <br />
              Anda akan memonitor semua link yang anda Tambahkan selama satu tahun kedepan.
            </p>
            <div className='flex justify-center flex-wrap text-center space-x-2'>
              <button
                type="button"
                className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-blue-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                onClick={interactor_monitoredLinkAddMultiple}
                disabled={loading}
              >
                Simpan
              </button>
              <button
                type="button"
                className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-red-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
                onClick={() => {
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
                      key={linkDetail.FrontendID} //edit form should use HubMonitoredLinkUrl
                      className="pt-4"
                    >
                      <div>
                        <label
                          htmlFor={`monitoredLinkAdd-url-${linkDetail.FrontendID}`}
                          className="sr-only block text-sm font-medium leading-6 text-gray-900"
                        >
                          Alamat URL
                        </label>
                        <div>
                          <input
                            id={`monitoredLinkAdd-url-${linkDetail.FrontendID}`}
                            name={`monitoredLinkAdd-url-${linkDetail.FrontendID}`}
                            type="text"
                            placeholder="Alamat URL (misal: https://ecommerce.com/produk-yang-akan-dimonitor)"
                            value={linkDetail.HubMonitoredLinkUrl}
                            onChange={event=>updateMonitoredLinkList({
                              ...linkDetail,
                              HubMonitoredLinkUrl: event.target.value,
                            })}
                            className="placeholder:text-gray-400 block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                          />
                        </div>
                      </div>
                      <div>
                        <label
                          htmlFor={`monitoredLinkAdd-url-${linkDetail.FrontendID}`}
                          className="block text-sm font-medium leading-6 text-gray-100"
                        >
                          Ingatkan saya pada harga: {displayCurrencyInUI(
                            convertNumberStringToCurrency(linkDetail.AlertPriceString,"IDR"))
                          }
                        </label>
                        <div>
                          <input
                            id={`monitoredLinkAdd-url-${linkDetail.FrontendID}`}
                            name={`monitoredLinkAdd-url-${linkDetail.FrontendID}`}
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
                          AvailableAlertMethodList.map(alertMethod=>{
                            return (
                              <AlertMethodChip
                                key={alertMethod.backendValue}
                                name={alertMethod.frontendValue}
                                active={linkDetail.AlertMethodList.includes(alertMethod.backendValue)}
                                onClick={(active)=>{
                                  let newData = linkDetail.AlertMethodList
                                    .filter(s=>s!==alertMethod.backendValue)
                                  if(active){
                                    newData.push(alertMethod.backendValue)
                                  }
                                  updateMonitoredLinkList({
                                    ...linkDetail,
                                    AlertMethodList: newData,
                                  })
                                }}
                              />
                            )
                          })
                        }
                      </div>
                      <div className='flex justify-end flex-wrap text-center space-x-2'>
                        <button
                          type="button"
                          className="mt-3 rounded-md w-40 disabled:bg-slate-600 bg-red-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
                          onClick={()=>{
                            setRequestDTO({
                              MonitoredLinkList: [...requestDTO.MonitoredLinkList.filter(o=>o.FrontendID!==linkDetail.FrontendID)]
                            })
                          }}
                          disabled={loading}
                        >
                          Hapus Link
                        </button>
                      </div>
                    </div>
                  )
                })
              }
              <div className='flex justify-center flex-wrap text-center space-x-2 mt-5'>
                <button
                  type="button"
                  className="my-2 rounded-md w-40 disabled:bg-slate-600 bg-green-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-green-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
                  onClick={()=>{
                    setRequestDTO({
                      MonitoredLinkList: [...requestDTO.MonitoredLinkList, {
                        FrontendID: v4(),
                        HubMonitoredLinkUrl: "",
                        AlertPrice: createZeroCurrency("IDR"),
                        AlertPriceString: "",
                        AlertMethodList: [AlertMethod.Email.backendValue]
                      }]
                    })
                  }}
                  disabled={loading}
                >
                  Tambah Link
                </button>
              </div>
            </div>
          </div>
        </>
      }
    </div>
  )
}
