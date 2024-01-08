export interface Currency {
  s: string, //significand
  e: string, //exponent
  c: "IDR", //currency
}

export function createZeroCurrency(currencyUnit: "IDR"):Currency{
  return {s: "0", e: "1", c: currencyUnit}
}

export function displayCurrencyInUI(currency:Currency|null):string{
  if (currency === null) {
    return "-"
  }
  switch (currency.c) {
    case "IDR":
      return new Intl.NumberFormat('id-ID', {style: 'currency', currency: 'IDR'}).format(convertCurrencyToNumber(currency))
    default:
      return `${convertCurrencyToNumber(currency)}`
  }
}

export function convertCurrencyToNumber(currency:Currency):number{
  const s = parseFloat(currency.s)
  const e = parseFloat(currency.e)
  return s * Math.pow(10,e)
}

export function convertCurrencyToIntegerString(currency:Currency):string{
  return `${convertCurrencyToNumber(currency)}`
}

function onlyContainsZeroForNumber(s: string){
	for (let i = 0; i < s.length; i++) {
		if(["1","2","3","4","5","6","7","8","9"].includes(s[i])){
			return false
		}
	}
	return true
}

const indexNotFound = -1

export function convertNumberStringToCurrency(value:string,currencyUnit:"IDR"):Currency{
  var significand = ""
	var exponent = ""
	if (value === "" || value.startsWith(".")) {
		value = `0${value}`
	}
	let valueFloat = parseFloat(value)
	if (valueFloat === 0) {
		//value is zero
		significand = "0"
		exponent = "1"
	} else if (value.startsWith("0.") || value.startsWith("-0.")) {
		//value has negative exponent
		//(between -1 exclusive and 1 exclusive; also excluding 0)

		if (value.startsWith("-")) {
			significand += "-"
			value = value.slice(1)
		}
		value = value.slice(2)
		let leadingZero = 1 //the 0. is the first leading zero
		for (let i = 0; i < value.length; i++) {
			if (value[i] === '0') {
				leadingZero++
			} else {
        leadingZero += value.length - 1 - i
				significand += value.slice(i)
				break
			}
		}
		exponent = `-${leadingZero}`
	} else {
		//value has positive exponent
		var dotIndex = indexNotFound
		for (let i = 0; i < value.length; i++) {
			if (value[i] === '.') {
				dotIndex = i
				break
			}
		}
		if (dotIndex === indexNotFound) {
			//value is integer only
			value = `${value}.0`
		}
		for (let i = 0; i < value.length; i++) {
			if (value[i] === '.') {
				exponent = `${i - 1}`
				break
			}
		}
		let valueWithoutDot = value.split(".").join("")
		significand = `${valueWithoutDot.slice(0,1)}.${valueWithoutDot.slice(1)}`
	}

	//trim trailing zeros after comma from significand
	let commaSpotted = false
	let firstZeroAfterCommaIndex = indexNotFound
	for (let i = 0; i < significand.length; i++) {
		if (significand[i] === '.') {
			commaSpotted = true
			continue
		}

		if (commaSpotted) {
			if (significand[i] === '0' && firstZeroAfterCommaIndex === indexNotFound) {
				firstZeroAfterCommaIndex = i
			} else if (significand[i] !== '0') {
				firstZeroAfterCommaIndex = indexNotFound
			}
		}

		if (i === (significand.length - 1) /*last index*/ && commaSpotted && firstZeroAfterCommaIndex !== indexNotFound) {
			significand = significand.slice(0,firstZeroAfterCommaIndex)
		}
	}
	if (significand[significand.length-1] === '.') { //if last char is .
		significand = significand + "0"
	}
  return {
    s: significand,
    e: exponent,
    c: currencyUnit,
  }
}
