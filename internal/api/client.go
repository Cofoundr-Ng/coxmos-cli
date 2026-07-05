package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL   string
	Token     string
	APIKey    string
	APISecret string
	HTTP      *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) authHeaders() http.Header {
	h := http.Header{}
	if c.Token != "" {
		h.Set("Authorization", "Bearer "+c.Token)
	}
	if c.APIKey != "" {
		h.Set("X-Api-Key", c.APIKey)
		h.Set("X-Api-Secret", c.APISecret)
	}
	return h
}

func (c *Client) Do(method, path string, body, out interface{}) (*http.Response, error) {
	return c.do(method, path, body, out)
}

func (c *Client) do(method, path string, body, out interface{}) (*http.Response, error) {
	u, _ := url.JoinPath(c.BaseURL, path)
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal: %w", err)
		}
		r = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, u, r)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	req.Header = c.authHeaders()
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	if res.StatusCode >= 400 {
		bodyB, _ := io.ReadAll(res.Body)
		res.Body.Close()
		return res, fmt.Errorf("%s %s: %s (%s)", method, path, res.Status, string(bodyB))
	}
	if out != nil {
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			return res, fmt.Errorf("decode: %w", err)
		}
	}
	return res, nil
}

// --- Auth ---

type LoginReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresIn *int   `json:"expires_in,omitempty"`
}

type LoginRes struct {
	Token string `json:"token"`
	User  struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

func (c *Client) Login(req LoginReq) (*LoginRes, error) {
	var res LoginRes
	_, err := c.do("POST", "/auth/login", req, &res)
	return &res, err
}

type RegisterDeviceCodeRes struct {
	ExpiresIn int `json:"expires_in"`
}

func (c *Client) RegisterDeviceCode(codeHash, codePrefix string) error {
	_, err := c.do("POST", "/auth/cli/register", map[string]string{"code_hash": codeHash, "code_prefix": codePrefix}, nil)
	return err
}

type PollDeviceCodeRes struct {
	Status    string `json:"status"`
	Token     string `json:"token,omitempty"`
	User      struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user,omitempty"`
}

func (c *Client) PollDeviceCode(code string) (*PollDeviceCodeRes, error) {
	var res PollDeviceCodeRes
	_, err := c.do("POST", "/auth/cli/poll", map[string]string{"code": code}, &res)
	return &res, err
}

type APIKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key,omitempty"`
	Secret      string `json:"secret,omitempty"`
	Service     string `json:"service"`
	ExpiresAt   string `json:"expires_at,omitempty"`
	LastUsedAt  string `json:"last_used_at,omitempty"`
	SecretHint  string `json:"secret_hint,omitempty"`
}

type CreateAPIKeyReq struct {
	Name     string `json:"name"`
	Service  string `json:"service"`
	ExpiresIn string `json:"expires_in,omitempty"`
}

func (c *Client) CreateAPIKey(req CreateAPIKeyReq) (*APIKey, error) {
	var res APIKey
	_, err := c.do("POST", "/auth/api-keys", req, &res)
	return &res, err
}

func (c *Client) ListAPIKeys() ([]APIKey, error) {
	var res []APIKey
	_, err := c.do("GET", "/auth/api-keys", nil, &res)
	return res, err
}

func (c *Client) DeleteAPIKey(id string) error {
	_, err := c.do("DELETE", "/auth/api-keys/"+id, nil, nil)
	return err
}

// --- Apps ---

type DeployReq struct {
	GitURL    string `json:"git_url"`
	Branch    string `json:"branch,omitempty"`
	Framework string `json:"framework,omitempty"`
}

type DeployRes struct {
	ID           string `json:"id"`
	Slug         string `json:"slug"`
	URL          string `json:"url"`
	Status       string `json:"status"`
	DeploymentID string `json:"deployment_id"`
}

type App struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	URL        string `json:"url"`
	Framework  string `json:"framework"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type AppStats struct {
	Slug    string `json:"slug"`
	Isolates int   `json:"isolates"`
	VMs     int    `json:"vms"`
}

func (c *Client) DeployApp(req DeployReq) (*DeployRes, error) {
	var res DeployRes
	_, err := c.do("POST", "/apps/deploy", req, &res)
	return &res, err
}

func (c *Client) ListApps() ([]App, error) {
	var stats []AppStats
	_, err := c.do("GET", "/apps/stats", nil, &stats)
	if err != nil {
		return nil, err
	}
	apps := make([]App, len(stats))
	for i, s := range stats {
		apps[i] = App{Slug: s.Slug, Status: "running"}
	}
	return apps, nil
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Line      string `json:"line"`
	Stream    string `json:"stream"`
}

func (c *Client) GetLogs(deploymentID string) ([]LogEntry, error) {
	var res []LogEntry
	_, err := c.do("GET", "/apps/logs?id="+deploymentID, nil, &res)
	return res, err
}

func (c *Client) StopApp(slug string) error {
	_, err := c.do("POST", "/apps/stop", map[string]string{"slug": slug}, nil)
	return err
}

func (c *Client) StartApp(slug string) error {
	_, err := c.do("POST", "/apps/start", map[string]string{"slug": slug}, nil)
	return err
}

func (c *Client) RestartApp(slug string) error {
	_, err := c.do("POST", "/apps/restart", map[string]string{"slug": slug}, nil)
	return err
}

// --- Databases ---

type CreateDBReq struct {
	DBType string `json:"db_type"`
	Kind   string `json:"kind"`
	Name   string `json:"name"`
}

type Database struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	DBType           string `json:"db_type"`
	Kind             string `json:"kind"`
	Status           string `json:"status"`
	ConnectionString string `json:"connection_string,omitempty"`
	Host             string `json:"host,omitempty"`
	Port             int    `json:"port,omitempty"`
}

func (c *Client) CreateDatabase(req CreateDBReq) (*Database, error) {
	var res Database
	_, err := c.do("POST", "/databases/new", req, &res)
	return &res, err
}

func (c *Client) ListDatabases() ([]Database, error) {
	var res []Database
	_, err := c.do("GET", "/databases", nil, &res)
	return res, err
}

func (c *Client) DeleteDatabase(id string) error {
	_, err := c.do("DELETE", "/databases/"+id, nil, nil)
	return err
}

// --- Redis ---

type CreateRedisReq struct {
	Name      string `json:"name"`
	MemoryMB  int    `json:"memory_mb,omitempty"`
}

type RedisInstance struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URI       string `json:"uri"`
	Status    string `json:"status"`
	MemoryMB  int    `json:"memory_mb"`
}

type RedisStats struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	MemoryMB  float64 `json:"memory_mb"`
	Status    string  `json:"status"`
}

func (c *Client) CreateRedis(req CreateRedisReq) (*RedisInstance, error) {
	var res RedisInstance
	_, err := c.do("POST", "/redis/new", req, &res)
	return &res, err
}

func (c *Client) ListRedis() ([]RedisStats, error) {
	var res []RedisStats
	_, err := c.do("GET", "/redis/stats", nil, &res)
	return res, err
}

func (c *Client) StopRedis(id string) error {
	_, err := c.do("POST", "/redis/stop", map[string]string{"id": id}, nil)
	return err
}

func (c *Client) StartRedis(id string) error {
	_, err := c.do("POST", "/redis/start", map[string]string{"id": id}, nil)
	return err
}

func (c *Client) RestartRedis(id string) error {
	_, err := c.do("POST", "/redis/restart", map[string]string{"id": id}, nil)
	return err
}

// --- Email ---

type CreateEmailAccountReq struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type EmailAccount struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type AddEmailDomainReq struct {
	Domain string `json:"domain"`
}

func (c *Client) CreateEmailAccount(req CreateEmailAccountReq) (*EmailAccount, error) {
	var res EmailAccount
	_, err := c.do("POST", "/mail/accounts", req, &res)
	return &res, err
}

func (c *Client) AddEmailDomain(domain string) error {
	_, err := c.do("POST", "/mail/domains", AddEmailDomainReq{Domain: domain}, nil)
	return err
}

func (c *Client) VerifyEmailDomain(domain string) error {
	_, err := c.do("POST", "/mail/domains/"+domain+"/verify", nil, nil)
	return err
}

// --- DNS ---

type RegisterDomainReq struct {
	Domain string `json:"domain"`
}

type DNSRecord struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}
type DNSRecordsRes struct {
	Records []DNSRecord `json:"records"`
}

type AttachDomainReq struct {
	Domain string `json:"domain"`
	AppSlug string `json:"app_slug"`
}

func (c *Client) RegisterDomain(domain string) error {
	_, err := c.do("POST", "/api/v1/dns/register", RegisterDomainReq{Domain: domain}, nil)
	return err
}

func (c *Client) RemoveDomain(domain string) error {
	_, err := c.do("DELETE", "/api/v1/dns/"+domain, nil, nil)
	return err
}

func (c *Client) VerifyDomain(domain string) error {
	_, err := c.do("POST", "/api/v1/dns/verify", RegisterDomainReq{Domain: domain}, nil)
	return err
}

func (c *Client) CheckDomainVerification(domain string) (bool, error) {
	var res map[string]bool
	_, err := c.do("POST", "/api/v1/dns/verify/check", RegisterDomainReq{Domain: domain}, &res)
	if err != nil {
		return false, err
	}
	return res["verified"], nil
}

func (c *Client) ListDNSRecords(domain string) ([]DNSRecord, error) {
	var res DNSRecordsRes
	_, err := c.do("GET", "/api/v1/dns/"+domain+"/records", nil, &res)
	return res.Records, err
}

func (c *Client) AddDKIMRecord(domain string) error {
	_, err := c.do("POST", "/api/v1/dns/dkim", RegisterDomainReq{Domain: domain}, nil)
	return err
}

func (c *Client) AttachDomain(domain, appSlug string) error {
	_, err := c.do("POST", "/api/v1/dns/domains/attach", AttachDomainReq{Domain: domain, AppSlug: appSlug}, nil)
	return err
}
