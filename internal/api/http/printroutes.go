package apihttp

import (
	"fmt"
	"sort"

	"github.com/danielgtaylor/huma/v2"
)

func PrintRoutes(api huma.API) {
	fmt.Println("\n=== Registered Huma Routes ===")

	// Get the OpenAPI spec
	spec := api.OpenAPI()

	// Create a slice to store routes for sorting
	var routes []string

	// Iterate through all paths and methods
	for path, pathItem := range spec.Paths {
		if pathItem.Get != nil {
			routes = append(routes, fmt.Sprintf("GET    %s", path))
		}
		if pathItem.Post != nil {
			routes = append(routes, fmt.Sprintf("POST   %s", path))
		}
		if pathItem.Put != nil {
			routes = append(routes, fmt.Sprintf("PUT    %s", path))
		}
		if pathItem.Patch != nil {
			routes = append(routes, fmt.Sprintf("PATCH  %s", path))
		}
		if pathItem.Delete != nil {
			routes = append(routes, fmt.Sprintf("DELETE %s", path))
		}
		if pathItem.Head != nil {
			routes = append(routes, fmt.Sprintf("HEAD   %s", path))
		}
		if pathItem.Options != nil {
			routes = append(routes, fmt.Sprintf("OPTIONS %s", path))
		}
	}

	// Sort routes alphabetically
	sort.Strings(routes)

	// Print all routes
	for _, route := range routes {
		fmt.Println(route)
	}

	fmt.Printf("\nTotal endpoints: %d\n", len(routes))
	fmt.Println("==============================")
}
