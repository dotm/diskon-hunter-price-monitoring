import { refreshIfNewAppVersionAvailable } from "@/utils/appversionutil";
import { LocalStorageKey, backendBaseUrl, backendHeadersForPostRequest } from "@/utils/constants";
import { convertCurrencyToNumber, convertCurrencyToNumberString, convertNumberStringToCurrency, displayCurrencyInUI } from "@/utils/currencyutil";
import { handleErrorInFrontend } from "@/utils/error";
import { disableChangingNumberValueOnScroll } from "@/utils/eventhandler";
import { LoggedInUserData, SearchedItemDetail, UserSearchesItemDetailEditRequestDTO, UserSearchesItemListEditRequestDTO, emptyUserSearchesItemListEditRequestDTO } from "@/utils/models";
import Link from "next/link";
import { useRouter } from 'next/router';
import { MouseEvent, useEffect, useState } from "react";
import useLocalStorage from "use-local-storage";

export default function SearchItemEditItems({
  appVersion,
  setAppVersion,
}:{
  appVersion: string,
  setAppVersion: (value: string)=>void,
}) {
  const router = useRouter()
  const [loggedInUserData, setLoggedInUserData] =
    useLocalStorage<LoggedInUserData | undefined>(LocalStorageKey.loggedInUser, undefined)
  const [requestDTO, setRequestDTO] = useLocalStorage<UserSearchesItemListEditRequestDTO>(
    LocalStorageKey.UserSearchesItemListEditRequestDTO,
    emptyUserSearchesItemListEditRequestDTO(),
  )
  const [loading, setLoading] = useState(false)

  function updateSearchedItemList(searchedItem: UserSearchesItemDetailEditRequestDTO){
    setRequestDTO({
      SearchedItemList: requestDTO.SearchedItemList.map(o=>{
        if(o.HubSearchedItemId !== searchedItem.HubSearchedItemId){
          return o
        }
        return searchedItem
      })
    })
  }
  function resetRequestDTO(){
    setRequestDTO(emptyUserSearchesItemListEditRequestDTO())
  }
  useEffect(() => {
    refreshIfNewAppVersionAvailable(appVersion, setAppVersion, router, function(){resetRequestDTO()})
  })

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
      if(!searchedItemListRespJson.ok || !searchedItemListRespJson.data){
        throw new Error(searchedItemListRespJson.err?.code ?? "error searchedItemListRespJson")
      }
      const initialValue: UserSearchesItemListEditRequestDTO = {
        SearchedItemList: searchedItemListRespJson.data.map((o: SearchedItemDetail):UserSearchesItemDetailEditRequestDTO=>{
          return {
            ...o,
            AlertPriceString: convertCurrencyToNumberString(o.AlertPrice)
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
    interactor_searchedItemList()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [true])

  async function interactor_searchedItemEditMultiple(event: MouseEvent<HTMLButtonElement, globalThis.MouseEvent>){
    event.preventDefault()
    try {
      setLoading(true)
      if(loggedInUserData === undefined){
        throw new Error("Mohon sign in terlebih dahulu")
      }
      
      if(requestDTO.SearchedItemList.find(o=>!o.Name) !== undefined){
        throw new Error("Mohon isi semua Nama Barang")
      }
      const processedSearchedItemList = requestDTO.SearchedItemList.map(o=>{
        return {
          ...o,
          AlertPrice: convertNumberStringToCurrency(o.AlertPriceString, "IDR")
        }
      })
      if(processedSearchedItemList.find(o=>convertCurrencyToNumber(o.AlertPrice) <= 0) !== undefined){
        throw new Error("Mohon isi semua Harga")
      }
      const validatedRequestDTO = {
        SearchedItemList: processedSearchedItemList
      }
      const searchedItemEditRespJson = await fetch(`${backendBaseUrl}/v1/searchedItem.editMultiple`, {
        method: 'POST',
        headers: backendHeadersForPostRequest(loggedInUserData.jwt),
        body: JSON.stringify(validatedRequestDTO),
      })
      .then(response => response.json())
      if(!searchedItemEditRespJson.ok || !searchedItemEditRespJson.data){
        throw new Error(searchedItemEditRespJson.err?.code ?? "error searchedItemEditRespJson")
      }
      //do something with searchedItemEditRespJson.data
      resetRequestDTO()
      alert(`Berhasil mengubah data`)
      router.replace('/searchItem')
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
        <></>
      }
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
        !loading
        ?
        <>
          <div className="divide-y divide-gray-400 bg-gray-800 text-white block w-[100%] px-3 pb-3 rounded-xl mx-auto">
            <h2 className="text-center font-bold p-2">Daftar Barang yang Anda Cari</h2>
            <p className="text-center text-sm pt-1 pb-2">
              Mohon klik Simpan untuk mengubah data yang anda telah edit.
            </p>
            <div className='flex justify-center flex-wrap text-center space-x-2'>
              <button
                type="button"
                className="my-2 rounded-md w-20 disabled:bg-slate-600 bg-blue-600 py-1 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                onClick={interactor_searchedItemEditMultiple}
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
                requestDTO.SearchedItemList.map((itemDetail)=>{
                  return (
                    <div
                      key={itemDetail.HubSearchedItemId} //edit form should use HubSearchedItemId
                      className="pt-4"
                    >
                      <div>
                        <label
                          htmlFor={`searchedItemEdit-name-${itemDetail.HubSearchedItemId}`}
                          className="sr-only block text-sm font-medium leading-6 text-gray-900"
                        >
                          Nama Barang
                        </label>
                        <div>
                          <input
                            id={`searchedItemEdit-name-${itemDetail.HubSearchedItemId}`}
                            name={`searchedItemEdit-name-${itemDetail.HubSearchedItemId}`}
                            type="text"
                            placeholder="Nama Barang"
                            value={itemDetail.Name}
                            onChange={event=>updateSearchedItemList({
                              ...itemDetail,
                              Name: event.target.value,
                            })}
                            className="placeholder:text-gray-400 block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                          />
                        </div>
                      </div>
                      <div>
                        <label
                          htmlFor={`searchedItemEdit-name-${itemDetail.HubSearchedItemId}`}
                          className="block text-sm font-medium leading-6 text-white"
                        >
                          Deskripsi Barang (max 2000 karakter)
                        </label>
                        <div>
                          <textarea
                            id={`searchedItemEdit-description-${itemDetail.HubSearchedItemId}`}
                            name={`searchedItemEdit-description-${itemDetail.HubSearchedItemId}`}
                            placeholder="misal: contoh link barang, link untuk gambar barang, spesifikasi teknis, tahun keluaran, jenis varian, warna, ukuran, dsb."
                            value={itemDetail.Description}
                            onChange={event=>updateSearchedItemList({
                              ...itemDetail,
                              Description: event.target.value,
                            })}
                            className="h-44 placeholder:text-gray-400 block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                          />
                        </div>
                      </div>
                      <div>
                        <label
                          htmlFor={`searchedItemEdit-url-${itemDetail.HubSearchedItemId}`}
                          className="block text-sm font-medium leading-6 text-gray-100"
                        >
                          Ingatkan saya pada harga: {displayCurrencyInUI(
                            convertNumberStringToCurrency(itemDetail.AlertPriceString,"IDR"))
                          }
                        </label>
                        <div>
                          <input
                            id={`searchedItemEdit-url-${itemDetail.HubSearchedItemId}`}
                            name={`searchedItemEdit-url-${itemDetail.HubSearchedItemId}`}
                            type="number"
                            placeholder="Masukkan harga (angka tanpa titik atau koma)"
                            value={itemDetail.AlertPriceString}
                            onWheel={disableChangingNumberValueOnScroll}
                            onKeyDown={event=>{
                              if(["e",".",",","-"].includes(event.key)){ //disable keys
                                event.preventDefault()
                              }
                            }}
                            onChange={event=>updateSearchedItemList({
                              ...itemDetail,
                              AlertPriceString: event.target.value,
                            })}
                            className="placeholder:text-gray-400 block w-full rounded-md border-0 bg-white/5 py-1.5 px-2.5 text-white shadow-sm ring-1 ring-inset ring-white/10 focus:ring-2 focus:ring-inset focus:ring-indigo-500 sm:text-sm sm:leading-6"
                          />
                        </div>
                      </div>
                    </div>
                  )
                })
              }
            </div>
          </div>
        </>
        :
        <></>
      }
    </div>
  )
}
