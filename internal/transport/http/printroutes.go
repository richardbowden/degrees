package thttp

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Helper function to extract handler names
func getHandlerName(handler http.Handler) string {
	if handler == nil {
		return "nil"
	}

	// Handle http.HandlerFunc
	if handlerFunc, ok := handler.(http.HandlerFunc); ok {
		return getFunctionName(handlerFunc)
	}

	// Handle other handler types
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() == reflect.Ptr {
		return handlerType.Elem().Name()
	}
	return handlerType.Name()
}

// Helper function to extract middleware names
func getMiddlewareName(middleware func(http.Handler) http.Handler) string {
	return getFunctionName(middleware)
}

// Helper function to get function names using reflection
func getFunctionName(fn interface{}) string {
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return "unknown"
	}

	fnPtr := fnValue.Pointer()
	fnInfo := runtime.FuncForPC(fnPtr)
	if fnInfo == nil {
		return "unknown"
	}

	name := fnInfo.Name()

	// Clean up the name
	if lastSlash := strings.LastIndex(name, "/"); lastSlash != -1 {
		name = name[lastSlash+1:]
	}
	if lastDot := strings.LastIndex(name, "."); lastDot != -1 {
		name = name[lastDot+1:]
	}

	// Remove common suffixes
	name = strings.TrimSuffix(name, "-fm")
	name = strings.TrimSuffix(name, ".func1")

	return name
}

// Helper function to recursively print route structure
func PrintRoutes(routes []chi.Route, indent string) {
	for _, route := range routes {
		// Print each HTTP method for this route pattern
		if len(route.Handlers) > 0 {
			for method := range route.Handlers {
				fmt.Printf("%s%s %s\n", indent, method, route.Pattern)
			}
		} else {
			// Route with no handlers (probably just a mount point)
			fmt.Printf("%s* %s (mount point)\n", indent, route.Pattern)
		}

		// Print sub-routes recursively
		if route.SubRoutes != nil && len(route.SubRoutes.Routes()) > 0 {
			PrintRoutes(route.SubRoutes.Routes(), indent+"-")
		}
	}
}

// func PrintRoutes(api huma.API) {
// 	fmt.Println("\n=== Registered Huma Routes ===")
//
// 	// Get the OpenAPI spec
// 	spec := api.OpenAPI()
//
// 	// Create a slice to store routes for sorting
// 	var routes []string
//
// 	// Iterate through all paths and methods
// 	for path, pathItem := range spec.Paths {
// 		if pathItem.Get != nil {
// 			routes = append(routes, fmt.Sprintf("GET    %s", path))
// 		}
// 		if pathItem.Post != nil {
// 			routes = append(routes, fmt.Sprintf("POST   %s", path))
// 		}
// 		if pathItem.Put != nil {
// 			routes = append(routes, fmt.Sprintf("PUT    %s", path))
// 		}
// 		if pathItem.Patch != nil {
// 			routes = append(routes, fmt.Sprintf("PATCH  %s", path))
// 		}
// 		if pathItem.Delete != nil {
// 			routes = append(routes, fmt.Sprintf("DELETE %s", path))
// 		}
// 		if pathItem.Head != nil {
// 			routes = append(routes, fmt.Sprintf("HEAD   %s", path))
// 		}
// 		if pathItem.Options != nil {
// 			routes = append(routes, fmt.Sprintf("OPTIONS %s", path))
// 		}
// 	}
//
// 	// Sort routes alphabetically
// 	sort.Strings(routes)
//
// 	// Print all routes
// 	for _, route := range routes {
// 		fmt.Println(route)
// 	}
//
// 	fmt.Printf("\nTotal endpoints: %d\n", len(routes))
// 	fmt.Println("==============================")
// }
