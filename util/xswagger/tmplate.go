package xswagger

import (
	"bytes"
	"html/template"
	"strings"
)

const sTemplate = `
{{- /* delete empty line */ -}}
swagger: '2.0'
info:
  description: 'swagger测试中'
  version: 1.0.0
  title: '{{ .Name }}'
host: {{ .Host }}
schemes:
  - {{ .Scheme }}
paths:

{{ range .Rpcs }}
  {{ .Route }}:
{{- if eq .Stype "post" }}
    post:
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
         - in: body
           name: body
           required: true
           schema:
             $ref: '#/definitions/{{ .Req }}'
      responses:
        '200':
          description:	正常返回
          schema:
            $ref: '#/definitions/{{ .Rep }}'

{{- else }}
    get:
      summary: ！！！参数结构体若是嵌套则只能使用第一级！！！
      produces:
        - application/json
      parameters:
	{{ range .MReq.Properties }}
		{{- if ne .Type "object" }}
        - name: {{ .Name }}
          in: query
          required: false
          type: {{ .Type }}
          format: {{ .Format }}
		{{- end }}
	{{- end }}
      responses:
        '200':
          description: 正常返回
          schema:
            $ref: '#/definitions/{{ .Rep }}'

{{- end }}
{{- end }}

definitions:
{{ range .Messages }}
  {{ .Name }}:
    type: object
    properties:
	{{ range .Properties }}
      {{ .Name }}:
		{{- if eq .Type "array" }}
        type: array
        items:
			{{- if eq .SubType "object" }}
           $ref: '#/definitions/{{ .SubFormat }}'
			{{- else }}
          type: {{ .SubType }}
			{{- end }}
		{{- else if ne  .Type "object" }}
        type: {{ .Type }}
        format: {{ .Format }}
		{{- else }}
        $ref: '#/definitions/{{ .Format }}'
		{{- end }}
	{{- end }}
{{- end }}


externalDocs:
  description: Find out more about Swagger
  url: http://swagger.io

`

var types = map[string]string{"int32": "integer", "int64": "integer", "uint32": "integer", "uint64": "integer", "sint32": "integer",
	"sint64": "integer", "fixed32": "integer", "fixed64": "integer", "sfixed32": "integer", "sfixed64": "integer", "bool": "boolean",
	"string": "string", "float": "float", "double": "double", "array": "array"}
var formats = map[string]string{"int32": "int32", "int64": "int64", "uint32": "int32", "uint64": "int64", "sint32": "int32",
	"sint64": "int64", "fixed32": "int32", "fixed64": "int64", "sfixed32": "int32", "sfixed64": "int64", "bool": "boolean",
	"string": "string", "float": "number", "double": "number", "array": "array"}

type MethodType uint8

// Service is a proto service.

func (s *service) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("swagger").Parse(sTemplate)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getType(src string) string {
	if t, ok := types[src]; ok {
		return t
	}
	return "object"
}
func getFormat(src string) string {
	if t, ok := formats[src]; ok {
		return t
	}
	return strings.Replace(src, "*", "", -1)
}
