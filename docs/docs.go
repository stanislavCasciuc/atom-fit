// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "LoginHandler",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "LoginHandler",
                "parameters": [
                    {
                        "description": "Login Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.TokenResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new",
                "parameters": [
                    {
                        "description": "Register User Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.registerUserPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.TokenResponse"
                        }
                    }
                }
            }
        },
        "/exercise/{exerciseID}/like": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Unlike exercise",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "likes"
                ],
                "summary": "Unlike exercise",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Exercise ID",
                        "name": "exerciseID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/exercises": {
            "get": {
                "description": "Get all Exercises",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exercises"
                ],
                "summary": "Get all Exercises",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Since",
                        "name": "since",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Until",
                        "name": "until",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Sort",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Tags",
                        "name": "tags",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/store.Exercise"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new Exercise",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exercises"
                ],
                "summary": "Create a new Exercise",
                "parameters": [
                    {
                        "description": "Exercise Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/store.Exercise"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/store.Exercise"
                        }
                    }
                }
            }
        },
        "/exercises/{exerciseID}/like": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Like exercise",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "likes"
                ],
                "summary": "Like exercise",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Exercise ID",
                        "name": "exerciseID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/exercises/{id}": {
            "get": {
                "description": "Get Exercise by id from param",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exercises"
                ],
                "summary": "Get Exercise by id from param",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Exercise ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.Exercise"
                        }
                    }
                }
            }
        },
        "/nutrients/daily-goal": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get macronutrients goal per day for the user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "nutrients"
                ],
                "summary": "Get macronutrients goal per day",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/nutrients.UserNutrientsGoal"
                        }
                    }
                }
            }
        },
        "/reviews/workout/{workoutID}": {
            "get": {
                "description": "Get workout reviews",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reviews"
                ],
                "summary": "Get workout reviews",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Workout ID",
                        "name": "workoutID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/store.WorkoutReviewWithMetadata"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Review workout",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "reviews"
                ],
                "summary": "Review workout",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Workout ID",
                        "name": "workoutID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Review workout payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.WorkoutReviewPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.WorkoutReview"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.User"
                        }
                    }
                }
            }
        },
        "/users/activate": {
            "put": {
                "description": "Activate a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Activate a ActivateUser",
                "parameters": [
                    {
                        "description": "Activate User Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.ActivationPayload"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/users/attributes": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a user with attributes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user with attributes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.User"
                        }
                    }
                }
            }
        },
        "/users/attributes/log/weight": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Log a user weight",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Log a LogWeight",
                "parameters": [
                    {
                        "description": "Log Weight Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LogWeightPayload"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/users/attributes/weight": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a user weight",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user weight",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/store.UserWeightByDate"
                            }
                        }
                    }
                }
            }
        },
        "/workouts": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new workout",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workouts"
                ],
                "summary": "Create a new workout",
                "parameters": [
                    {
                        "description": "Create Workout Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.CreateWorkoutPayload"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/workouts/": {
            "get": {
                "description": "Get all workouts",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workouts"
                ],
                "summary": "Get all workouts",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Sort",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Tags",
                        "name": "tags",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/store.Workout"
                            }
                        }
                    }
                }
            }
        },
        "/workouts/{workoutID}": {
            "get": {
                "description": "Get workout by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workouts"
                ],
                "summary": "Get workout by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Workout ID",
                        "name": "workoutID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.Workout"
                        }
                    }
                }
            }
        },
        "/workouts/{workoutID}/like": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Like exercise",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "likes"
                ],
                "summary": "Like exercise",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Workout ID",
                        "name": "workoutID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Unlike Workout",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "likes"
                ],
                "summary": "Unlike Workout",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Workout ID",
                        "name": "workoutID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.ActivationPayload": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "handlers.CreateWorkoutExercisePayload": {
            "type": "object",
            "required": [
                "duration",
                "exercise_id"
            ],
            "properties": {
                "duration": {
                    "type": "integer"
                },
                "exercise_id": {
                    "type": "integer"
                }
            }
        },
        "handlers.CreateWorkoutPayload": {
            "type": "object",
            "required": [
                "exercises",
                "name"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "exercises": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handlers.CreateWorkoutExercisePayload"
                    }
                },
                "name": {
                    "type": "string"
                },
                "tutorial_link": {
                    "type": "string"
                }
            }
        },
        "handlers.LogWeightPayload": {
            "type": "object",
            "required": [
                "weight"
            ],
            "properties": {
                "weight": {
                    "type": "number"
                }
            }
        },
        "handlers.LoginPayload": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "handlers.TokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "handlers.WorkoutReviewPayload": {
            "type": "object",
            "required": [
                "content",
                "rating",
                "title"
            ],
            "properties": {
                "content": {
                    "type": "string"
                },
                "rating": {
                    "type": "integer",
                    "enum": [
                        1,
                        2,
                        3,
                        4,
                        5
                    ]
                },
                "title": {
                    "type": "string",
                    "maxLength": 255
                }
            }
        },
        "handlers.registerUserPayload": {
            "type": "object",
            "required": [
                "email",
                "goal",
                "is_male",
                "password",
                "username"
            ],
            "properties": {
                "age": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "goal": {
                    "type": "string",
                    "enum": [
                        "lose",
                        "gain",
                        "maintain"
                    ]
                },
                "height": {
                    "type": "integer",
                    "maximum": 250,
                    "minimum": 100
                },
                "is_male": {
                    "type": "boolean"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "username": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 4
                },
                "weight": {
                    "type": "number"
                },
                "weight_goal": {
                    "type": "number"
                }
            }
        },
        "nutrients.UserNutrientsGoal": {
            "type": "object",
            "properties": {
                "calories": {
                    "type": "number"
                },
                "carbohydrats": {
                    "type": "number"
                },
                "fats": {
                    "type": "number"
                },
                "proteins": {
                    "type": "number"
                }
            }
        },
        "response.SuccessResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        },
        "store.Exercise": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "is_duration": {
                    "type": "boolean"
                },
                "like": {
                    "type": "integer"
                },
                "muscles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "tutorial_link": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "store.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_active": {
                    "type": "boolean"
                },
                "user_attr": {
                    "$ref": "#/definitions/store.UserAttributes"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "store.UserAttributes": {
            "type": "object",
            "required": [
                "goal"
            ],
            "properties": {
                "age": {
                    "type": "integer"
                },
                "goal": {
                    "type": "string",
                    "enum": [
                        "lose",
                        "gain",
                        "maintain"
                    ]
                },
                "height": {
                    "type": "integer"
                },
                "is_male": {
                    "type": "boolean"
                },
                "user_id": {
                    "type": "integer"
                },
                "weight": {
                    "type": "number"
                },
                "weight_goal": {
                    "type": "number"
                }
            }
        },
        "store.UserWeightByDate": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                },
                "weight": {
                    "type": "number"
                }
            }
        },
        "store.Workout": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "likes": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "rating": {
                    "type": "number"
                },
                "reviews_count": {
                    "type": "integer"
                },
                "tutorial_link": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                },
                "workout_exercises": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/store.WorkoutExercises"
                    }
                }
            }
        },
        "store.WorkoutExercises": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "integer"
                },
                "exercise": {
                    "$ref": "#/definitions/store.Exercise"
                },
                "exercise_id": {
                    "type": "integer"
                },
                "workout_id": {
                    "type": "integer"
                }
            }
        },
        "store.WorkoutReview": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "rating": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                },
                "workout_id": {
                    "type": "integer"
                }
            }
        },
        "store.WorkoutReviewWithMetadata": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "rating": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                },
                "workout_id": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Atom Fit API",
	Description:      "This is a sample server for Atom Fit API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
