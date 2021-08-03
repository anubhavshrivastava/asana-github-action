package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/razorpay/asana-github-action/httpclient"
)

type ICore interface {
	GetPR(ctx context.Context, prLink string) (*PR, error)
	UpdatePR(ctx context.Context, updatePrRequest UpdatePR) (*PR, error)
}

type impl struct{}

var core ICore

func SetCore() {
	core = &impl{}
}

func GetCore() ICore {
	return core
}

func (i impl) GetPR(ctx context.Context, prLink string) (*PR, error) {
	// this should call github and get the PR and fill in the details required
	pr := PR{}
	err := httpclient.Call(prLink, nil, &pr, headers(), http.MethodGet)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (i impl) UpdatePR(ctx context.Context, updatePrRequest UpdatePR) (*PR, error) {
	// for now this will just blindly update things
	// this is the one that will be used for updating the body with all the task details
	return nil, nil
}

func headers() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", GetConfig().AccessToken),
		"Content-Type":  JsonContentType,
	}
}
