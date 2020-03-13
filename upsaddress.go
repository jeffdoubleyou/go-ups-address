package UPSAddress

import (
	"encoding/json"
	"fmt"
	"go-ups-address/UPSClient"
	"go-ups-address/models"
	"strconv"
)

type AddressAPI struct {
	client               *UPSClient.Client
	maximumCandidateSize int
	requestOption        int
}

type AddressValidationResult struct {
	XAVResponse *models.XAVResponse
}

func NewUPSAddress(username, password, accessKey string) *AddressAPI {
	return &AddressAPI{UPSClient.NewClient(username, password, accessKey), 10, 3}
}

func (a *AddressAPI) MaximumCandidateSize(size ...int) int {
	if len(size) == 1 {
		if size[0] > 50 {
			size[0] = 50
		}
		a.maximumCandidateSize = size[0]
	}
	return a.maximumCandidateSize
}

func (a *AddressAPI) Debug(debug ...bool) bool {
	if len(debug) == 1 {
		a.client.Debug(debug[0])
	}
	return a.client.Debug()
}

func (a *AddressAPI) Timeout(timeout ...int) int {
	if len(timeout) == 1 {
		a.client.Timeout(timeout[0])
	}
	return a.client.Timeout()
}

func (a *AddressAPI) ValidateAddress(address *models.Address) (result *AddressValidationResult, err error) {
	addressKeyFormat := &models.AddressKeyFormat{
		AddressLine:         []string{address.AddressLine1, address.AddressLine2, address.AddressLine3},
		PoliticalDivision1:  address.StateProv,
		PoliticalDivision2:  address.City,
		PostcodePrimaryLow:  address.PostalCode,
		PostcodeExtendedLow: address.PostalCodeExtended,
		CountryCode:         address.CountryCode,
	}

	xavRequest := &models.XAVRequest{addressKeyFormat}

	res, err := a.client.Post("addressvalidation", "XAVRequest", xavRequest, strconv.Itoa(a.requestOption))

	if err != nil {
		return
	}

	var response *models.XAVResponse
	err = json.Unmarshal(res, &response)

	if err != nil {
		return
	}

	// Request failure
	if response.XAVResponse.Response.ResponseStatus.Code != "1" {
		return nil, fmt.Errorf("Unable to call address validation service: %s", response.XAVResponse.Response.ResponseStatus.Description)
	}

	result = &AddressValidationResult{response}
	return
}

func (r *AddressValidationResult) XAVResponseObject() *models.XAVResponseObject {
	return r.XAVResponse.XAVResponse
}

func (r *AddressValidationResult) ValidAddress() bool {
	return r.XAVResponse.XAVResponse.ValidAddressIndicator
}

func (r *AddressValidationResult) AmbiguousAddress() bool {
	return r.XAVResponse.XAVResponse.AmbiguousAddressIndicator
}

func (r *AddressValidationResult) NoCandidate() bool {
	return r.XAVResponse.XAVResponse.NoCandidatesIndicator
}

func (r *AddressValidationResult) AddressClassification() string {
	return r.XAVResponse.XAVResponse.AddressClassification.Description
}

// If the address is valid, this will return the formatted address in the models.Address format.  If the address was not valid, this will return the first candidate in the models.Address format.  In the event that you call this without a valid address, nil will be returned
func (r *AddressValidationResult) Address() *models.Address {
	if len(r.XAVResponse.XAVResponse.Candidate) > 0 {
		candidate := r.XAVResponse.XAVResponse.Candidate[0].AddressKeyFormat
		address := &models.Address{
			AddressLine1:       candidate.AddressLine[0],
			City:               candidate.PoliticalDivision2,
			StateProv:          candidate.PoliticalDivision1,
			PostalCode:         candidate.PostcodePrimaryLow,
			PostalCodeExtended: candidate.PostcodeExtendedLow,
			CountryCode:        candidate.CountryCode,
		}

		if len(candidate.AddressLine) >= 2 {
			address.AddressLine2 = candidate.AddressLine[1]
		}
		if len(candidate.AddressLine) >= 3 {
			address.AddressLine3 = candidate.AddressLine[2]
		}
		return address
	}
	return nil

}

// This will return the first candidate's address key format.  It should be used to either take the first recommended address, or retrieve the valid address as formatted by UPS in the Address Key Format
func (r *AddressValidationResult) AddressKeyFormat() *models.AddressKeyFormat {
	if len(r.XAVResponse.XAVResponse.Candidate) > 0 {
		return r.XAVResponse.XAVResponse.Candidate[0].AddressKeyFormat
	}
	return nil
}
