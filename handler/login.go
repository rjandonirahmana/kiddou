package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"kiddou/base"
	"kiddou/domain"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type LoginSSOandler struct {
	configGoogle *oauth2.Config
	configGithub *oauth2.Config
	usecaseUser  domain.UsecaseUser
}

func NewLoginSSO(configGoogle *oauth2.Config, configGithub *oauth2.Config, usecaseUSer domain.UsecaseUser) *LoginSSOandler {
	return &LoginSSOandler{configGoogle: configGoogle, configGithub: configGithub, usecaseUser: usecaseUSer}
}

func (h *LoginSSOandler) LoginGoogle(c *gin.Context) {
	URL, err := url.Parse(h.configGoogle.Endpoint.AuthURL)
	if err != nil {
		base.APIResponse(c, "UNAUTHORIZED", http.StatusUnauthorized, err.Error(), nil)
		return
	}
	log.Println(URL.String())
	parameters := url.Values{}
	parameters.Add("client_id", h.configGoogle.ClientID)
	parameters.Add("scope", strings.Join(h.configGoogle.Scopes, " "))
	parameters.Add("redirect_uri", h.configGoogle.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", "cobacoba")
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	log.Println(url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *LoginSSOandler) CallbackGoogleLogin(c *gin.Context) {
	state := c.Request.FormValue("state")
	log.Println(state)
	if state != "cobacoba" {
		log.Println("invalid oauth state, expected cobacoba" + ", got " + state + "\n")
		c.Redirect(http.StatusTemporaryRedirect, "/api/v1/home")
		return
	}
	code := c.Request.FormValue("code")
	log.Printf("ini code %s", code)

	if code == "" {
		log.Printf("code not foundd")
		reason := c.Request.FormValue("error_reason")
		if reason == "user_denied" {
			base.APIResponse(c, "UNAUTHORIZED", http.StatusUnauthorized, "User has denied Permission..", nil)
			return
		}
		c.JSON(403, reason)
		return
		// User has denied access..
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		token, err := h.configGoogle.Exchange(context.Background(), code)
		if err != nil {
			base.APIResponse(c, "UNAUTHORIZED", http.StatusUnauthorized, err.Error(), nil)
			return
		}
		log.Println("TOKEN>> AccessToken>> " + token.AccessToken)
		log.Println("TOKEN>> Expiration Time>> " + token.Expiry.String())
		log.Println("TOKEN>> RefreshToken>> " + token.RefreshToken)

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/api/v1/home")
			return

		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// logger.Log.Error("ReadAll: " + err.Error() + "\n")
			c.Redirect(http.StatusTemporaryRedirect, "/api/v1/home")
			return
		}
		var user map[string]interface{}
		err = json.Unmarshal(response, &user)
		if err != nil {
			base.APIResponse(c, "Failed to unmarshal", http.StatusUnauthorized, err.Error(), nil)
			return
		}
		log.Printf("ini map user %v", user)
		users := &domain.Users{
			Name:      user["given_name"].(string) + user["family_name"].(string),
			Email:     user["email"].(string),
			Avatar:    sql.NullString{String: user["picture"].(string), Valid: true},
			Password:  "google",
			Salt:      "google",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		tokenStr, err := h.usecaseUser.LoginGoogle(c, users, user["id"].(string))
		if err != nil {
			base.APIResponse(c, "Failed to login google", http.StatusUnauthorized, err.Error(), nil)
			return
		}

		base.ResponseAPIToken(c, "success login google", 200, "success", nil, tokenStr)
		return
	}
}

func (h *LoginSSOandler) HomeLogin(c *gin.Context) {

	c.Writer.WriteString(`<!DOCTYPE html>
	<html>
		<head>
			<title>OAuth-2 Test</title>
		</head>
		<body>
			<h2>OAuth-2 Test</h2>
			<p>
				Login with the following,
			</p>
			<ul>
				<li><a href="/login-google">Google</a></li>
				<li><a href="/login-fb">Facebook</a></li>
				<li><a href="/login-github">Github</a></li>
			</ul>
		</body>
	</html>`)

}
