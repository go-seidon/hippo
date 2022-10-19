package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-seidon/hippo/internal/hashing"
	"github.com/go-seidon/hippo/internal/repository"
	"github.com/go-seidon/provider/encoding"
)

type BasicAuth interface {
	ParseAuthToken(ctx context.Context, p ParseAuthTokenParam) (*ParseAuthTokenResult, error)
	CheckCredential(ctx context.Context, p CheckCredentialParam) (*CheckCredentialResult, error)
}

type CheckCredentialParam struct {
	AuthToken string
}

type CheckCredentialResult struct {
	TokenValid bool
}

func (r *CheckCredentialResult) IsValid() bool {
	return r.TokenValid
}

type ParseAuthTokenParam struct {
	Token string
}

type ParseAuthTokenResult struct {
	ClientId     string
	ClientSecret string
}

type basicAuth struct {
	oAuthRepo repository.AuthRepository
	encoder   encoding.Encoder
	hasher    hashing.Hasher
}

func (a *basicAuth) ParseAuthToken(ctx context.Context, p ParseAuthTokenParam) (*ParseAuthTokenResult, error) {
	if strings.TrimSpace(p.Token) == "" {
		return nil, fmt.Errorf("invalid token")
	}

	d, err := a.encoder.Decode(p.Token)
	if err != nil {
		return nil, err
	}

	auth := strings.Split(string(d), ":")
	if len(auth) != 2 {
		return nil, fmt.Errorf("invalid auth encoding")
	}

	if strings.TrimSpace(auth[0]) == "" {
		return nil, fmt.Errorf("invalid client id")
	}
	if strings.TrimSpace(auth[1]) == "" {
		return nil, fmt.Errorf("invalid client secret")
	}

	res := &ParseAuthTokenResult{
		ClientId:     auth[0],
		ClientSecret: auth[1],
	}
	return res, nil
}

func (a *basicAuth) CheckCredential(ctx context.Context, p CheckCredentialParam) (*CheckCredentialResult, error) {

	client, err := a.ParseAuthToken(ctx, ParseAuthTokenParam{
		Token: p.AuthToken,
	})
	if err != nil {
		return nil, err
	}

	res := &CheckCredentialResult{
		TokenValid: false,
	}
	oClient, err := a.oAuthRepo.FindClient(ctx, repository.FindClientParam{
		ClientId: client.ClientId,
	})
	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			return res, nil
		}
		return nil, err
	}

	err = a.hasher.Verify(oClient.ClientSecret, client.ClientSecret)
	if err != nil {
		return res, nil
	}

	res.TokenValid = true
	return res, nil
}

type NewBasicAuthParam struct {
	AuthRepo repository.AuthRepository
	Encoder  encoding.Encoder
	Hasher   hashing.Hasher
}

func NewBasicAuth(p NewBasicAuthParam) (*basicAuth, error) {
	if p.AuthRepo == nil {
		return nil, fmt.Errorf("auth repo is not specified")
	}
	if p.Encoder == nil {
		return nil, fmt.Errorf("encoder is not specified")
	}
	if p.Hasher == nil {
		return nil, fmt.Errorf("hasher is not specified")
	}

	a := &basicAuth{
		oAuthRepo: p.AuthRepo,
		encoder:   p.Encoder,
		hasher:    p.Hasher,
	}
	return a, nil
}
