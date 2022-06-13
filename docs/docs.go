// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/oauth": {
            "get": {
                "description": "Get OAuth Login URL",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth"
                ],
                "summary": "Get OAuth Login URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "redirect URL",
                        "name": "redirect",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "303": {
                        "description": ""
                    }
                }
            }
        },
        "/oauth/callback": {
            "get": {
                "description": "OAuth Callback",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth"
                ],
                "summary": "OAuth Callback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "state",
                        "name": "state",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "scope",
                        "name": "scope",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "auth user",
                        "name": "authuser",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "prompt",
                        "name": "prompt",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "303": {
                        "description": ""
                    }
                }
            }
        },
        "/oauth/token": {
            "post": {
                "security": [
                    {
                        "BearerToken": []
                    }
                ],
                "description": "OAuth Token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth"
                ],
                "summary": "OAuth Token",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TokenRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.TokenSchema"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.TokenRequest": {
            "type": "object",
            "required": [
                "code",
                "grant_type"
            ],
            "properties": {
                "code": {
                    "type": "string"
                },
                "grant_type": {
                    "type": "string",
                    "enum": [
                        "code",
                        "refresh_token"
                    ]
                }
            }
        },
        "model.TokenSchema": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expiry": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "token_type": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:7001",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Baal API",
	Description:      "Baal API Doc",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
