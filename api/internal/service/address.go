package service

import (
	"github.com/mislu/market-api/internal/db"
	"github.com/mislu/market-api/internal/types/exceptions"
	"github.com/mislu/market-api/internal/types/models"
	"github.com/mislu/market-api/internal/types/request"
	"github.com/mislu/market-api/internal/types/response"
	"gorm.io/gorm"
)

func CreateAddress(req *request.CreateAddressReq) exceptions.APIError {
	userAddress := models.UserAddress{
		UserID:    req.UserID,
		IsDefault: *req.IsDefault,
		Phone:     req.Phone,
		Receiver:  req.Receiver,
		Detail:    req.Detail,
	}
	err := db.WithTransaction(func(tx *gorm.DB) error {
		address := models.Address{
			Address:      req.Address,
			City:         req.City,
			District:     req.District,
			Province:     req.Province,
			Street:       req.Street,
			StreetNumber: req.StreetNumber,
			Latitude:     req.Latitude,
			Longitude:    req.Longitude,
		}

		if err := db.Create(&address, tx); err != nil {
			return err
		}

		userAddress.AddressID = address.ID
		if err := db.Create(&userAddress, tx); err != nil {
			return err
		}

		if !*req.IsDefault {
			return nil
		}

		defaultAddress, err := db.GetOne[models.UserAddress](
			db.Equal("user_id", req.UserID),
			db.Equal("is_default", true),
		)

		if err != nil {
			return err
		}

		if !defaultAddress.Exists() {
			return nil
		}

		defaultAddress.IsDefault = false
		if err := db.Update(&defaultAddress, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func UpdateAddress(req *request.UpdateAddressReq) exceptions.APIError {
	err := db.WithTransaction(func(tx *gorm.DB) error {
		userAddress, err := db.GetOne[models.UserAddress](
			db.Equal("address_id", req.AddressID),
		)

		if err != nil {
			return err
		}

		if !userAddress.Exists() {
			return nil
		}

		userAddress.Phone = req.Phone
		userAddress.Receiver = req.Receiver
		userAddress.Detail = req.Detail

		if err := db.Update(&userAddress, tx); err != nil {
			return err
		}

		address, err := db.GetOne[models.Address](
			db.Equal("id", userAddress.AddressID),
		)

		if err != nil {
			return err
		}

		address.Address = req.Address
		address.City = req.City
		address.District = req.District
		address.Province = req.Province
		address.Street = req.Street
		address.StreetNumber = req.StreetNumber
		address.Latitude = req.Latitude
		address.Longitude = req.Longitude

		if err := db.Update(&address, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func DeleteAddress(req *request.DeleteAddressReq) exceptions.APIError {
	err := db.WithTransaction(func(tx *gorm.DB) error {
		userAddress, err := db.GetOne[models.UserAddress](
			db.Equal("address_id", req.AddressID),
		)

		if err != nil {
			return err
		}

		if !userAddress.Exists() {
			return nil
		}

		if err := db.Delete(&userAddress, tx); err != nil {
			return err
		}

		address, err := db.GetOne[models.Address](
			db.Equal("id", userAddress.AddressID),
		)

		if err != nil {
			return err
		}

		if err := db.Delete(&address, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}

func GetAddress(req *request.GetAddressReq) (response.GetAddressResp, exceptions.APIError) {
	var resp response.GetAddressResp
	userAddresses, err := db.GetAll[models.UserAddress](
		db.Equal("user_id", req.UserID),
	)

	if err != nil {
		return resp, exceptions.InternalServerError(err)
	}

	for _, userAddress := range userAddresses {
		address, err := db.GetOne[models.Address](
			db.Equal("id", userAddress.AddressID),
		)

		if err != nil {
			return resp, exceptions.InternalServerError(err)
		}

		resp.Addresses = append(resp.Addresses, response.WrappedAddress{
			UserAddress: userAddress,
			Address:     address,
		})
	}

	resp.Page = req.Page
	resp.Size = req.Size

	return resp, nil
}

func SetDefaultAddress(req *request.SetDefaultAddressReq) exceptions.APIError {
	err := db.WithTransaction(func(tx *gorm.DB) error {
		userAddress, err := db.GetOne[models.UserAddress](
			db.Equal("address_id", req.AddressID),
		)

		if err != nil {
			return err
		}

		if !userAddress.Exists() {
			return nil
		}

		userAddress.IsDefault = *req.IsDefault
		if err := db.Update(&userAddress, tx); err != nil {
			return err
		}

		if !*req.IsDefault {
			return nil
		}

		defaultAddress, err := db.GetOne[models.UserAddress](
			db.Equal("user_id", req.UserID),
		)

		if !defaultAddress.Exists() {
			return nil
		}

		defaultAddress.IsDefault = false
		if err := db.Update(&defaultAddress, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return exceptions.InternalServerError(err)
	}

	return nil
}
