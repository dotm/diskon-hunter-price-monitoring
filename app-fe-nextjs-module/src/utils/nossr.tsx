import dynamic from 'next/dynamic'

const NoSSRWrapper = function(props: any) {
  return ( 
    <>{props.children}</> 
  )
}
export default dynamic(() => Promise.resolve(NoSSRWrapper), { 
  ssr: false 
})