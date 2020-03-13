# go-ups-address
Golang library for validating addresses with UPS API

## Usage

```
package main

import (
	"fmt"

	UPSAddress "github.com/jeffdoubleyou/go-ups-address"
)

func main() {
	fmt.Println("vim-go")
	ups := UPSAddress.NewUPSAddress("yourusername", "yourpassword", "0DBCDEFGHIJKL")
	address := ups.NewAddress(map[string]string{
		"AddressLine1": "200 Corporate Pointe",
		"AddressLine2": "Suite 350",
		"City":         "Culver City",
		"StateProv":    "CA",
		"PostalCode":   "90230",
		"CountryCode":  "US",
	})

	ups.Debug(true)

	v, err := ups.ValidateAddress(address)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if v.ValidAddress() {
			fmt.Printf("The address is valid\n")
			fmt.Printf("This is a %s address\n", v.AddressClassification())
		} else {
			if v.AmbiguousAddress() {
				firstCandidate := v.Address()
				fmt.Printf("Maybe try street address: %s\n", firstCandidate.AddressLine1)
			}
		}
	}

}
```
