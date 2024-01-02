import { Price } from "./models"

export function displayPriceInUI(price:Price):string{
  switch (price.currency) {
    case "IDR":
      return new Intl.NumberFormat('id-ID', {style: 'currency', currency: 'IDR'}).format(price.amount)
    default:
      return `${price.amount}`
  }
}