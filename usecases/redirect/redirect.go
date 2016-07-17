package redirect

import (
	//"github.com/pkg/errors"
	//	"fmt"
	. "github.com/sjhitchner/trackly/domain"
	//	"log"
)

type RedirectInteractorImpl struct {
}

func NewRedirectInteractor() RedirectInteractor {
	return &RedirectInteractorImpl{}
}

func (t *RedirectInteractorImpl) Redirect(ShortURL, DeviceId, *DeviceInfo) (LongURL, error) {
	return "", nil
}

func (t *RedirectInteractorImpl) Register(url LongURL) (ShortURL, error) {
	return "", nil
}
