package url

import (
	"short/app/entity"
	"short/app/usecase/keygen"
	"short/app/usecase/repository"
	"short/app/usecase/validator"
)

var _ Creator = (*CreatorPersist)(nil)

// ErrAliasExist represents alias unavailable error
type ErrAliasExist string

func (e ErrAliasExist) Error() string {
	return string(e)
}

// ErrInvalidLongLink represents incorrect long link format error
type ErrInvalidLongLink string

func (e ErrInvalidLongLink) Error() string {
	return string(e)
}

// ErrInvalidCustomAlias represents incorrect custom alias format error
type ErrInvalidCustomAlias string

func (e ErrInvalidCustomAlias) Error() string {
	return string(e)
}

// Creator represents a URL alias creator
type Creator interface {
	CreateURL(url entity.URL, alias *string, user entity.User, isPublic bool) (entity.URL, error)
}

// CreatorPersist represents a URL alias creator which persist the generated
// alias in the repository
type CreatorPersist struct {
	urlRepo             repository.URL
	userURLRelationRepo repository.UserURLRelation
	keyGen              keygen.KeyGenerator
	longLinkValidator   validator.LongLink
	aliasValidator      validator.CustomAlias
}

// CreateURL persists a new url with a given or auto generated alias in the repository.
// TODO(issue#235): add functionality for public URLs
func (c CreatorPersist) CreateURL(url entity.URL, customAlias *string, user entity.User, isPublic bool) (entity.URL, error) {
	longLink := url.OriginalURL
	if !c.longLinkValidator.IsValid(&longLink) {
		return entity.URL{}, ErrInvalidLongLink(longLink)
	}

	if customAlias == nil {
		return c.createURLWithAutoAlias(url, user)
	}

	if !c.aliasValidator.IsValid(customAlias) {
		return entity.URL{}, ErrInvalidCustomAlias(*customAlias)
	}
	return c.createURLWithCustomAlias(url, *customAlias, user)
}

func (c CreatorPersist) createURLWithAutoAlias(url entity.URL, user entity.User) (entity.URL, error) {
	key, err := c.keyGen.NewKey()
	if err != nil {
		return entity.URL{}, err
	}
	randomAlias := string(key)
	return c.createURLWithCustomAlias(url, randomAlias, user)
}

func (c CreatorPersist) createURLWithCustomAlias(url entity.URL, alias string, user entity.User) (entity.URL, error) {
	url.Alias = alias

	isExist, err := c.urlRepo.IsAliasExist(alias)
	if err != nil {
		return entity.URL{}, err
	}

	if isExist {
		return entity.URL{}, ErrAliasExist("url alias already exist")
	}

	err = c.urlRepo.Create(url)
	if err != nil {
		return entity.URL{}, err
	}

	err = c.userURLRelationRepo.CreateRelation(user, url)
	return url, err
}

// NewCreatorPersist creates CreatorPersist
func NewCreatorPersist(
	urlRepo repository.URL,
	userURLRelationRepo repository.UserURLRelation,
	keyGen keygen.KeyGenerator,
	longLinkValidator validator.LongLink,
	aliasValidator validator.CustomAlias,
) CreatorPersist {
	return CreatorPersist{
		urlRepo:             urlRepo,
		userURLRelationRepo: userURLRelationRepo,
		keyGen:              keyGen,
		longLinkValidator:   longLinkValidator,
		aliasValidator:      aliasValidator,
	}
}
