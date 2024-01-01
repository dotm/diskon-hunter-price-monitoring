import '@/styles/globals.css'
import NoSSRWrapper from '@/utils/nossr'
import type { AppProps } from 'next/app'
import Head from 'next/head'

export default function App({ Component, pageProps }: AppProps) {
  return (
    <NoSSRWrapper>
      <Head>
        <title>Diskon Hunter</title>
      </Head>
      <Component {...pageProps} />
    </NoSSRWrapper>
  )
}
