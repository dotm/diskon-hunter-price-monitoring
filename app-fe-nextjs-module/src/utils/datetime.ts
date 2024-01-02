import { format } from "date-fns";

export function displayDateInUI(date: Date): string {
  return format(
    date,
    'dd MMM yyyy',
  )
}
export function displayDateTimeInUI(date: Date): string {
  return format(
    date,
    'dd MMM yyyy HH:mm',
  )
}