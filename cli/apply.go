package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

type ApplyResult struct {
	Kind   string
	Name   string
	Result string
}

func (r *Operations) ApplyCmd() *cobra.Command {
	var filePath string
	var recursive bool

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply a configuration to a resource by file",
		Long:  "Apply a configuration to a resource by file",
		Example: `
			bl apply -f ./my-deployment.yaml
			# Or using stdin
			cat file.yaml | bl apply -f -
		`,
		Run: func(cmd *cobra.Command, args []string) {
			_, err := r.Apply(filePath, recursive)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&filePath, "filename", "f", "", "Path to YAML file to apply")
	cmd.Flags().BoolVarP(&recursive, "recursive", "R", false, "Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.")
	err := cmd.MarkFlagRequired("filename")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return cmd
}

func (r *Operations) Apply(filePath string, recursive bool) ([]ApplyResult, error) {
	results, err := getResults(filePath, recursive)
	if err != nil {
		return nil, fmt.Errorf("error getting results: %w", err)
	}
	applyResults := []ApplyResult{}

	// À ce stade, results contient tous vos documents YAML
	for _, result := range results {
		for _, resource := range resources {
			if resource.Kind == result.Kind {
				name := result.Metadata.(map[string]interface{})["name"].(string)
				status := resource.PutFn(resource.Kind, name, result)
				applyResults = append(applyResults, ApplyResult{
					Kind:   resource.Kind,
					Name:   name,
					Result: status,
				})
			}
		}
	}
	return applyResults, nil
}

// Helper function to handle common resource operations
func (resource Resource) handleResourceOperation(name string, resourceObject interface{}, operation string) (*http.Response, error) {
	ctx := context.Background()

	// Get the appropriate function based on operation
	var fn reflect.Value
	if operation == "put" {
		fn = reflect.ValueOf(resource.Put)
	} else {
		fn = reflect.ValueOf(resource.Post)
	}

	if fn.Kind() != reflect.Func {
		return nil, fmt.Errorf("fn is not a valid function")
	}

	// Handle spec conversion
	specJson, err := json.Marshal(resourceObject)
	if err != nil {
		return nil, fmt.Errorf("error marshaling spec: %v", err)
	}

	destBody := reflect.New(resource.SpecType).Interface()
	if err := json.Unmarshal(specJson, destBody); err != nil {
		return nil, fmt.Errorf("error unmarshaling to target type: %v", err)
	}

	// Call the function
	var results []reflect.Value
	switch operation {
	case "put":
		results = fn.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(name),
			reflect.ValueOf(destBody).Elem(),
		})
	case "post":
		results = fn.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(destBody).Elem(),
		})
	default:
		return nil, fmt.Errorf("invalid operation: %s", operation)
	}

	if len(results) <= 1 {
		return nil, nil
	}

	if err, ok := results[1].Interface().(error); ok && err != nil {
		return nil, err
	}

	response, ok := results[0].Interface().(*http.Response)
	if !ok {
		return nil, fmt.Errorf("the result is not a pointer to http.Response")
	}

	return response, nil
}

func (resource Resource) PutFn(resourceName string, name string, resourceObject interface{}) string {
	formattedError := fmt.Sprintf("Resource %s:%s error: ", resourceName, name)
	response, err := resource.handleResourceOperation(name, resourceObject, "put")
	if err != nil {
		fmt.Printf("%s%v", formattedError, err)
		return "failed"
	}
	if response == nil {
		return "failed"
	}

	defer response.Body.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, response.Body); err != nil {
		fmt.Printf("%s%v", formattedError, err)
		return "failed"
	}

	if response.StatusCode == 404 {
		// Need to create the resource
		return resource.PostFn(resourceName, name, resourceObject)
	}

	if response.StatusCode >= 400 {
		ErrorHandler(response.Request, resourceName, name, buf.String())
		return "failed"
	}

	var res interface{}
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		fmt.Printf("%s%v", formattedError, err)
		return "failed"
	}
	fmt.Printf("Resource %s:%s configured\n", resourceName, name)
	return "configured"
}

func (resource Resource) PostFn(resourceName string, name string, resourceObject interface{}) string {
	formattedError := fmt.Sprintf("Resource %s:%s error: ", resourceName, name)
	response, err := resource.handleResourceOperation(name, resourceObject, "post")
	if err != nil {
		fmt.Printf("%s%v\n", formattedError, err)
		return "failed"
	}
	if response == nil {
		return "failed"
	}

	defer response.Body.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, response.Body); err != nil {
		fmt.Printf("%s%v\n", formattedError, err)
		return "failed"
	}

	if response.StatusCode >= 400 {
		ErrorHandler(response.Request, resourceName, name, buf.String())
		return "failed"
	}

	var res interface{}
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		fmt.Printf("%s%v\n", formattedError, err)
		return "failed"
	}
	fmt.Printf("Resource %s:%s created\n", resourceName, name)
	return "created"
}
