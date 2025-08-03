package server

import (
	"errors"
	"strings"

	pb "wallet-services/proto"
)

// ValidateGetWalletByUserIDRequest validates the GetWalletByUserID request
func ValidateGetWalletByUserIDRequest(req *pb.GetWalletByUserIDRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("userID is required")
	}
	return nil
}

// ValidateCreateWalletRequest validates the CreateWallet request
func ValidateCreateWalletRequest(req *pb.CreateWalletRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("userID is required")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return errors.New("currency is required")
	}
	if !ValidateCurrency(req.Currency) {
		return errors.New("invalid currency code")
	}
	return nil
}

// ValidateFundWalletRequest validates the FundWallet request
func ValidateFundWalletRequest(req *pb.FundWalletRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("userID is required")
	}
	if strings.TrimSpace(req.Amount) == "" {
		return errors.New("amount is required")
	}
	if !ValidateAmount(req.Amount) {
		return errors.New("invalid amount format or amount must be positive")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return errors.New("currency is required")
	}
	if !ValidateCurrency(req.Currency) {
		return errors.New("invalid currency code")
	}
	return nil
}

// ValidateDeductFromWalletRequest validates the DeductFromWallet request
func ValidateDeductFromWalletRequest(req *pb.DeductFromWalletRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("userID is required")
	}
	if strings.TrimSpace(req.Amount) == "" {
		return errors.New("amount is required")
	}
	if !ValidateAmount(req.Amount) {
		return errors.New("invalid amount format or amount must be positive")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return errors.New("currency is required")
	}
	if !ValidateCurrency(req.Currency) {
		return errors.New("invalid currency code")
	}
	return nil
}

// ValidateGetWalletBalanceRequest validates the GetWalletBalance request
func ValidateGetWalletBalanceRequest(req *pb.GetWalletBalanceRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("userID is required")
	}
	return nil
}

// ValidateGetWalletTransactionsRequest validates the GetWalletTransactions request
func ValidateGetWalletTransactionsRequest(req *pb.GetWalletTransactionsRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("userID is required")
	}
	return nil
}

// ValidateRefundToWalletRequest validates the RefundToWallet request
func ValidateRefundToWalletRequest(req *pb.RefundToWalletRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if strings.TrimSpace(req.TransactionID) == "" {
		return errors.New("transactionID is required")
	}
	if strings.TrimSpace(req.Amount) == "" {
		return errors.New("amount is required")
	}
	if !ValidateAmount(req.Amount) {
		return errors.New("invalid amount format or amount must be positive")
	}
	return nil
}