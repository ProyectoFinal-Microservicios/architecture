package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

var apiContext *APIGatewayContext

// APIGatewayContext contiene el contexto compartido de pruebas
type APIGatewayContext struct {
	baseURL       string
	client        *http.Client
	lastResponse  *http.Response
	lastResponseBody []byte
	lastError     error
	customHeaders map[string]string
	token         string
	username      string
}

// NewAPIGatewayContext crea un nuevo contexto
func NewAPIGatewayContext() *APIGatewayContext {
	return &APIGatewayContext{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		customHeaders: make(map[string]string),
	}
}

// ============ IMPLEMENTACIÓN DE STEPS REALES ============

func queElGatewayEstDisponibleEn(baseURL string) error {
	apiContext = NewAPIGatewayContext()
	apiContext.baseURL = baseURL
	// Verificar que el gateway esté realmente disponible
	resp, err := apiContext.client.Get(baseURL + "/health")
	if err != nil {
		return fmt.Errorf("el gateway no está disponible en %s: %v", baseURL, err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return fmt.Errorf("el gateway respondió con estado %d", resp.StatusCode)
	}
	
	return nil
}

func queElServicioDeAutenticacinEstDisponible() error {
	// En una implementación real, podrías verificar el health del servicio de autenticación
	// Por ahora asumimos que está disponible si el gateway responde
	return nil
}

func queElServicioDeProfilesEstDisponible() error {
	// Similar al anterior, asumimos disponibilidad
	return nil
}

func queExisteUnUsuarioConUsername(username string) error {
	apiContext.username = username
	// En una implementación real, podrías crear el usuario o verificar que existe
	return nil
}

func queEstoyAutenticadoComo(username string) error {
	apiContext.username = username
	// Simulamos un token válido - en pruebas reales obtendrías un token real
	apiContext.token = "simulated-valid-token-for-" + username
	apiContext.customHeaders["Authorization"] = "Bearer " + apiContext.token
	return nil
}

func incluyoTokenVlido() error {
	if apiContext.token == "" {
		apiContext.token = "simulated-valid-token"
	}
	apiContext.customHeaders["Authorization"] = "Bearer " + apiContext.token
	return nil
}

func noIncluyoToken() error {
	delete(apiContext.customHeaders, "Authorization")
	apiContext.token = ""
	return nil
}

func queTengoUnTokenExpirado() error {
	apiContext.customHeaders["Authorization"] = "Bearer expired-token-123"
	return nil
}

func hagoUnaSolicitudGETA(endpoint string) error {
	return makeRequest("GET", endpoint, nil)
}

func hagoUnaSolicitudGETASinToken(endpoint string) error {
	// Guardar token temporalmente
	tempToken := apiContext.token
	delete(apiContext.customHeaders, "Authorization")
	
	err := makeRequest("GET", endpoint, nil)
	
	// Restaurar token
	if tempToken != "" {
		apiContext.customHeaders["Authorization"] = "Bearer " + tempToken
	}
	return err
}

func hagoUnaSolicitudGETAConTokenVlido(endpoint string) error {
	incluyoTokenVlido()
	return makeRequest("GET", endpoint, nil)
}

func hagoUnaSolicitudGETAConTokenExpirado(endpoint string) error {
	queTengoUnTokenExpirado()
	return makeRequest("GET", endpoint, nil)
}

func hagoUnaSolicitudGETAConTokenInvlido(endpoint string) error {
	apiContext.customHeaders["Authorization"] = "Bearer invalid-token-123"
	return makeRequest("GET", endpoint, nil)
}

func hagoUnaSolicitudDELETEAConTokenVlido(endpoint string) error {
	incluyoTokenVlido()
	return makeRequest("DELETE", endpoint, nil)
}

func hagoUnaSolicitudDELETEASinToken(endpoint string) error {
	return hagoUnaSolicitudGETASinToken(endpoint)
}

func hagoUnaSolicitudPOSTAConDatos(endpoint string, dataTable *godog.Table) error {
	body := convertTableToJSON(dataTable)
	return makeRequest("POST", endpoint, body)
}

func hagoUnaSolicitudPATCHAConDatos(endpoint string, dataTable *godog.Table) error {
	body := convertTableToJSON(dataTable)
	return makeRequest("PATCH", endpoint, body)
}

func hagoUnaSolicitudPUTAConDatos(endpoint string, dataTable *godog.Table) error {
	body := convertTableToJSON(dataTable)
	return makeRequest("PUT", endpoint, body)
}

// ============ VALIDACIONES DE RESPUESTA ============

func laRespuestaDebeTenerEstado(expectedStatus int) error {
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	if apiContext.lastResponse.StatusCode != expectedStatus {
		bodyStr := string(apiContext.lastResponseBody)
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}
		return fmt.Errorf("se esperaba estado %d pero se obtuvo %d. Respuesta: %s", 
			expectedStatus, apiContext.lastResponse.StatusCode, bodyStr)
	}
	
	return nil
}

func laRespuestaDebeTenerEstadoO(status1, status2 int) error {
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	actualStatus := apiContext.lastResponse.StatusCode
	if actualStatus != status1 && actualStatus != status2 {
		return fmt.Errorf("se esperaba estado %d o %d pero se obtuvo %d", status1, status2, actualStatus)
	}
	
	return nil
}

func laRespuestaDebeSerJSONVlido() error {
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	if len(apiContext.lastResponseBody) == 0 {
		return fmt.Errorf("el cuerpo de la respuesta está vacío")
	}
	
	var data interface{}
	if err := json.Unmarshal(apiContext.lastResponseBody, &data); err != nil {
		return fmt.Errorf("la respuesta no es JSON válido: %v. Body: %s", err, string(apiContext.lastResponseBody))
	}
	
	return nil
}

func laRespuestaDebeContener(field string) error {
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal(apiContext.lastResponseBody, &data); err != nil {
		return fmt.Errorf("no se pudo parsear JSON: %v", err)
	}
	
	if _, exists := data[field]; !exists {
		return fmt.Errorf("la respuesta no contiene el campo '%s'", field)
	}
	
	return nil
}

func laRespuestaDebeContenerConValor(field, expectedValue string) error {
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal(apiContext.lastResponseBody, &data); err != nil {
		return fmt.Errorf("no se pudo parsear JSON: %v", err)
	}
	
	value, exists := data[field]
	if !exists {
		return fmt.Errorf("la respuesta no contiene el campo '%s'", field)
	}
	
	// Convertir el valor a string para comparar
	actualValue := fmt.Sprintf("%v", value)
	if actualValue != expectedValue {
		return fmt.Errorf("se esperaba '%s' = '%s' pero se obtuvo '%s'", field, expectedValue, actualValue)
	}
	
	return nil
}

func laRespuestaDebeContenerAccess_token() error {
	return laRespuestaDebeContener("access_token")
}

func laRespuestaDebeContenerUserData() error {
	return laRespuestaDebeContener("user")
}

func laRespuestaDebeContenerError() error {
	return laRespuestaDebeContener("error")
}

func laRespuestaDebeContenerErrorMsg(expectedError string) error {
	return laRespuestaDebeContenerConValor("error", expectedError)
}

func laRespuestaDebeContenerMessage(expectedMessage string) error {
	return laRespuestaDebeContenerConValor("message", expectedMessage)
}

func laRespuestaDebeContenerStatus(expectedStatus string) error {
	return laRespuestaDebeContenerConValor("status", expectedStatus)
}

func laRespuestaDebeContenerUsername(expectedUsername string) error {
	return laRespuestaDebeContenerConValor("username", expectedUsername)
}

func laRespuestaDebeContenerFirstName(expectedName string) error {
	// Buscar en user.firstName
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal(apiContext.lastResponseBody, &data); err != nil {
		return fmt.Errorf("no se pudo parsear JSON: %v", err)
	}
	
	user, exists := data["user"].(map[string]interface{})
	if !exists {
		return fmt.Errorf("la respuesta no contiene objeto 'user'")
	}
	
	firstName, exists := user["firstName"]
	if !exists {
		return fmt.Errorf("el usuario no contiene campo 'firstName'")
	}
	
	actualValue := fmt.Sprintf("%v", firstName)
	if actualValue != expectedName {
		return fmt.Errorf("se esperaba firstName = '%s' pero se obtuvo '%s'", expectedName, actualValue)
	}
	
	return nil
}

func elHeaderContentTypeDebeSer(expectedContentType string) error {
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	contentType := apiContext.lastResponse.Header.Get("Content-Type")
	if contentType != expectedContentType {
		return fmt.Errorf("se esperaba Content-Type '%s' pero se obtuvo '%s'", expectedContentType, contentType)
	}
	
	return nil
}

func elHeaderAccessControlAllowOriginDebeExistir() error {
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	acao := apiContext.lastResponse.Header.Get("Access-Control-Allow-Origin")
	if acao == "" {
		return fmt.Errorf("el header Access-Control-Allow-Origin no existe")
	}
	
	return nil
}

func elHeaderAccessControlAllowMethodsDebeIncluirGET() error {
	if apiContext.lastResponse == nil {
		return fmt.Errorf("no se recibió respuesta")
	}
	
	acam := apiContext.lastResponse.Header.Get("Access-Control-Allow-Methods")
	if acam == "" {
		return fmt.Errorf("el header Access-Control-Allow-Methods no existe")
	}
	
	if !strings.Contains(acam, "GET") {
		return fmt.Errorf("el header Access-Control-Allow-Methods no incluye GET: %s", acam)
	}
	
	return nil
}

// ============ FUNCIONES AUXILIARES ============

func makeRequest(method, endpoint string, body []byte) error {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, apiContext.baseURL+endpoint, reqBody)
	if err != nil {
		return err
	}

	// Headers por defecto
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Headers personalizados
	for key, value := range apiContext.customHeaders {
		req.Header.Set(key, value)
	}

	resp, err := apiContext.client.Do(req)
	apiContext.lastResponse = resp
	apiContext.lastError = err

	if err != nil {
		return err
	}

	// Leer el cuerpo de la respuesta
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	apiContext.lastResponseBody = responseBody

	return nil
}

func convertTableToJSON(table *godog.Table) []byte {
	if table == nil || len(table.Rows) == 0 {
		return []byte("{}")
	}

	data := make(map[string]interface{})
	headers := table.Rows[0].Cells

	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i].Cells
		for j, header := range headers {
			if j < len(row) {
				key := header.Value
				value := row[j].Value
				
				// Intentar convertir a número si es posible
				if intVal, err := strconv.Atoi(value); err == nil {
					data[key] = intVal
				} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
					data[key] = floatVal
				} else if value == "true" {
					data[key] = true
				} else if value == "false" {
					data[key] = false
				} else {
					data[key] = value
				}
			}
		}
	}

	jsonData, _ := json.Marshal(data)
	return jsonData
}

// ============ INICIALIZACIÓN DE SCENARIO ============

func InitializeScenario(ctx *godog.ScenarioContext, apiCtx *APIGatewayContext) {
	apiContext = apiCtx
	// Configurar todos los steps
	ctx.Step(`^que el gateway está disponible en "([^"]*)"$`, queElGatewayEstDisponibleEn)
	ctx.Step(`^que el servicio de autenticación está disponible$`, queElServicioDeAutenticacinEstDisponible)
	ctx.Step(`^que el servicio de profiles está disponible$`, queElServicioDeProfilesEstDisponible)
	ctx.Step(`^que existe un usuario con username "([^"]*)"$`, queExisteUnUsuarioConUsername)
	ctx.Step(`^que estoy autenticado como "([^"]*)"$`, queEstoyAutenticadoComo)
	ctx.Step(`^que tengo un token expirado$`, queTengoUnTokenExpirado)
	
	ctx.Step(`^hago una solicitud GET a "([^"]*)"$`, hagoUnaSolicitudGETA)
	ctx.Step(`^hago una solicitud GET a "([^"]*)" sin token$`, hagoUnaSolicitudGETASinToken)
	ctx.Step(`^hago una solicitud GET a "([^"]*)" con token válido$`, hagoUnaSolicitudGETAConTokenVlido)
	ctx.Step(`^hago una solicitud GET a "([^"]*)" con token expirado$`, hagoUnaSolicitudGETAConTokenExpirado)
	ctx.Step(`^hago una solicitud GET a "([^"]*)" con token inválido$`, hagoUnaSolicitudGETAConTokenInvlido)
	ctx.Step(`^hago una solicitud DELETE a "([^"]*)" con token válido$`, hagoUnaSolicitudDELETEAConTokenVlido)
	ctx.Step(`^hago una solicitud DELETE a "([^"]*)" sin token$`, hagoUnaSolicitudDELETEASinToken)
	ctx.Step(`^hago una solicitud POST a "([^"]*)" con datos:$`, hagoUnaSolicitudPOSTAConDatos)
	ctx.Step(`^hago una solicitud PATCH a "([^"]*)" con datos:$`, hagoUnaSolicitudPATCHAConDatos)
	ctx.Step(`^hago una solicitud PUT a "([^"]*)" con datos:$`, hagoUnaSolicitudPUTAConDatos)
	
	ctx.Step(`^incluyo token válido$`, incluyoTokenVlido)
	ctx.Step(`^no incluyo token$`, noIncluyoToken)
	
	ctx.Step(`^la respuesta debe tener estado (\d+)$`, laRespuestaDebeTenerEstado)
	ctx.Step(`^la respuesta debe tener estado (\d+) o (\d+)$`, laRespuestaDebeTenerEstadoO)
	ctx.Step(`^la respuesta debe ser JSON válido$`, laRespuestaDebeSerJSONVlido)
	ctx.Step(`^la respuesta debe contener "([^"]*)"$`, laRespuestaDebeContener)
	ctx.Step(`^la respuesta debe contener "([^"]*)" con valor "([^"]*)"$`, laRespuestaDebeContenerConValor)
	ctx.Step(`^la respuesta debe contener access_token$`, laRespuestaDebeContenerAccess_token)
	ctx.Step(`^la respuesta debe contener user data$`, laRespuestaDebeContenerUserData)
	ctx.Step(`^la respuesta debe contener error$`, laRespuestaDebeContenerError)
	ctx.Step(`^la respuesta debe contener error "([^"]*)"$`, laRespuestaDebeContenerErrorMsg)
	ctx.Step(`^la respuesta debe contener message "([^"]*)"$`, laRespuestaDebeContenerMessage)
	ctx.Step(`^la respuesta debe contener status "([^"]*)"$`, laRespuestaDebeContenerStatus)
	ctx.Step(`^la respuesta debe contener username "([^"]*)"$`, laRespuestaDebeContenerUsername)
	ctx.Step(`^la respuesta debe contener firstName "([^"]*)"$`, laRespuestaDebeContenerFirstName)
	
	ctx.Step(`^el header Content-Type debe ser "([^"]*)"$`, elHeaderContentTypeDebeSer)
	ctx.Step(`^el header Access-Control-Allow-Origin debe existir$`, elHeaderAccessControlAllowOriginDebeExistir)
	ctx.Step(`^el header Access-Control-Allow-Methods debe incluir GET$`, elHeaderAccessControlAllowMethodsDebeIncluirGET)

	// Hook para limpiar el contexto antes de cada escenario
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		apiContext = NewAPIGatewayContext()
	})
}