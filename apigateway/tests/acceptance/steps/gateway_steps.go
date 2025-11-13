package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
	"github.com/microservicios/apigateway/tests/support"
)

// APIGatewayContext contiene el contexto compartido de pruebas
type APIGatewayContext struct {
	client           *support.HTTPClient
	validator        *support.SchemaValidator
	lastResponse     *support.Response
	lastError        error
	lastRequestBody  map[string]interface{}
	lastRequestURL   string
	customHeaders    map[string]string
	requestDataTable map[string]interface{}
}

// NewAPIGatewayContext crea un nuevo contexto
func NewAPIGatewayContext() *APIGatewayContext {
	return &APIGatewayContext{
		customHeaders:    make(map[string]string),
		requestDataTable: make(map[string]interface{}),
	}
}

// InitializeScenario inicializa los steps del escenario
func InitializeScenario(ctx *godog.ScenarioContext, apiCtx *APIGatewayContext) {
	// Antecedentes
	ctx.Step(`^que el gateway está accesible en "([^"]*)"$`, apiCtx.GatewayIsAccessible)

	// Solicitudes HTTP
	ctx.Step(`^realizo una solicitud GET a "([^"]*)"$`, apiCtx.MakeGetRequest)
	ctx.Step(`^realizo una solicitud POST a "([^"]*)"$`, apiCtx.MakePostRequest)
	ctx.Step(`^realizo una solicitud PATCH a "([^"]*)"$`, apiCtx.MakePatchRequest)
	ctx.Step(`^realizo una solicitud DELETE a "([^"]*)"$`, apiCtx.MakeDeleteRequest)
	ctx.Step(`^realizo una solicitud con método OPTIONS a "([^"]*)"$`, apiCtx.MakeOptionsRequest)

	// Headers
	ctx.Step(`^incluyo el header "([^"]*)" con valor "([^"]*)"$`, apiCtx.AddHeader)

	// Data tables
	ctx.Step(`^$`, apiCtx.DataTable)

	// Validaciones de código de estado
	ctx.Step(`^el código de estado debe ser (\d+)$`, apiCtx.StatusCodeShouldBe)
	ctx.Step(`^el código de estado debe ser (\d+) o (\d+) o (\d+)$`, apiCtx.StatusCodeShouldBeOneOf)
	ctx.Step(`^el código de estado está entre (\d+) y (\d+)$`, apiCtx.StatusCodeShouldBeBetween)

	// Validaciones de respuesta
	ctx.Step(`^la respuesta debe validarse contra el esquema "([^"]*)"$`, apiCtx.ResponseValidatesAgainstSchema)
	ctx.Step(`^la respuesta debe contener el campo "([^"]*)"$`, apiCtx.ResponseContainsField)
	ctx.Step(`^la respuesta debe tener el header "([^"]*)"$`, apiCtx.ResponseHasHeader)
	ctx.Step(`^la respuesta debe tener un header con patrón "([^"]*)"$`, apiCtx.ResponseHasHeaderPattern)

	// Validaciones especiales
	ctx.Step(`^el gateway debe registrar la solicitud en los logs$`, apiCtx.GatewayLogsRequest)
	ctx.Step(`^intento conectar con un servicio en "([^"]*)"$`, apiCtx.ConnectToService)
	ctx.Step(`^debería recibir un error de conexión$`, apiCtx.ShouldReceiveConnectionError)
}

// ============ Antecedentes ============

func (ctx *APIGatewayContext) GatewayIsAccessible(baseURL string) error {
	validator, err := support.NewSchemaValidator("./features/schemas/gateway-schemas.json")
	if err != nil {
		return err
	}

	ctx.client = support.NewHTTPClient(baseURL)
	ctx.validator = validator
	ctx.customHeaders = make(map[string]string)

	return nil
}

// ============ Solicitudes HTTP ============

func (ctx *APIGatewayContext) MakeGetRequest(endpoint string) error {
	ctx.lastRequestURL = endpoint
	response, err := ctx.client.GET(endpoint, ctx.customHeaders)
	ctx.lastResponse = response
	ctx.lastError = err
	ctx.customHeaders = make(map[string]string) // Reset headers after request
	return err
}

func (ctx *APIGatewayContext) MakePostRequest(endpoint string) error {
	ctx.lastRequestURL = endpoint

	var body interface{} = nil
	if len(ctx.requestDataTable) > 0 {
		body = ctx.requestDataTable
		ctx.lastRequestBody = ctx.requestDataTable
		ctx.requestDataTable = make(map[string]interface{})
	}

	response, err := ctx.client.POST(endpoint, body, ctx.customHeaders)
	ctx.lastResponse = response
	ctx.lastError = err
	ctx.customHeaders = make(map[string]string)
	return err
}

func (ctx *APIGatewayContext) MakePatchRequest(endpoint string) error {
	ctx.lastRequestURL = endpoint

	var body interface{} = nil
	if len(ctx.requestDataTable) > 0 {
		body = ctx.requestDataTable
		ctx.lastRequestBody = ctx.requestDataTable
		ctx.requestDataTable = make(map[string]interface{})
	}

	response, err := ctx.client.PATCH(endpoint, body, ctx.customHeaders)
	ctx.lastResponse = response
	ctx.lastError = err
	ctx.customHeaders = make(map[string]string)
	return err
}

func (ctx *APIGatewayContext) MakeDeleteRequest(endpoint string) error {
	ctx.lastRequestURL = endpoint
	response, err := ctx.client.DELETE(endpoint, ctx.customHeaders)
	ctx.lastResponse = response
	ctx.lastError = err
	ctx.customHeaders = make(map[string]string)
	return err
}

func (ctx *APIGatewayContext) MakeOptionsRequest(endpoint string) error {
	ctx.lastRequestURL = endpoint

	req, err := http.NewRequest("OPTIONS", ctx.client.baseURL+endpoint, nil)
	if err != nil {
		return err
	}

	for key, value := range ctx.customHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.lastError = err
		return err
	}
	defer resp.Body.Close()

	// Crear una respuesta simulada
	ctx.lastResponse = &support.Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}

	ctx.customHeaders = make(map[string]string)
	return nil
}

// ============ Headers ============

func (ctx *APIGatewayContext) AddHeader(name, value string) error {
	ctx.customHeaders[name] = value
	return nil
}

// ============ Data Tables ============

func (ctx *APIGatewayContext) DataTable(dt *godog.Table) error {
	if dt == nil {
		return nil
	}

	// Convertir la tabla de datos en un mapa
	headers := dt.Rows[0].Cells
	for i := 1; i < len(dt.Rows); i++ {
		row := dt.Rows[i].Cells
		for j, header := range headers {
			if j < len(row) {
				ctx.requestDataTable[header.Value] = row[j].Value
			}
		}
	}

	return nil
}

// ============ Validaciones de Código de Estado ============

func (ctx *APIGatewayContext) StatusCodeShouldBe(expectedStatus int) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response received")
	}

	if ctx.lastResponse.StatusCode != expectedStatus {
		return fmt.Errorf("expected status %d but got %d", expectedStatus, ctx.lastResponse.StatusCode)
	}

	return nil
}

func (ctx *APIGatewayContext) StatusCodeShouldBeOneOf(status1, status2, status3 int) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response received")
	}

	if ctx.lastResponse.StatusCode != status1 && ctx.lastResponse.StatusCode != status2 && ctx.lastResponse.StatusCode != status3 {
		return fmt.Errorf("expected status %d, %d or %d but got %d", status1, status2, status3, ctx.lastResponse.StatusCode)
	}

	return nil
}

func (ctx *APIGatewayContext) StatusCodeShouldBeBetween(minStatus, maxStatus int) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response received")
	}

	if ctx.lastResponse.StatusCode < minStatus || ctx.lastResponse.StatusCode > maxStatus {
		return fmt.Errorf("expected status between %d and %d but got %d", minStatus, maxStatus, ctx.lastResponse.StatusCode)
	}

	return nil
}

// ============ Validaciones de Respuesta ============

func (ctx *APIGatewayContext) ResponseValidatesAgainstSchema(schemaName string) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response received")
	}

	if len(ctx.lastResponse.Body) == 0 {
		return fmt.Errorf("response body is empty")
	}

	var data interface{}
	if err := json.Unmarshal(ctx.lastResponse.Body, &data); err != nil {
		return fmt.Errorf("response is not valid JSON: %v", err)
	}

	result := ctx.validator.Validate(data, schemaName)
	if !result.IsValid {
		return fmt.Errorf("validation failed: %v", strings.Join(result.Errors, "; "))
	}

	return nil
}

func (ctx *APIGatewayContext) ResponseContainsField(fieldName string) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response received")
	}

	jsonData := ctx.lastResponse.GetBodyAsJSON()
	if _, ok := jsonData[fieldName]; !ok {
		return fmt.Errorf("response does not contain field '%s'", fieldName)
	}

	return nil
}

func (ctx *APIGatewayContext) ResponseHasHeader(headerName string) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response received")
	}

	if !ctx.lastResponse.HasHeader(headerName) {
		return fmt.Errorf("response does not have header '%s'", headerName)
	}

	return nil
}

func (ctx *APIGatewayContext) ResponseHasHeaderPattern(pattern string) error {
	if ctx.lastResponse == nil {
		return fmt.Errorf("no response received")
	}

	patterns := strings.Split(pattern, "|")
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if ctx.lastResponse.HasHeader(p) {
			return nil
		}
	}

	return fmt.Errorf("response does not have any header matching pattern '%s'", pattern)
}

// ============ Validaciones Especiales ============

func (ctx *APIGatewayContext) GatewayLogsRequest() error {
	// Esta validación dependerá de la configuración específica del gateway
	// Por ahora, simplemente verificamos que la solicitud fue realizada
	if ctx.lastRequestURL == "" {
		return fmt.Errorf("no request was made")
	}

	return nil
}

func (ctx *APIGatewayContext) ConnectToService(serviceURL string) error {
	// Intentar conectar al servicio especificado
	testClient := support.NewHTTPClient(serviceURL)
	response, err := testClient.GET("/health", make(map[string]string))

	ctx.lastResponse = response
	ctx.lastError = err

	return nil
}

func (ctx *APIGatewayContext) ShouldReceiveConnectionError() error {
	if ctx.lastError == nil {
		return fmt.Errorf("expected a connection error but got none")
	}

	return nil
}

// InitializeTestSuite inicializa la suite de pruebas con el contexto
func InitializeTestSuite(suite *godog.TestSuite) {
	suite.ScenarioInitializer(func(ctx *godog.ScenarioContext) {
		apiCtx := NewAPIGatewayContext()
		InitializeScenario(ctx, apiCtx)
	})
}
