package hyperledger

import (
	"go-auth-api/internal/jwt"
	"log"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

type loginData struct {
	User        string `json:"user"`
	Org         string `json:"org"`
	IdentityBts string `json:"identity_bts"`
}

type Client interface {
	CreateUserAndLogin(user, secret, org, affiliation string) (token string, err error)
	// Helpers
	FindOrCreateAffilation(org, affiliation string) (*msp.AffiliationResponse, error)
	CreateIdentity(user, secret, org, affiliation, identityType string) (*msp.IdentityResponse, error)
	LoginAndGetToken(user, secret, org string) (token string, err error)
}

func NewClient(configPath string) (Client, error) {
	sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		log.Printf("unable to get config from file err: %s\n", err.Error())
		return nil, errors.WithMessage(err, "unable to initialize hyperledger client")
	}

	return &client{sdk: sdk}, nil
}

type client struct {
	sdk *fabsdk.FabricSDK
}

func (c *client) CreateUserAndLogin(user, secret, org, affiliation string) (token string, err error) {
	if _, err := c.FindOrCreateAffilation(org, affiliation); err != nil {
		return "", err
	}

	if _, err := c.CreateIdentity(user, secret, org, affiliation, "user"); err != nil {
		return "", err
	}

	return c.LoginAndGetToken(user, secret, org)
}

func (c *client) FindOrCreateAffilation(org, affiliation string) (*msp.AffiliationResponse, error) {
	mspclient, err := msp.New(c.sdk.Context(), msp.WithOrg(org))
	if err != nil {
		return nil, err
	}

	res, err := mspclient.GetAffiliation(affiliation)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return mspclient.AddAffiliation(&msp.AffiliationRequest{
				Name:  affiliation,
				Force: true,
			})
		}
		return nil, err
	}

	return res, nil
}

func (c *client) CreateIdentity(user, secret, org, affiliation, identityType string) (*msp.IdentityResponse, error) {
	mspclient, err := msp.New(c.sdk.Context(), msp.WithOrg(org))
	if err != nil {
		return nil, err
	}

	return mspclient.CreateIdentity(&msp.IdentityRequest{
		ID:          user,
		Affiliation: affiliation,
		Type:        identityType,
		Secret:      secret,
	})
}

func (c *client) LoginAndGetToken(user, secret, org string) (token string, err error) {
	mspclient, err := msp.New(c.sdk.Context(), msp.WithOrg(org))
	if err != nil {
		return "", err
	}

	if err = mspclient.Enroll(user, msp.WithSecret(secret)); err != nil {
		return "", err
	}

	signingIdentity, err := mspclient.GetSigningIdentity(user)
	if err != nil {
		return "", err
	}

	bts, err := signingIdentity.Serialize()
	if err != nil {
		return "", err
	}

	return jwt.GetToken(user, org, bts), nil
}
