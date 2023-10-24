package model

import (
	"time"
)

type TokenPlacement string

const (
	TokenPlacementHeader TokenPlacement = "header"
	TokenPlacementCookie TokenPlacement = "cookie"
)


type General struct {
	ID                      uint          `json:"id"`
	Name                    string        `json:"name"`
	CanLogin                bool          `json:"canLogin"`
	CanRegister             bool          `json:"canRegister"`
	HttpOnlyCookie          bool          `json:"httpOnlyCookie"`
	AccessTokenExpiryTime   time.Duration `json:"accessTokenExpireTime"`
	RefreshTokenExpiryTime  time.Duration `json:"refreshTokenExpireTime"`
	OrganizationEmailDomain string        `json:"OrganizationEmailDomain"`
	TokenPlacement          TokenPlacement `json:"tokenPlacement"`
}
