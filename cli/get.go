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

func (r *Operations) GetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a resource",
	}
	for _, resource := range resources {
		subcmd := &cobra.Command{
			Use:     resource.Plural,
			Aliases: []string{resource.Singular, resource.Short},
			Short:   fmt.Sprintf("Get a %s", resource.Kind),
			Run: func(cmd *cobra.Command, args []string) {
				options := map[string]string{
					"environment": environment,
				}
				if len(args) == 0 {
					resource.ListFn(options)
					return
				}
				if len(args) == 1 {
					resource.GetFn(args[0], options)
				}
			},
		}
		cmd.AddCommand(subcmd)
	}

	return cmd
}

func (resource Resource) GetFn(name string, options map[string]string) {
	ctx := context.Background()
	formattedError := fmt.Sprintf("Resource %s:%s error: ", resource.Kind, name)
	// Use reflect to call the function
	funcValue := reflect.ValueOf(resource.Get)
	if funcValue.Kind() != reflect.Func {
		fmt.Printf("%s%s", formattedError, "fn is not a valid function")
		os.Exit(1)
	}
	// Create a slice for the arguments
	fnargs := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(name)} // Add the context and the resource name

	if resource.GetParamsType != nil {
		paramsValue := retrieveListParams(resource.GetParamsType, options)
		fnargs = append(fnargs, paramsValue)
	}

	// Call the function with the arguments
	results := funcValue.Call(fnargs)

	// Handle the results based on your needs
	if len(results) <= 1 {
		return
	}

	if err, ok := results[1].Interface().(error); ok && err != nil {
		fmt.Printf("%s%v", formattedError, err)
		os.Exit(1)
	}

	// Check if the first result is a pointer to http.Response
	response, ok := results[0].Interface().(*http.Response)
	if !ok {
		fmt.Printf("%s%s", formattedError, "the result is not a pointer to http.Response")
		os.Exit(1)
	}
	// Read the content of http.Response.Body
	defer response.Body.Close() // Ensure to close the ReadCloser
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, response.Body); err != nil {
		fmt.Printf("%s%v", formattedError, err)
		os.Exit(1)
	}

	if response.StatusCode >= 400 {
		ErrorHandler(resource.Kind, name, buf.String())
		os.Exit(1)
	}

	// Check if the content is an array or an object
	var res interface{}
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		fmt.Printf("%s%v", formattedError, err)
		os.Exit(1)
	}
	output(resource, []interface{}{res}, outputFormat)
}

func (resource Resource) ListFn(options map[string]string) {
	formattedError := fmt.Sprintf("Resource %s error: ", resource.Kind)
	ctx := context.Background()
	// Use reflect to call the function
	funcValue := reflect.ValueOf(resource.List)
	if funcValue.Kind() != reflect.Func {
		fmt.Println("fn is not a valid function")
		os.Exit(1)
	}
	// Create a slice for the arguments
	fnargs := []reflect.Value{reflect.ValueOf(ctx)} // Add the context

	// Handle the options if the resource has ListParamsType
	if resource.ListParamsType != nil {
		paramsValue := retrieveListParams(resource.ListParamsType, options)
		fnargs = append(fnargs, paramsValue)
	}

	// Call the function with the arguments
	results := funcValue.Call(fnargs)
	// Handle the results based on your needs
	if len(results) <= 1 {
		return
	}
	if err, ok := results[1].Interface().(error); ok && err != nil {
		fmt.Printf("%s%v", formattedError, err)
		os.Exit(1)
	}
	// Check if the first result is a pointer to http.Response
	response, ok := results[0].Interface().(*http.Response)
	if !ok {
		fmt.Printf("%s%s", formattedError, "the result is not a pointer to http.Response")
		os.Exit(1)
	}
	// Read the content of http.Response.Body
	defer response.Body.Close() // Ensure to close the ReadCloser
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, response.Body); err != nil {
		fmt.Printf("%s%v", formattedError, err)
		os.Exit(1)
	}
	if response.StatusCode >= 400 {
		ErrorHandler(resource.Kind, "", buf.String())
		os.Exit(1)
	}

	// Check if the content is an array or an object
	var slices []interface{}
	if err := json.Unmarshal(buf.Bytes(), &slices); err != nil {
		fmt.Printf("%s%v", formattedError, err)
		os.Exit(1)
	}
	// Check the output format
	output(resource, slices, outputFormat)
}
