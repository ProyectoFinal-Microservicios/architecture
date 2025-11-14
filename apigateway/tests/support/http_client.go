package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient especializado para pruebas BDD
type HTTPClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewHTTPClient crea un nuevo cliente HTTP
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBaseURL retorna la URL base del cliente
func (c *HTTPClient) GetBaseURL() string {
	return c.baseURL
}

// SetToken establece el token JWT para autenticación
func (c *HTTPClient) SetToken(token string) {
	c.token = token
}

// GetToken obtiene el token actual
func (c *HTTPClient) GetToken() string {
	return c.token
}

// Response encapsula la respuesta HTTP
type Response struct {
	StatusCode int
	Body       []byte
	Data       interface{}
	Headers    http.Header
}

// Request realiza una solicitud HTTP genérica
func (c *HTTPClient) Request(method, endpoint string, body interface{}, headers map[string]string) (*Response, error) {
	url := c.baseURL + endpoint

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Headers por defecto
	req.Header.Set("Content-Type", "application/json")

	// Agregar token si está disponible
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	// Headers personalizados
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}

	// Intentar parsear JSON
	if len(respBody) > 0 {
		var data interface{}
		if err := json.Unmarshal(respBody, &data); err == nil {
			response.Data = data
		}
	}

	return response, nil
}

// GET realiza un GET request
func (c *HTTPClient) GET(endpoint string, headers map[string]string) (*Response, error) {
	return c.Request("GET", endpoint, nil, headers)
}

// POST realiza un POST request
func (c *HTTPClient) POST(endpoint string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request("POST", endpoint, body, headers)
}

// PATCH realiza un PATCH request
func (c *HTTPClient) PATCH(endpoint string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request("PATCH", endpoint, body, headers)
}

// PUT realiza un PUT request
func (c *HTTPClient) PUT(endpoint string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request("PUT", endpoint, body, headers)
}

// DELETE realiza un DELETE request
func (c *HTTPClient) DELETE(endpoint string, headers map[string]string) (*Response, error) {
	return c.Request("DELETE", endpoint, nil, headers)
}

// GetBodyAsJSON retorna el body parseado como JSON
func (r *Response) GetBodyAsJSON() map[string]interface{} {
	if r.Data == nil {
		return make(map[string]interface{})
	}
	if data, ok := r.Data.(map[string]interface{}); ok {
		return data
	}
	return make(map[string]interface{})
}

// GetBodyAsString retorna el body como string
func (r *Response) GetBodyAsString() string {
	return string(r.Body)
}

// HasHeader verifica si existe un header
func (r *Response) HasHeader(name string) bool {
	_, ok := r.Headers[name]
	return ok
}

// GetHeader obtiene el valor de un header
func (r *Response) GetHeader(name string) string {
	return r.Headers.Get(name)
}