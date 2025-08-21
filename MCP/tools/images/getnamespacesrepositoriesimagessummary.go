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

func GetnamespacesrepositoriesimagessummaryHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		if val, ok := args["active_from"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("active_from=%v", val))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/v2/namespaces/%s/repositories/%s/images-summary%s", cfg.BaseURL, namespace, repository, queryString)
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
		var result models.GetNamespaceRepositoryImagesSummaryResponse
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

func CreateGetnamespacesrepositoriesimagessummaryTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_v2_namespaces_namespace_repositories_repository_images-summary",
		mcp.WithDescription("Get summary of repository's images"),
		mcp.WithString("namespace", mcp.Required(), mcp.Description("Namespace of the repository.")),
		mcp.WithString("repository", mcp.Required(), mcp.Description("Name of the repository.")),
		mcp.WithString("active_from", mcp.Description("Sets the time from which an image must have been pushed or pulled to\nbe counted as active.\n\nDefaults to 1 month before the current time.\n")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetnamespacesrepositoriesimagessummaryHandler(cfg),
	}
}
