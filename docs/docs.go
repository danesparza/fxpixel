// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/config": {
            "get": {
                "description": "Get the system configuration information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "Get the system configuration information",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/timelines": {
            "get": {
                "description": "List all timelines in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "List all timelines in the system",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Adds a timeline to the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "Adds a timeline to the system",
                "parameters": [
                    {
                        "description": "The timeline to add",
                        "name": "endpoint",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.Timeline"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/timelines/tag/{tag}": {
            "get": {
                "description": "Gets timelines that have a tag",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "Gets timelines that have a tag",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The tag to use when fetching timelines",
                        "name": "tag",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/timelines/{id}": {
            "get": {
                "description": "Gets a single timeline",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "Gets a single timeline",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The timeline id to get",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Updates tags for a timeline",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "Updates tags for a timeline",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The timeline id to update tags for",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "The tags to set for the timeline",
                        "name": "endpoint",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UpdateTagsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a single timeline",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timeline"
                ],
                "summary": "Delete a single timeline",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The timeline id to delete",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "api.SystemResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "api.Timeline": {
            "type": "object",
            "properties": {
                "created": {
                    "description": "Timeline create time",
                    "type": "string"
                },
                "enabled": {
                    "description": "Timeline enabled or not",
                    "type": "boolean"
                },
                "gpio": {
                    "description": "The GPIO device to play the timeline on.  Optional.  If not set, uses the default",
                    "type": "integer"
                },
                "id": {
                    "description": "Unique Timeline ID",
                    "type": "string"
                },
                "name": {
                    "description": "Timeline name",
                    "type": "string"
                },
                "steps": {
                    "description": "Steps for the timeline",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.TimelineStep"
                    }
                },
                "tags": {
                    "description": "List of Tags to associate with this timeline",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "api.TimelineStep": {
            "type": "object",
            "properties": {
                "effect": {
                    "description": "The Effect type (if Type=effect)",
                    "type": "string"
                },
                "id": {
                    "description": "The timeline step id",
                    "type": "string"
                },
                "leds": {
                    "description": "Leds to use for the scene (optional) If not set and is required for the type, defaults to entire strip",
                    "type": "string"
                },
                "meta-info": {
                    "description": "Additional information required for specific types"
                },
                "number": {
                    "description": "The step number (ordinal position in the timeline)",
                    "type": "integer"
                },
                "time": {
                    "description": "Time (in milliseconds).  Some things (like trigger) don't require time",
                    "type": "integer"
                },
                "type": {
                    "description": "Timeline frame type (effect/sleep/trigger/loop)",
                    "type": "string"
                }
            }
        },
        "api.UpdateTagsRequest": {
            "type": "object",
            "properties": {
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "fxPixel",
	Description:      "fxPixel LED lighting effects REST service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
