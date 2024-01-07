import { NextRouter } from "next/router";

type AppVersionResponseData = {
  version: string
}

export async function refreshIfNewAppVersionAvailable(
  appVersion: string,
  setAppVersion: (value: string)=>void,
  router: NextRouter,
  preRefreshCallback: ()=>void = ()=>{}, //default value: empty function
) {
  //uncomment this after GTM ~kodok
  // const response = await fetch(
  //   `${backendBaseUrl}/appVersion`,
  //   {
  //     method: "GET",
  //   }
  // );
  // const respJson: AppVersionResponseData = await response.json();
  // const currentVersion = respJson.version
  // if(appVersion !== currentVersion){
  //   preRefreshCallback()
  //   alert("Kami akan memuat ulang halaman ini agar anda dapat mendapatkan versi terbaru aplikasi kami.")
  //   setAppVersion(currentVersion)
  //   router.reload()
  // }
}
