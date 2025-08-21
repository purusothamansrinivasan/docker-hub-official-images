package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"

	"github.com/docker-hub-api/mcp-server/config"
	"github.com/docker-hub-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func PostnamespacesdeleteimagesHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		namespaceVal, ok := args["namespace"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: namespace"), nil
		}
		namespace, ok := namespaceVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: namespace"), nil
		}
		// Create properly typed request body using the generated schema
		var requestBody models.PostNamespacesDeleteImagesRequest
		
		// Optimized: Single marshal/unmarshal with JSON tags handling field mapping
		if argsJSON, err := json.Marshal(args); err == nil {
			if err := json.Unmarshal(argsJSON, &requestBody); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to convert arguments to request type: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal arguments: %v", err)), nil
		}
		
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to encode request body", err), nil
		}
		url := fmt.Sprintf("%s/v2/namespaces/%s/delete-images", cfg.BaseURL, namespace)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// No authentication required for this endpoint
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.PostNamespacesDeleteImagesResponseSuccess
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreatePostnamespacesdeleteimagesTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_v2_namespaces_namespace_delete-images",
		mcp.WithDescription("Delete images"),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace of the repository.")),
		mcp.WithString("active_from", mcp.Description("Input parameter: Sets the time from which an image must have been pushed or pulled to\nbe counted as active.\n\nDefaults to 1 month before the current time.\n")),
		mcp.WithBoolean("dry_run", mcp.Description("Input parameter: If `true` then will check and return errors and unignored warnings for the deletion request but will not delete any images.")),
		mcp.WithArray("ignore_warnings", mcp.Description("Input parameter: Warnings to ignore. If a warning is not ignored then no deletions will happen and the \nwarning is returned in the response.\n\nThese warnings include:\n\n- is_active: warning when attempting to delete an image that is marked as active.\n- current_tag: warning when attempting to delete an image that has one or more current \ntags in the repository.\n\nWarnings can be copied from the response to the request.\n")),
		mcp.WithArray("manifests", mcp.Description("Input parameter: Image manifests to delete.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    PostnamespacesdeleteimagesHandler(cfg),
	}
}
