import { NextRouter } from "next/router";
import { WheelEvent } from "react";

export function disableChangingNumberValueOnScroll(e: WheelEvent<HTMLInputElement>){
  if (e.target instanceof HTMLElement) {
    e.target.blur()
  }else{
    console.log("failed bluring element")
  }
}

export function backButtonClicked(router: NextRouter, href: string = '/'){
  router.back()
  // how to check previous page is in your domain???
}