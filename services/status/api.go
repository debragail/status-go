package status

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

var (
	// ErrInvalidStatusLoginParameters is returned when the number of parameters
	// for status_login is not valid.
	ErrInvalidStatusLoginParameters = errors.New("invalid number of parameters for status_login (2 expected)")

	// ErrInvalidStatusJoinPublicChannelParameters is returned when the number
	// of parameters for status_login is not valid.
	ErrInvalidStatusJoinPublicChannelParameters = errors.New("invalid number of parameters for status_joinpublicchannel (1 expected)")

	// ErrInvalidStatusSignupParameters is returned when the number of parameters
	// for status_signup is not valid.
	ErrInvalidStatusSignupParameters = errors.New("invalid number of parameters for status_signup (1 expected)")

	// ErrCouldNotCreateAnAccount is returned when an account based on a password
	// can't be created
	ErrCouldNotCreateAnAccount = errors.New("could not create the specified account")

	// ErrCouldNotJoinPublicChannel is returned when an account based on a password
	// can't be created
	ErrCouldNotJoinPublicChannel = errors.New("could not join the specified public channel")
)

// PublicAPI represents a set of APIs from the `web3.status` namespace.
type PublicAPI struct {
	s *Service
}

// NewAPI creates an instance of the status API.
func NewAPI(s *Service) *PublicAPI {
	return &PublicAPI{s: s}
}

// LoginRequest : json request for status_login.
type LoginRequest struct {
	Addr string `json:"address"`
	Pwd  string `json:"password"`
}

// LoginResponse : json response returned by status_login.
type LoginResponse struct {
	AccountKey   *keystore.Key `json:"-"`
	AddressKeyID string        `json:"address_key_id"`
}

// Login is an implementation of `status_login` or `web3.status.login` API
func (api *PublicAPI) Login(context context.Context, req LoginRequest) (LoginResponse, error) {
	var err error
	res := LoginResponse{}

	if _, res.AccountKey, err = api.s.am.AddressToDecryptedAccount(req.Addr, req.Pwd); err != nil {
		return res, err
	}

	if res.AddressKeyID, err = api.s.w.AddKeyPair(res.AccountKey.PrivateKey); err != nil {
		return res, err
	}

	if err = api.s.am.SelectAccount(req.Addr, req.Pwd); err != nil {
		return res, err
	}

	return res, err
}

// SignupRequest : json request for status_signup.
type SignupRequest struct {
	Password string `json:"password"`
}

// SignupResponse : json response returned by status_signup.
type SignupResponse struct {
	Address  string `json:"address"`
	Pubkey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic"`
}

// Signup is an implementation of `status_signup` or `web3.status.signup` API
func (api *PublicAPI) Signup(context context.Context, req SignupRequest) (SignupResponse, error) {
	var err error
	var res SignupResponse

	if res.Address, res.Pubkey, res.Mnemonic, err = api.s.am.CreateAccount(req.Password); err != nil {
		return res, ErrCouldNotCreateAnAccount
	}

	return res, nil
}
