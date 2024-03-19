// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/configuration": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener las configuraciones",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.AllConfigurations"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para crear una configuracion",
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/configuration/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener una configuracion por id",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.Configuration"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para modificar una configuracion",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para eliminar una configuracion por id",
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/errores/por-dia": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener las errores detectados por dia",
                "parameters": [
                    {
                        "type": "string",
                        "format": "2006-01-02",
                        "description": "Fecha de comienzo del periodo - valor por defecto: 7 dias atras",
                        "name": "timeStart",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "2006-01-02",
                        "description": "Fecha del final del periodo - valor por defecto: hoy",
                        "name": "timeEnd",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.ErrorsCountPerDayAndType"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/healthcheck": {
            "get": {
                "description": "get the status of the server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Show the status of the server.",
                "responses": {
                    "200": {
                        "description": "Server is up and running",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/inputs/metricas-generales": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener las metricas generales de inputs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.InputsGeneralMetrics"
                        }
                    }
                }
            }
        },
        "/series/curadas/{serie_id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener los valores de una serie curada por id",
                "parameters": [
                    {
                        "type": "string",
                        "format": "2006-01-02",
                        "description": "Fecha de comienzo del periodo - valor por defecto: 7 dias atras",
                        "name": "timeStart",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "2006-01-02",
                        "description": "Fecha del final del periodo - valor por defecto: 5 dias despues",
                        "name": "timeEnd",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.StreamsDataResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/series/estaciones": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener el resumen de las series agrupado por estacion",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.StreamsPerStationResponse"
                        }
                    }
                }
            }
        },
        "/series/observadas/{serie_id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener los valores de una serie observada por id",
                "parameters": [
                    {
                        "type": "string",
                        "format": "2006-01-02",
                        "description": "Fecha de comienzo del periodo - valor por defecto: 7 dias atras",
                        "name": "timeStart",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "2006-01-02",
                        "description": "Fecha del final del periodo - valor por defecto: mañana",
                        "name": "timeEnd",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.StreamsDataResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/series/pronosticadas/{calibrado_id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener los valores de una serie pronosticadas por id",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.CalibratedStreamsDataResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/series/redes": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Endpoint para obtener el resumen de las series agrupado por red",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.StreamsPerNetworkResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.AllConfigurations": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "dtos.CalibratedStreamsData": {
            "type": "object",
            "properties": {
                "qualifier": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "dtos.CalibratedStreamsDataResponse": {
            "type": "object",
            "properties": {
                "mainStreams": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.CalibratedStreamsData"
                    }
                },
                "p05Streams": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.CalibratedStreamsData"
                    }
                },
                "p25Streams": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.CalibratedStreamsData"
                    }
                },
                "p75Streams": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.CalibratedStreamsData"
                    }
                },
                "p95Streams": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.CalibratedStreamsData"
                    }
                }
            }
        },
        "dtos.Configuration": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "dtos.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        },
        "dtos.ErrorsCountPerDayAndType": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                },
                "errorType": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "dtos.InputsGeneralMetrics": {
            "type": "object",
            "properties": {
                "totalNetworks": {
                    "type": "integer"
                },
                "totalStations": {
                    "type": "integer"
                },
                "totalStreams": {
                    "type": "integer"
                }
            }
        },
        "dtos.StreamsData": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "dtos.StreamsDataResponse": {
            "type": "object",
            "properties": {
                "streams": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.StreamsData"
                    }
                }
            }
        },
        "dtos.StreamsPerNetwork": {
            "type": "object",
            "properties": {
                "networkId": {
                    "type": "string"
                },
                "networkName": {
                    "type": "string"
                },
                "streamsCount": {
                    "type": "integer"
                }
            }
        },
        "dtos.StreamsPerNetworkResponse": {
            "type": "object",
            "properties": {
                "networks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.StreamsPerNetwork"
                    }
                }
            }
        },
        "dtos.StreamsPerStation": {
            "type": "object",
            "properties": {
                "stationId": {
                    "type": "string"
                },
                "stationName": {
                    "type": "string"
                },
                "streamsCount": {
                    "type": "integer"
                }
            }
        },
        "dtos.StreamsPerStationResponse": {
            "type": "object",
            "properties": {
                "stations": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.StreamsPerStation"
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
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Inputs API",
	Description:      "This API manages the inputs of the forecast model",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
