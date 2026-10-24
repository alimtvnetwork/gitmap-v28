package cmddb

import (
	"encoding/json"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"net/http"
	"strings"
)

// OpenAPISpec represents the root OpenAPI 3.0.3 specification structure.
type OpenAPISpec struct {
	OpenAPI    string                     `json:"openapi"`
	Info       OpenAPIInfo                `json:"info"`
	Servers    []OpenAPIServer            `json:"servers"`
	Tags       []OpenAPITag               `json:"tags"`
	Paths      map[string]OpenAPIPathItem `json:"paths"`
	Components OpenAPIComponents          `json:"components"`
}

type OpenAPIInfo struct {
	Title, Version, Description string
}

func (i OpenAPIInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"title": i.Title, "version": i.Version, "description": i.Description})
}

type OpenAPIServer struct {
	URL, Description string
}

func (s OpenAPIServer) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"url": s.URL, "description": s.Description})
}

type OpenAPITag struct {
	Name, Description string
}

func (t OpenAPITag) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"name": t.Name, "description": t.Description})
}

type OpenAPIPathItem struct {
	Get, Post, Put, Delete *OpenAPIOperation
}

func (p OpenAPIPathItem) MarshalJSON() ([]byte, error) {
	m := make(map[string]*OpenAPIOperation)
	if p.Get != nil {
		m["get"] = p.Get
	}
	if p.Post != nil {
		m["post"] = p.Post
	}
	if p.Put != nil {
		m["put"] = p.Put
	}
	if p.Delete != nil {
		m["delete"] = p.Delete
	}
	return json.Marshal(m)
}

type OpenAPIOperation struct {
	Tags        []string                   `json:"tags"`
	Summary     string                     `json:"summary"`
	Parameters  []OpenAPIParameter         `json:"parameters,omitempty"`
	RequestBody *OpenAPIRequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]OpenAPIResponse `json:"responses"`
}

type OpenAPIParameter struct {
	Name, In string
	Required bool
	Schema   *OpenAPISchema
}

func (p OpenAPIParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{"name": p.Name, "in": p.In, "required": p.Required, "schema": p.Schema})
}

type OpenAPIRequestBody struct {
	Required bool                          `json:"required"`
	Content  map[string]OpenAPIMediaObject `json:"content"`
}

type OpenAPIMediaObject struct {
	Schema *OpenAPISchema `json:"schema,omitempty"`
}

type OpenAPIResponse struct {
	Description string                        `json:"description"`
	Content     map[string]OpenAPIMediaObject `json:"content,omitempty"`
}

type OpenAPIComponents struct {
	Schemas map[string]OpenAPISchema `json:"schemas"`
}

type OpenAPISchema struct {
	Type                 string                   `json:"type,omitempty"`
	Format               string                   `json:"format,omitempty"`
	Ref                  string                   `json:"$ref,omitempty"`
	Required             []string                 `json:"required,omitempty"`
	Properties           map[string]OpenAPISchema `json:"properties,omitempty"`
	Items                *OpenAPISchema           `json:"items,omitempty"`
	AdditionalProperties any                      `json:"additionalProperties,omitempty"`
}

// GenerateOpenAPISpec builds the complete OpenAPI 3.0.3 specification.
func GenerateOpenAPISpec(mgr *MultiDbManager) (*OpenAPISpec, error) {
	spec := newBaseOpenAPISpec()
	populateComponents(&spec.Components)
	populateSystemPaths(spec.Paths)
	if mgr != nil {
		populateDatabasePaths(mgr, spec.Paths, spec.Components.Schemas)
	}
	return spec, nil
}

func newBaseOpenAPISpec() *OpenAPISpec {
	return &OpenAPISpec{
		OpenAPI: "3.0.3",
		Info:    OpenAPIInfo{Title: "GitMap Database REST API", Version: "1.0.0", Description: "Dynamic REST API and schema explorer."},
		Servers: []OpenAPIServer{{URL: "/", Description: "Local GitMap REST Server"}},
		Tags: []OpenAPITag{
			{Name: "System", Description: "Database catalog and server health"},
			{Name: "Schema", Description: "Table and column discovery endpoints"},
			{Name: "Records", Description: "Generic CRUD operations across tables"},
			{Name: "SQL Engine", Description: "Direct SQL queries and batch execution"},
			{Name: "Views", Description: "Dynamic view creation and query endpoints"},
			{Name: "DDL & Migrations", Description: "Table creation, column alterations, and migrations"},
		},
		Paths:      make(map[string]OpenAPIPathItem),
		Components: OpenAPIComponents{Schemas: make(map[string]OpenAPISchema)},
	}
}

func populateComponents(comp *OpenAPIComponents) {
	comp.Schemas["RestStatus"] = OpenAPISchema{
		Type:     "object",
		Required: []string{"IsSuccess", "IsFailed", "Code", "Message", "Timestamp"},
		Properties: map[string]OpenAPISchema{
			"IsSuccess": {Type: "boolean"}, "IsFailed": {Type: "boolean"},
			"Code": {Type: "integer"}, "Message": {Type: "string"},
			"Timestamp": {Type: "string", Format: "date-time"},
		},
	}
	comp.Schemas["UniversalEnvelope"] = OpenAPISchema{
		Type:     "object",
		Required: []string{"Status", "Attributes", "Results"},
		Properties: map[string]OpenAPISchema{
			"Status":     {Ref: "#/components/schemas/RestStatus"},
			"Attributes": {Type: "object", AdditionalProperties: true},
			"Results":    {Type: "array", Items: &OpenAPISchema{}},
			"Errors":     {Type: "array", Items: &OpenAPISchema{Type: "string"}},
		},
	}
}

func populateSystemPaths(paths map[string]OpenAPIPathItem) {
	paths["/api/v1/databases"] = OpenAPIPathItem{
		Get: &OpenAPIOperation{
			Tags: []string{"System"}, Summary: "List mounted databases",
			Responses: map[string]OpenAPIResponse{"200": {Description: "Catalog of databases"}},
		},
	}
	paths["/api/v1/db/{alias}/tables"] = OpenAPIPathItem{Get: newDbOp("Schema", "List tables", pathParam("alias"))}
	paths["/api/v1/db/{alias}/schema"] = OpenAPIPathItem{Get: newDbOp("Schema", "Inspect schema", pathParam("alias"))}
	paths["/api/v1/db/{alias}/find-columns"] = OpenAPIPathItem{Get: newDbOp("Schema", "Search columns", pathParam("alias"), queryParam("name"))}
	paths["/api/v1/db/{alias}/query"] = OpenAPIPathItem{
		Get:  newDbOp("SQL Engine", "Read-only query", pathParam("alias"), queryParam("sql")),
		Post: newDbOp("SQL Engine", "Parameterized query", pathParam("alias")),
	}
	populateGenericDdlPaths(paths)
}

func populateGenericDdlPaths(paths map[string]OpenAPIPathItem) {
	paths["/api/v1/db/{alias}/views"] = OpenAPIPathItem{
		Get:  newDbOp("Views", "List dynamic views", pathParam("alias")),
		Post: newDbOp("Views", "Create dynamic view", pathParam("alias")),
	}
	paths["/api/v1/db/{alias}/views/{name}"] = OpenAPIPathItem{Delete: newDbOp("Views", "Drop dynamic view", pathParam("alias"), pathParam("name"))}
	paths["/api/v1/db/{alias}/alter"] = OpenAPIPathItem{Post: newDbOp("DDL & Migrations", "Alter table", pathParam("alias"))}
	paths["/api/v1/db/{alias}/migrate"] = OpenAPIPathItem{Post: newDbOp("DDL & Migrations", "Run migration", pathParam("alias"))}
	paths["/api/v1/db/{alias}/clone"] = OpenAPIPathItem{Post: newDbOp("DDL & Migrations", "Clone database", pathParam("alias"))}
	paths["/api/v1/db/{alias}/{table}"] = OpenAPIPathItem{
		Get:  newDbOp("Records", "Query records", pathParam("alias"), pathParam("table")),
		Post: newDbOp("Records", "Insert record", pathParam("alias"), pathParam("table")),
	}
	paths["/api/v1/db/{alias}/{table}/{id}"] = OpenAPIPathItem{
		Get:    newDbOp("Records", "Fetch record by PK", pathParam("alias"), pathParam("table"), pathParam("id")),
		Put:    newDbOp("Records", "Update record by PK", pathParam("alias"), pathParam("table"), pathParam("id")),
		Delete: newDbOp("Records", "Delete record by PK", pathParam("alias"), pathParam("table"), pathParam("id")),
	}
}

func populateDatabasePaths(mgr *MultiDbManager, paths map[string]OpenAPIPathItem, schemas map[string]OpenAPISchema) {
	for _, alias := range mgr.ListAliases() {
		tables, err := mgr.GetTables(alias)
		if err != nil {
			continue
		}
		for _, tab := range tables {
			schemaKey := alias + "_" + tab.Name
			cols, _ := mgr.GetColumns(alias, tab.Name)
			schemas[schemaKey] = buildTableSchema(cols)
		}
	}
}

func buildTableSchema(cols []RestColumnInfo) OpenAPISchema {
	props := make(map[string]OpenAPISchema)
	var req []string
	for _, c := range cols {
		schemaType, schemaFmt := MapSqliteTypeToOpenAPI(c.Type)
		props[c.Name] = OpenAPISchema{Type: schemaType, Format: schemaFmt}
		if !c.IsNullable {
			req = append(req, c.Name)
		}
	}
	return OpenAPISchema{Type: "object", Required: req, Properties: props}
}

// MapSqliteTypeToOpenAPI normalizes SQLite affinities to OpenAPI 3.0 types.
func MapSqliteTypeToOpenAPI(sqliteType string) (string, string) {
	upper := strings.ToUpper(strings.TrimSpace(sqliteType))
	if strings.Contains(upper, "INT") {
		return "integer", "int64"
	}
	if strings.Contains(upper, "REAL") || strings.Contains(upper, "FLOA") || strings.Contains(upper, "DOUBLE") || strings.Contains(upper, "NUM") {
		return "number", "double"
	}
	if strings.Contains(upper, "BOOL") {
		return "boolean", ""
	}
	if strings.Contains(upper, "DATE") || strings.Contains(upper, "TIME") {
		return "string", "date-time"
	}
	if strings.Contains(upper, "BLOB") {
		return "string", "byte"
	}
	return "string", ""
}

func newDbOp(tag, summary string, params ...OpenAPIParameter) *OpenAPIOperation {
	return &OpenAPIOperation{
		Tags: []string{tag}, Summary: summary, Parameters: params,
		Responses: map[string]OpenAPIResponse{"200": {Description: "Success envelope"}},
	}
}

func pathParam(name string) OpenAPIParameter {
	return OpenAPIParameter{Name: name, In: "path", Required: true, Schema: &OpenAPISchema{Type: "string"}}
}

func queryParam(name string) OpenAPIParameter {
	return OpenAPIParameter{Name: name, In: "query", Required: false, Schema: &OpenAPISchema{Type: "string"}}
}

// HandleSwaggerDoc serves the OpenAPI 3.0.3 specification JSON.
func HandleSwaggerDoc(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spec, err := GenerateOpenAPISpec(mgr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(spec)
	}
}

// HandleSwaggerUI serves the embedded Swagger UI single-page application.
func HandleSwaggerUI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(SwaggerUIHtml))
	}
}

// RegisterSwaggerRoutes registers Swagger JSON and UI endpoints onto the ServeMux.
func RegisterSwaggerRoutes(mux *http.ServeMux, mgr *MultiDbManager) {
	docHandler := HandleSwaggerDoc(mgr)
	uiHandler := HandleSwaggerUI()
	mux.HandleFunc("GET /swagger/doc.json", docHandler)
	mux.HandleFunc("GET /swagger.json", docHandler)
	mux.HandleFunc("GET /swagger", uiHandler)
	mux.HandleFunc("GET /swagger/", uiHandler)
}
