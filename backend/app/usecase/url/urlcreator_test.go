// +build !integration all

package url

import (
	"short/app/entity"
	"short/app/usecase/keygen"
	"short/app/usecase/repository"
	"short/app/usecase/validator"
	"testing"
	"time"

	"github.com/byliuyang/app/mdtest"
)

func TestURLCreatorPersist_CreateURL(t *testing.T) {
	t.Parallel()

	now := time.Now()

	alias := "220uFicCJj"
	longAlias := "an-alias-cannot-be-used-to-specify-default-arguments"

	testCases := []struct {
		name          string
		urls          urlMap
		alias         *string
		availableKeys []string
		user          entity.User
		url           entity.URL
		relationUsers []entity.User
		relationURLs  []entity.URL
		isPublic      bool
		expHasErr     bool
		expectedURL   entity.URL
	}{
		{
			name: "alias exists",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &now,
				},
			},
			alias: &alias,
			user: entity.User{
				Email: "alpha@example.com",
			},
			url:       entity.URL{},
			isPublic:  false,
			expHasErr: true,
		},
		{
			name: "alias too long",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &now,
				},
			},
			alias: &longAlias,
			user: entity.User{
				Email: "alpha@example.com",
			},
			url: entity.URL{
				OriginalURL: "https://www.google.com",
			},
			expHasErr: true,
		},
		{
			name:  "create alias successfully",
			urls:  urlMap{},
			alias: &alias,
			user: entity.User{
				Email: "alpha@example.com",
			},
			url: entity.URL{
				Alias:       "220uFicCJj",
				OriginalURL: "https://www.google.com",
				ExpireAt:    &now,
			},
			isPublic:  false,
			expHasErr: false,
			expectedURL: entity.URL{
				Alias:       "220uFicCJj",
				OriginalURL: "https://www.google.com",
				ExpireAt:    &now,
			},
		},
		{
			name: "automatically generate alias",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &now,
				},
			},
			availableKeys: []string{
				"test",
			},
			alias: nil,
			user: entity.User{
				Email: "alpha@example.com",
			},
			url: entity.URL{
				OriginalURL: "https://www.google.com",
			},
			expHasErr: false,
			expectedURL: entity.URL{
				Alias:       "test",
				OriginalURL: "https://www.google.com",
			},
		},
		{
			name: "no available key",
			urls: urlMap{
				"220uFicCJj": entity.URL{
					Alias:    "220uFicCJj",
					ExpireAt: &now,
				},
			},
			availableKeys: []string{},
			alias:         nil,
			user: entity.User{
				Email: "alpha@example.com",
			},
			url: entity.URL{
				OriginalURL: "https://www.google.com",
			},
			expHasErr: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			urlRepo := repository.NewURLFake(testCase.urls)
			userURLRepo := repository.NewUserURLRepoFake(
				testCase.relationUsers,
				testCase.relationURLs,
			)
			keyGen := keygen.NewFake(testCase.availableKeys)
			longLinkValidator := validator.NewLongLink()
			aliasValidator := validator.NewCustomAlias()

			creator := NewCreatorPersist(
				&urlRepo,
				&userURLRepo,
				&keyGen,
				longLinkValidator,
				aliasValidator,
			)

			_, err := urlRepo.GetByAlias(testCase.url.Alias)
			mdtest.NotEqual(t, nil, err)

			isExist := userURLRepo.IsRelationExist(testCase.user, testCase.url)
			mdtest.Equal(t, false, isExist)

			url, err := creator.CreateURL(testCase.url, testCase.alias, testCase.user, testCase.isPublic)
			if testCase.expHasErr {
				mdtest.NotEqual(t, nil, err)

				_, err = urlRepo.GetByAlias(testCase.expectedURL.Alias)
				mdtest.NotEqual(t, nil, err)

				isExist := userURLRepo.IsRelationExist(testCase.user, testCase.expectedURL)
				mdtest.Equal(t, false, isExist)
				return
			}
			mdtest.Equal(t, nil, err)
			mdtest.Equal(t, testCase.expectedURL, url)

			savedURL, err := urlRepo.GetByAlias(testCase.expectedURL.Alias)
			mdtest.Equal(t, nil, err)
			mdtest.Equal(t, testCase.expectedURL, savedURL)

			isExist = userURLRepo.IsRelationExist(testCase.user, testCase.expectedURL)
			mdtest.Equal(t, true, isExist)
		})
	}
}
