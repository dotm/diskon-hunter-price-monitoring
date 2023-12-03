package currencyutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

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
				significand += value[i:]
				break
			}
		}
		exponent = fmt.Sprintf("-%d", leadingZero)
	} else {
		//value has positive exponent
		var dotIndex = -1
		for i := 0; i < len(value); i++ {
			if value[i] == '.' {
				dotIndex = i
				break
			}
		}
		if dotIndex == -1 {
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
	firstZeroAfterCommaIndex := -1
	for i := 0; i < len(significand); i++ {
		if significand[i] == '.' {
			commaSpotted = true
			continue
		}

		if commaSpotted {
			if significand[i] == '0' && firstZeroAfterCommaIndex == -1 {
				firstZeroAfterCommaIndex = i
			} else if significand[i] != '0' {
				firstZeroAfterCommaIndex = -1
			}
		}

		if i == len(significand)-1 /*last index*/ && firstZeroAfterCommaIndex != 1 {
			significand = significand[:firstZeroAfterCommaIndex]
		}
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
