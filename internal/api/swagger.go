// internal/api/swagger.go
package api

import (
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

type SwaggerConfig struct {
	SpecPath string
	UITitle  string
}

func SetupSwagger(app *fiber.App, config SwaggerConfig) {
	if config.SpecPath == "" {
		config.SpecPath = "./contract.yaml"
	}

	if config.UITitle == "" {
		config.UITitle = "ShareTrip API Documentation"
	}

	// Загружаем спецификацию
	spec, err := loadOpenAPISpec(config.SpecPath)
	if err != nil {
		spec = getDefaultSpec()
	}

	// Эндпоинт для JSON спецификации (используется Swagger UI)
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.JSON(spec)
	})

	// Эндпоинт для YAML спецификации
	app.Get("/swagger/spec.yaml", func(c *fiber.Ctx) error {
		return c.SendFile(config.SpecPath)
	})

	// HTML страница Swagger UI
	app.Get("/swagger", func(c *fiber.Ctx) error {
		html := getSwaggerHTML(config.UITitle)
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	// Редирект с /swagger/index.html на /swagger
	app.Get("/swagger/index.html", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger")
	})
}

func loadOpenAPISpec(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(data, &spec); err != nil {
		var yamlSpec map[string]interface{}
		if err := yaml.Unmarshal(data, &yamlSpec); err != nil {
			return nil, err
		}
		return yamlSpec, nil
	}
	return spec, nil
}

func getDefaultSpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       "ShareTrip API",
			"version":     "1.0.0",
			"description": "API for ShareTrip service. Please load contract.yaml for full specification.",
		},
		"paths": map[string]interface{}{
			"/api/ready": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Health check endpoint",
					"description": "Returns service ready status",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Service is ready",
						},
					},
				},
			},
		},
		"servers": []map[string]interface{}{
			{
				"url": "http://localhost:8080",
			},
		},
	}
}

func getSwaggerHTML(title string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + title + `</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                validatorUrl: null,
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`
}