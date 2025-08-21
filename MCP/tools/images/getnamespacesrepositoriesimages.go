package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/docker-hub-api/mcp-server/config"
	"github.com/docker-hub-api/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetnamespacesrepositoriesimagesHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		repositoryVal, ok := args["repository"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: repository"), nil
		}
		repository, ok := repositoryVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: repository"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["status"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("status=%v", val))
		}
		if val, ok := args["currently_tagged"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("currently_tagged=%v", val))
		}
		if val, ok := args["ordering"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("ordering=%v", val))
		}
		if val, ok := args["active_from"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("active_from=%v", val))
		}
		if val, ok := args["page"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("page=%v", val))
		}
		if val, ok := args["page_size"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("page_size=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/v2/namespaces/%s/repositories/%s/images%s", cfg.BaseURL, namespace, repository, queryString)
		req, err := http.NewRequest("GET", url, nil)
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
		var result models.GetNamespaceRepositoryImagesResponse
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

func CreateGetnamespacesrepositoriesimagesTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_v2_namespaces_namespace_repositories_repository_images",
		mcp.WithDescription("Get details of repository's images"),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace of the repository.")),
		mcp.WithString("repository", mcp.Required(), mcp.Description("Name of the repository.")),
		mcp.WithString("status", mcp.Description("Filters to only show images of this status.")),
		mcp.WithBoolean("currently_tagged", mcp.Description("Filters to only show images with:\n- `true`: at least 1 current tag.\n- `false`: no current tags.\n")),
		mcp.WithString("ordering", mcp.Description("Orders the results by this property.\n\nPrefixing with `-` sorts by descending order.\n")),
		mcp.WithString("active_from", mcp.Description("Sets the time from which an image must have been pushed or pulled to\nbe counted as active.\n\nDefaults to 1 month before the current time.\n")),
		mcp.WithNumber("page", mcp.Description("Page number to get. Defaults to 1.")),
		mcp.WithNumber("page_size", mcp.Description("Number of images to get per page. Defaults to 10. Max of 100.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetnamespacesrepositoriesimagesHandler(cfg),
	}
}
