// Copyright 2020 Jeffrey Weitz \&lt;jeffdoubleyou@gmail.com\&gt;
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"encoding/json"
	"reflect"
)

type Address struct {
	AddressLine1       string
	AddressLine2       string
	AddressLine3       string
	City               string
	StateProv          string
	PostalCode         string
	PostalCodeExtended string
	AddressType        string
	CountryCode        string
	Classification     string
}

type AddressKeyFormat struct {
	ConsigneeName       string
	BuildingName        string
	AddressLine         []string
	Region              string
	PoliticalDivision1  string
	PoliticalDivision2  string
	PostcodePrimaryLow  string
	PostcodeExtendedLow string
	Urbanization        string
	CountryCode         string
}

type AddressClassification struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}

type Candidate struct {
	AddressClassification *AddressClassification `json:"AddressClassification"`
	AddressKeyFormat      *AddressKeyFormat      `json:"AddressKeyFormat"`
}

type XAVRequest struct {
	AddressKeyFormat *AddressKeyFormat `json:"AddressKeyFormat"`
}

type XAVResponseObject struct {
	Response                  *Response              `json:"Response"`
	NoCandidatesIndicator     bool                   `json:"NoCandidatesIndicator"`
	ValidAddressIndicator     bool                   `json:"ValidAddressIndicator"`
	AmbiguousAddressIndicator bool                   `json:"AmbiguousAddressIndicator"`
	AddressClassification     *AddressClassification `json:"AddressClassification"`
	Candidate                 []*Candidate           `json:"Candidate"`
}

type XAVResponse struct {
	XAVResponse *XAVResponseObject `json:"XAVResponse"`
}

// Convert the weird UPS response to something usable
// For example, the "Indicator" flags are either not present = false or an empty string = true
// And Candidates may or may not be an array
func (x *XAVResponseObject) UnmarshalJSON(b []byte) error {
	var decoded map[string]interface{}
	err := json.Unmarshal(b, &decoded)

	if err != nil {
		return err
	}

	v := reflect.ValueOf(*x)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if decoded[typeOfS.Field(i).Name] == nil {
			switch typeOfS.Field(i).Name {
			case "AddressClassification":
				decoded["AddressClassification"] = &AddressClassification{"0", "Unknown"}
			default:
				if typeOfS.Field(i).Type.String() == "bool" {
					decoded[typeOfS.Field(i).Name] = false
				} else {
					decoded[typeOfS.Field(i).Name] = ""
				}
			}
		} else {
			if typeOfS.Field(i).Type.String() == "bool" {
				decoded[typeOfS.Field(i).Name] = true
			}
		}
	}

	// Convert candidates to slice if necessary
	if decoded["Candidate"] != nil {
		if reflect.ValueOf(decoded["Candidate"]).Type().String() == "map[string]interface {}" {
			candidate := []interface{}{decoded["Candidate"]}
			decoded["Candidate"] = candidate
		}
	}

	jsonCandidate, _ := json.Marshal(decoded["Candidate"])
	json.Unmarshal(jsonCandidate, &x.Candidate)

	jsonResponse, _ := json.Marshal(decoded["Response"])
	json.Unmarshal(jsonResponse, &x.Response)

	jsonAddressClassification, _ := json.Marshal(decoded["AddressClassification"])
	json.Unmarshal(jsonAddressClassification, &x.AddressClassification)

	x.NoCandidatesIndicator = decoded["NoCandidatesIndicator"].(bool)
	x.ValidAddressIndicator = decoded["ValidAddressIndicator"].(bool)
	x.AmbiguousAddressIndicator = decoded["AmbiguousAddressIndicator"].(bool)

	return nil
}

// All of this to handle the poor JSON response of either an aray or string for the AddressLine property
// There might be / probably is a better way to do this
func (x *AddressKeyFormat) UnmarshalJSON(b []byte) error {
	var decoded map[string]interface{}

	err := json.Unmarshal(b, &decoded)

	if err != nil {
		return err
	}

	v := reflect.ValueOf(*x)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if decoded[typeOfS.Field(i).Name] == nil {
			decoded[typeOfS.Field(i).Name] = ""
		}
	}

	if decoded["AddressLine"] != nil {
		switch decoded["AddressLine"].(type) {
		case string:
			decoded["AddressLine"] = []string{decoded["AddressLine"].(string)}
		default:
			var s []string
			for _, x := range decoded["AddressLine"].([]interface{}) {
				s = append(s, x.(string))
			}
			decoded["AddressLine"] = s
		}
		x.AddressLine = decoded["AddressLine"].([]string)
	}

	x.ConsigneeName = decoded["ConsigneeName"].(string)
	x.BuildingName = decoded["BuildingName"].(string)
	x.Region = decoded["Region"].(string)
	x.PoliticalDivision1 = decoded["PoliticalDivision1"].(string)
	x.PoliticalDivision2 = decoded["PoliticalDivision2"].(string)
	x.PostcodePrimaryLow = decoded["PostcodePrimaryLow"].(string)
	x.PostcodeExtendedLow = decoded["PostcodeExtendedLow"].(string)
	x.Urbanization = decoded["Urbanization"].(string)
	x.CountryCode = decoded["CountryCode"].(string)

	return nil
}
