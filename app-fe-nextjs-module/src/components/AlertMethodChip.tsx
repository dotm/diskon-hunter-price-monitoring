import { CheckCircleIcon, XCircleIcon } from "@heroicons/react/20/solid";

export default function AlertMethodChip({
  name,
  active,
  onClick,
}: {
  name: string,
  active: boolean,
  onClick?: ((currentlyActive:boolean)=>void),
}){
  return (
    <button
      type="button"
      className={`
        inline-flex items-center gap-x-1.5 rounded-full px-3 py-2 text-sm font-semibold text-white shadow-sm
        ${active ?
          "bg-indigo-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          :
          "bg-transparent ring-1 ring-inset ring-indigo-600 "
        }
        ${!onClick ? "cursor-default" : ""}
        ${onClick && active ? "hover:bg-indigo-500": ""}
        ${onClick && !active ? "hover:bg-gray-500": ""}
      `}
      onClick={function(){ if(onClick){ onClick(!active) } }}
    >
      {name}
      {active ?
      <CheckCircleIcon className="-mr-0.5 h-5 w-5" aria-hidden="true" />
      :
      <XCircleIcon className="-mr-0.5 h-5 w-5" aria-hidden="true" />
      }
    </button>
  )
}