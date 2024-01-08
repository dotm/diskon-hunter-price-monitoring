package currencyutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// add IDR, USD as CurrencyUnit enum later ~kodok
const IDR = "IDR"
const USD = "USD"

type Currency struct {
	Significand  string `json:"s"`
	Exponent     string `json:"e"`
	CurrencyUnit string `json:"c"`
}

func NewZeroCurrency(currencyUnit string) Currency {
	return Currency{
		Significand:  "0",
		Exponent:     "1",
		CurrencyUnit: currencyUnit,
	}
}

const indexNotFound = -1

func NewFromExcelString(value, currencyUnit string) Currency {
	value = strings.Replace(value, ",", ".", -1) //replace all comma with dot
	return NewFromNumberString(value, currencyUnit)
}

func NewFromNumberString(value, currencyUnit string) Currency {
	var significand = ""
	var exponent = ""
	if value == "" || strings.HasPrefix(value, ".") {
		value = fmt.Sprintf("0%s", value)
	}
	valueFloat, _ := strconv.ParseFloat(value, 64)
	if valueFloat == 0 {
		//value is zero
		significand = "0"
		exponent = "1"
	} else if strings.HasPrefix(value, "0.") || strings.HasPrefix(value, "-0.") {
		//value has negative exponent
		//(between -1 exclusive and 1 exclusive; also excluding 0)

		if strings.HasPrefix(value, "-") {
			significand += "-"
			value = value[1:]
		}
		value = value[2:]
		leadingZero := 1 //the 0. is the first leading zero
		for i := 0; i < len(value); i++ {
			if value[i] == '0' {
				leadingZero++
			} else {
				leadingZero += len(value) - 1 - i
				significand += value[i:]
				break
			}
		}
		exponent = fmt.Sprintf("-%d", leadingZero)
	} else {
		//value has positive exponent
		var dotIndex = indexNotFound
		for i := 0; i < len(value); i++ {
			if value[i] == '.' {
				dotIndex = i
				break
			}
		}
		if dotIndex == indexNotFound {
			//value is integer only
			value = fmt.Sprintf("%s.0", value)
		}
		for i := 0; i < len(value); i++ {
			if value[i] == '.' {
				exponent = strconv.Itoa(i - 1)
				break
			}
		}
		valueWithoutDot := strings.Join(strings.Split(value, "."), "")
		significand = fmt.Sprintf("%s.%s", valueWithoutDot[0:1], valueWithoutDot[1:])
	}

	//trim trailing zeros after comma from significand
	commaSpotted := false
	firstZeroAfterCommaIndex := indexNotFound
	for i := 0; i < len(significand); i++ {
		if significand[i] == '.' {
			commaSpotted = true
			continue
		}

		if commaSpotted {
			if significand[i] == '0' && firstZeroAfterCommaIndex == indexNotFound {
				firstZeroAfterCommaIndex = i
			} else if significand[i] != '0' {
				firstZeroAfterCommaIndex = indexNotFound
			}
		}

		if i == len(significand)-1 /*last index*/ && commaSpotted && firstZeroAfterCommaIndex != indexNotFound {
			significand = significand[:firstZeroAfterCommaIndex]
		}
	}
	if significand[len(significand)-1] == '.' { //if last char is .
		significand = significand + "0"
	}

	return Currency{
		Significand:  significand,
		Exponent:     exponent,
		CurrencyUnit: currencyUnit,
	}
}

func (x Currency) ToDouble() float64 {
	s, _ := strconv.ParseFloat(x.Significand, 64)
	e, _ := strconv.ParseFloat(x.Exponent, 64)
	return s * math.Pow(10, e)
}

func (x Currency) Add(y Currency) Currency {
	return NewFromNumberString(fmt.Sprintf("%f", x.ToDouble()+y.ToDouble()), x.CurrencyUnit)
}

func (x Currency) Substract(y Currency) Currency {
	//validate unit is the same, implement manually (not using double operation) ~kodok
	return NewFromNumberString(fmt.Sprintf("%f", x.ToDouble()-y.ToDouble()), x.CurrencyUnit)
}

func (x Currency) IsNotZero() bool {
	return !x.IsZero()
}
func (x Currency) IsZero() bool {
	for i := 0; i < len(x.Significand); i++ {
		if x.Significand[i] != '0' && x.Significand[i] != '.' && x.Significand[i] != '-' {
			return false
		}
	}
	return true
}

func (x Currency) IsNegative() bool {
	return strings.HasPrefix(x.Significand, "-")
}

func (x Currency) IsPositive() bool {
	return !x.IsNegative()
}

func (x Currency) IsLessThanOrEqualTo(b Currency) bool {
	return x.IsLessThan(b) || x.IsEqualTo(b)
}

func (x Currency) IsLessThan(b Currency) bool {
	//TODO: implement directly using significand and exponent
	//and then copy to Currency ~kodok
	return x.Substract(b).IsNegative()
}

func (x Currency) IsEqualTo(b Currency) bool {
	return x.ToDouble() == b.ToDouble()

	//the code below hasn't account for when significand is different (2.0 and 2).
	//this can happen when the currency is loaded from JSON, thus bypassing NewFromNumberString.
	// return significand == b.significand &&
	//     exponent == b.exponent &&
	//     unit == b.unit;
}
