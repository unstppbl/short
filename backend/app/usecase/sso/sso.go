package sso

import (
	"errors"
	"short/app/entity"
	"short/app/usecase/account"
	"short/app/usecase/auth"
	"short/app/usecase/service"
)

// SingleSignOn enables sign in through external identity providers, such as
// Github, Facebook, and Google.
type SingleSignOn struct {
	identityProvider  service.IdentityProvider
	ssoAccountService service.SSOAccount
	accountProvider   account.Provider
	authenticator     auth.Authenticator
}

// SignIn generates access token for a user using authorization code obtained
// from external identity provider.
func (o SingleSignOn) SignIn(authorizationCode string) (string, error) {
	if len(authorizationCode) < 1 {
		return "", errors.New("authorizationCode can't be empty")
	}

	accessToken, err := o.identityProvider.RequestAccessToken(authorizationCode)
	if err != nil {
		return "", err
	}

	ssoUser, err := o.ssoAccountService.GetSingleSignOnUser(accessToken)
	if err != nil {
		return "", err
	}

	email := ssoUser.Email
	isExist, err := o.accountProvider.IsAccountExist(email)
	if err != nil {
		return "", err
	}

	user := entity.User{
		Email: email,
	}
	authToken, err := o.authenticator.GenerateToken(user)
	if err != nil {
		return "", err
	}

	if isExist {
		return authToken, nil
	}

	err = o.accountProvider.CreateAccount(email, ssoUser.Name)
	if err != nil {
		return "", nil
	}

	return authToken, nil
}

// NewSingleSignOn creates SingleSignOn service for a given external
// identity provider.
func NewSingleSignOn(
	identityProvider service.IdentityProvider,
	ssoAccountService service.SSOAccount,
	accountProvider account.Provider,
	authenticator auth.Authenticator,
) SingleSignOn {
	return SingleSignOn{
		identityProvider:  identityProvider,
		ssoAccountService: ssoAccountService,
		accountProvider:   accountProvider,
		authenticator:     authenticator,
	}
}
