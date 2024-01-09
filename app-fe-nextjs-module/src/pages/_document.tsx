import { Head, Html, Main, NextScript } from 'next/document'

export default function Document() {
  return (
    <Html lang="en" className='h-full'>
      <Head />
      <body className='h-full bg-black max-w-screen-2xl mx-auto'>
        <Main />
        <NextScript />
      </body>
    </Html>
  )
}
