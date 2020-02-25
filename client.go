package cas

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/tmsong/hlog"
	"net/http"
	"net/url"
)

// Options : Client configuration options
type Options struct {
	LoginURL       *url.URL     // LoginURL to the CAS login
	BaseUrl        *url.URL     // LoginURL to the CAS service
	Store          TicketStore  // Custom TicketStore, if nil a MemoryStore will be used
	Client         *http.Client // Custom http client to allow options for http connections
	SendService    bool         // Custom sendService to determine whether you need to send service param
	URLScheme      URLScheme    // Custom url scheme, can be used to modify the request urls for the client
	Cookie         *http.Cookie // http.Cookie options, uses Path, Domain, MaxAge, HttpOnly, & Secure
	SessionStore   SessionStore
	ValidationType string //CAS1,CAS2,CAS3
	AppId          int64
	AppKey         string
	ClientHost     *url.URL
	OpenUrl        *url.URL
	Logger         *hlog.Logger
}

type cli struct {
	tickets   TicketStore
	client    *http.Client
	urlScheme URLScheme
	cookie    *http.Cookie

	sessions    SessionStore
	sendService bool

	stValidator *ServiceTicketValidator
	pmValidator *PermissionValidator

	appId      int64
	appKey     string
	clientHost *url.URL
}

// Client implements the main protocol
type Client struct {
	*cli
	logger *hlog.Logger
}

// NewClient creates a Client with the provided Options.
func NewClient(options *Options) *Client {
	var tickets TicketStore
	if options.Store != nil {
		tickets = options.Store
	} else {
		tickets = &MemoryStore{}
	}

	var sessions SessionStore
	if options.SessionStore != nil {
		sessions = options.SessionStore
	} else {
		sessions = NewMemorySessionStore()
	}

	var urlScheme URLScheme
	if options.URLScheme != nil {
		urlScheme = options.URLScheme
	} else {
		urlScheme = NewDefaultURLScheme(options.LoginURL)
	}

	var client *http.Client
	if options.Client != nil {
		client = options.Client
	} else {
		client = &http.Client{}
	}

	var cookie *http.Cookie
	if options.Cookie != nil {
		cookie = options.Cookie
	} else {
		cookie = &http.Cookie{
			MaxAge:   86400,
			HttpOnly: false,
			Secure:   false,
		}
	}
	c := &Client{logger: options.Logger}
	c.cli = &cli{
		tickets:     tickets,
		client:      client,
		urlScheme:   urlScheme,
		cookie:      cookie,
		sessions:    sessions,
		sendService: options.SendService,
		stValidator: NewServiceTicketValidator(client, options.BaseUrl, options.ValidationType, c),
		pmValidator: NewPermissionValidator(client, options.OpenUrl, c),
		appId:       options.AppId,
		appKey:      options.AppKey,
		clientHost:  options.ClientHost,
	}
	return c
}

// Handle wraps a http.Handler to provide CAS authentication for the handler.
func (c *Client) Handle(h http.Handler) http.Handler {
	return &clientHandler{
		c: c,
		h: h,
	}
}

// HandleFunc wraps a function to provide CAS authentication for the handler function.
func (c *Client) HandleFunc(h func(http.ResponseWriter, *http.Request)) http.Handler {
	return c.Handle(http.HandlerFunc(h))
}

// requestURL determines an absolute LoginURL from the http.Request.
func requestURL(r *http.Request) (*url.URL, error) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		return nil, err
	}

	u.Host = r.Host
	if host := r.Header.Get("X-Forwarded-Host"); host != "" {
		u.Host = host
	}

	u.Scheme = "http"
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		u.Scheme = scheme
	} else if r.TLS != nil {
		u.Scheme = "https"
	}

	return u, nil
}

// LoginUrlForRequest determines the CAS login LoginURL for the http.Request.
func (c *Client) LoginUrlForRequest(r *http.Request) (string, error) {
	u, err := c.urlScheme.Index()
	if err != nil {
		return "", err
	}

	//service, err := requestURL(r)
	//if err != nil {
	//	return "", err
	//}

	q := u.Query()
	//q.Add("service", c.clientHost)
	//q.Add("redirectUrl", c.clientHost.String())
	jumpTo, err := requestURL(r)
	if ref := r.Referer(); len(ref) > 0 {
		jumpTo, err = url.Parse(ref)
	}
	if err == nil {
		q.Add("jumpto", jumpTo.String())
	}
	q.Add("appId", fmt.Sprintf("%d", c.appId))

	u.RawQuery = q.Encode()

	return u.String(), nil
}

// LogoutUrlForRequest determines the CAS logout LoginURL for the http.Request.
func (c *Client) LogoutUrlForRequest(r *http.Request) (string, error) {
	u, err := c.urlScheme.Logout()
	if err != nil {
		return "", err
	}

	if c.sendService {
		service, err := requestURL(r)
		if err != nil {
			return "", err
		}

		q := u.Query()
		q.Add("service", sanitisedURLString(service))
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

// performSingleLogout processes a single logout request
func (c *Client) PerformSingleLogout(w http.ResponseWriter, r *http.Request) {
	rawXML := r.FormValue("logoutRequest")
	logoutRequest, err := parseLogoutRequest([]byte(rawXML))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.tickets.Delete(logoutRequest.SessionIndex); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//c.deleteSession(logoutRequest.SessionIndex)	wrong use, this 'session index' is not in session store

	w.WriteHeader(http.StatusOK)
}

// ServiceValidateUrlForRequest determines the CAS serviceValidate LoginURL for the ticket and http.Request.
func (c *Client) ServiceValidateUrlForRequest(ticket string, r *http.Request) (string, error) {
	service, err := requestURL(r)
	if err != nil {
		return "", err
	}
	return c.stValidator.ServiceValidateUrl(service, ticket)
}

// ValidateUrlForRequest determines the CAS validate LoginURL for the ticket and http.Request.
func (c *Client) ValidateUrlForRequest(ticket string, r *http.Request) (string, error) {
	service, err := requestURL(r)
	if err != nil {
		return "", err
	}
	return c.stValidator.ValidateUrl1(service, ticket)
}

// RedirectToLogout replies to the request with a redirect LoginURL to log out of CAS.
func (c *Client) RedirectToLogout(w http.ResponseWriter, r *http.Request) {
	u, err := c.LogoutUrlForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.logger.Debugf("Logging out, redirecting client to %v with status %v",
		u, http.StatusFound)
	c.clearSession(w, r)
	http.Redirect(w, r, u, http.StatusFound)
}

// RedirectToLogin replies to the request with a redirect LoginURL to authenticate with CAS.
func (c *Client) RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	u, err := c.LoginUrlForRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.logger.Debugf("Redirecting client to %v with status %v", u, http.StatusFound)
	http.Redirect(w, r, u, http.StatusFound)
}

// validateTicket performs CAS ticket validation with the given ticket and service.
func (c *Client) validateTicket(ticket string, service *http.Request) error {
	//serviceURL, err := requestURL(service)
	//if err != nil {
	//	return err
	//}
	success, err := c.stValidator.ValidateTicket(c.clientHost, ticket)
	if err != nil {
		return err
	}

	if err := c.tickets.Write(ticket, success); err != nil {
		return err
	}

	return nil
}

// getSession finds or creates a session for the request.
//
// A cookie is set on the response if one is not provided with the request.
// Validates the ticket if the LoginURL parameter is provided.
func (c *Client) GetSession(w http.ResponseWriter, r *http.Request) {
	cookie := c.getCookie(w, r)

	if s, ok := c.sessions.Get(cookie.Value); ok {
		if t, err := c.tickets.Read(s); err == nil {
			c.logger.Debugf("Re-used ticket %s for %s", s, t.User)

			c.setSession(cookie.Value, s) //每次有操作都刷新过期时间
			c.tickets.Refresh(s)          //刷新ticket过期时间
			setAuthenticationResponse(r, t)
			return
		} else {
			c.logger.Warningf("Ticket %v not in %T: %v", s, c.tickets, err)

			c.logger.Infof("Clearing ticket %s, no longer exists in ticket store", s)

			clearCookie(w, cookie)
		}
	}

	if ticket := r.URL.Query().Get("ticket"); ticket != "" {
		if err := c.validateTicket(ticket, r); err != nil {
			c.logger.Errorf("Error validating ticket: %v", err)
			return // allow ServeHTTP()
		}

		c.setSession(cookie.Value, ticket)

		if t, err := c.tickets.Read(ticket); err == nil {
			c.logger.Debugf("Validated ticket %s for %s", ticket, t.User)

			setAuthenticationResponse(r, t)
			if jumpto := r.URL.Query().Get("jumpto"); jumpto != "" {
				http.Redirect(w, r, jumpto, http.StatusFound)
			}
			return
		} else {
			c.logger.Warningf("Ticket %v not in %T: %v", ticket, c.tickets, err)

			c.logger.Infof("Clearing ticket %s, no longer exists in ticket store", ticket)

			clearCookie(w, cookie)
		}
	}
}

// getCookie finds or creates the session cookie on the response.
func (c *Client) getCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		// NOTE: Intentionally not enabling HttpOnly so the cookie can
		//       still be used by Ajax requests.
		cookie = &http.Cookie{
			Name:     sessionCookieName,
			Value:    newSessionID(),
			Path:     c.cookie.Path,
			Domain:   c.cookie.Domain,
			MaxAge:   c.cookie.MaxAge,
			HttpOnly: c.cookie.HttpOnly,
			Secure:   c.cookie.Secure,
		}

		c.logger.Debugf("Setting %v cookie with value: %v", cookie.Name, cookie.Value)

		r.AddCookie(cookie) // so we can find it later if required
		http.SetCookie(w, cookie)
	}

	return cookie
}

func (c *Client) GetLogger() *hlog.Logger {
	return c.logger
}

func (c *Client) SetSessionStore(store SessionStore) {
	c.sessions = store
}

func (c *Client) SetTicketStore(store TicketStore) {
	c.tickets = store
}

// newSessionId generates a new opaque session identifier for use in the cookie.
func newSessionID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// generate 64 character string
	bytes := make([]byte, 64)
	rand.Read(bytes)

	for k, v := range bytes {
		bytes[k] = alphabet[v%byte(len(alphabet))]
	}

	return string(bytes)
}

// clearCookie invalidates and removes the cookie from the client.
func clearCookie(w http.ResponseWriter, c *http.Cookie) {
	c.MaxAge = -1
	http.SetCookie(w, c)
}

// setSession stores the session id to ticket mapping in the Client.
func (c *Client) setSession(id string, ticket string) {
	c.logger.Debugf("Recording session, %v -> %v", id, ticket)

	c.sessions.Set(id, ticket)
}

// clearSession removes the session from the client and clears the cookie.
func (c *Client) clearSession(w http.ResponseWriter, r *http.Request) {
	cookie := c.getCookie(w, r)

	if s, ok := c.sessions.Get(cookie.Value); ok {
		if err := c.tickets.Delete(s); err != nil {
			c.logger.Errorf("Failed to remove %v from %T: %v", cookie.Value, c.tickets, err)
		}

		c.deleteSession(s)
	}

	clearCookie(w, cookie)
}

// deleteSession removes the session from the client
func (c *Client) deleteSession(id string) {
	c.sessions.Delete(id)
}

// 判断是否有权限 from the client
func (c *Client) PermissionValidateForRequest(r *http.Request) error {
	service, err := requestURL(r)
	if err != nil {
		return err
	}
	uid := GetCurrentUserId(r)
	if uid <= 0 {
		return errors.New("please login")
	}
	return c.pmValidator.HasPermission(uid, service.Path)
}

// 获取角色列表 from the client
func (c *Client) RoleList(r *http.Request) ([]RoleListResponse, error) {
	uid := GetCurrentUserId(r)
	if uid <= 0 {
		return nil, errors.New("please login")
	}
	return c.pmValidator.RoleList(uid)
}

// 获取功能列表 from the client
func (c *Client) PermissionList(r *http.Request, roleId int64) ([]PermissionListResponse, error) {
	uid := GetCurrentUserId(r)
	if uid <= 0 {
		return nil, errors.New("please login")
	}
	return c.pmValidator.PermissionList(uid, roleId)
}

// 获取用户信息 from the client
func (c *Client) UserInfo(userId int64) (*UserInfoResponse, error) {
	return c.pmValidator.UserInfo(userId)
}

// 获取用户信息详情 from the client
func (c *Client) UserInfoDetail(userId int64, employeeId string) (*UserInfoDetailResponse, error) {
	return c.pmValidator.UserInfoDetail(userId, employeeId)
}

// 获取部门信息详情 from the client
func (c *Client) DepartmentInfo(departmentId int64) (*DepartmentInfoRespose, error) {
	return c.pmValidator.DepartmentInfo(departmentId)
}
