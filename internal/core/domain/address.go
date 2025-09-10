package domain

import "errors"

type Address struct {
	Street  string `json:"street" validate:"required" gorm:"not null"`
	ZipCode string `json:"zipcode" validate:"required" gorm:"not null"`
	ExtNum  string `json:"ext_num" validate:"required" gorm:"not null"`
	IntNum  string `json:"int_num,omitempty"`
	City    string `json:"city" validate:"required" gorm:"not null"`
	State   string `json:"state" validate:"required" gorm:"not null"`
	Country string `json:"country" validate:"required" gorm:"not null"`
}

func NewAddress(street, zipCode, extNum, intNum, city, state, country string) (*Address, error) {
	addr := &Address{
		Street:  street,
		ZipCode: zipCode,
		ExtNum:  extNum,
		IntNum:  intNum,
		City:    city,
		State:   state,
		Country: country,
	}

	if err := addr.Validate(); err != nil {
		return nil, err
	}

	return addr, nil
}

func (a *Address) Validate() error {
	if a.Street == "" || a.ZipCode == "" || a.ExtNum == "" {
		return errors.New("street, zipcode and ext_num are required")
	}
	if a.City == "" || a.State == "" || a.Country == "" {
		return errors.New("city, state and country are required")
	}
	return nil
}

func (a *Address) FullAddress() string {
	full := a.Street + " " + a.ExtNum
	if a.IntNum != "" {
		full += " Int. " + a.IntNum
	}
	full += ", " + a.City + ", " + a.State + " " + a.ZipCode + ", " + a.Country
	return full
}
