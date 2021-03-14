// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Name Test",
            "url": "http://contact.sample/text",
            "email": "sample@mail.io"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/books": {
            "get": {
                "description": "Get the books in the repository",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Books"
                ],
                "summary": "Get Books",
                "responses": {
                    "200": {
                        "description": "List of Books",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Book"
                            }
                        }
                    },
                    "500": {
                        "description": "err.repo_ops",
                        "schema": {
                            "$ref": "#/definitions/dto.ApiError"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new book from the passed data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Books"
                ],
                "summary": "Create a new book",
                "parameters": [
                    {
                        "description": "Book Data",
                        "name": "book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Book"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Book"
                        }
                    },
                    "422": {
                        "description": "err.duplicate_key || Invalid data\"\t// TODO learn to make validation of params and body, make the response 400 (bad request)",
                        "schema": {
                            "$ref": "#/definitions/dto.ApiError"
                        }
                    },
                    "500": {
                        "description": "err.repo_ops || Internal error, same struct as Iris.Problem",
                        "schema": {
                            "$ref": "#/definitions/dto.ApiError"
                        }
                    }
                }
            }
        },
        "/books/{id}": {
            "get": {
                "description": "Get a book through its Id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Books"
                ],
                "summary": "Get book by Id",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "uint32",
                        "description": "Requested Book Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Book"
                        }
                    },
                    "404": {
                        "description": "err.not_found",
                        "schema": {
                            "$ref": "#/definitions/dto.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal error, same struct as Iris.Problem",
                        "schema": {
                            "$ref": "#/definitions/dto.ApiError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a Book by its Id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Books"
                ],
                "summary": "Delete a Book",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "uint32",
                        "description": "Account ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "err.not_found",
                        "schema": {
                            "$ref": "#/definitions/dto.ApiError"
                        }
                    },
                    "500": {
                        "description": "err.repo_ops",
                        "schema": {
                            "$ref": "#/definitions/dto.ApiError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.ApiError": {
            "type": "object",
            "properties": {
                "detail": {
                    "type": "string",
                    "example": "Some error details"
                },
                "status": {
                    "type": "integer",
                    "example": 503
                },
                "title": {
                    "type": "string",
                    "example": "err_code"
                }
            }
        },
        "models.Book": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "example": "2021-03-12T02:11:03.292442-05:00"
                },
                "id": {
                    "type": "integer",
                    "example": 24
                },
                "items": {
                    "type": "integer",
                    "example": 46
                },
                "name": {
                    "type": "string",
                    "example": "The Book of Eli"
                },
                "updatedAt": {
                    "type": "string",
                    "example": "0001-01-01T00:00:00Z"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.0",
	Host:        "localhost:8080",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Shell Project",
	Description: "Api description shell project",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
