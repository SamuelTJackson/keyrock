package keyrock

import (
	"fmt"
	"github.com/rs/xid"
	"net/http"
	"time"
)

// retrieve keyrock Token, keyrock refresh_token and user information
func (c HTTPClient) CallbackHandler(r *http.Request) (*Token, error) {
	queryState := r.FormValue("state")
	if queryState == "" {
		return nil, fmt.Errorf("state is empty")
	}
	if !c.ValidateState(queryState) {
		return nil, fmt.Errorf("state is not Valid")
	}
	expiresAt := r.FormValue("expires_at")
	if len(expiresAt) == 0 {
		return nil, fmt.Errorf("expires_at is empty")
	}
	token := r.FormValue("token")
	if len(token) == 0 {
		return nil, fmt.Errorf("token is empty")
	}
	return &Token{
		Value:     token,
		ExpiresAt: time.Time{},
	}, nil
}

func (c HTTPClient) RedirectToKeyrock(w http.ResponseWriter, r *http.Request, appID ID) error {
	state := xid.New().String()
	url := fmt.Sprintf("/oauth2/authorize?response_type=token&client_id=%s&redirect_uri=%s&state=%s",
		appID.Value, c.RedirectURL, state)
	if err := c.SaveState(state); err != nil {
		return err
	}
	http.Redirect(w,r,fmt.Sprintf("%s%s",c.KeyrockBaseURL, url), http.StatusFound)
	return nil
}
