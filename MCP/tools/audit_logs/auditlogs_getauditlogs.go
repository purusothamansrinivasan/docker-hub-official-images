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

func Auditlogs_getauditlogsHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		accountVal, ok := args["account"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: account"), nil
		}
		account, ok := accountVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: account"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["action"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("action=%v", val))
		}
		if val, ok := args["name"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("name=%v", val))
		}
		if val, ok := args["actor"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("actor=%v", val))
		}
		if val, ok := args["from"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("from=%v", val))
		}
		if val, ok := args["to"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("to=%v", val))
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
		url := fmt.Sprintf("%s/v2/auditlogs/%s%s", cfg.BaseURL, account, queryString)
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
		var result models.GetAuditLogsResponse
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

func CreateAuditlogs_getauditlogsTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_v2_auditlogs_account",
		mcp.WithDescription("Returns list of audit log  events."),
		mcp.WithString("account", mcp.Required(), mcp.Description("Namespace to query audit logs for.")),
		mcp.WithString("action", mcp.Description("action name one of [\"repo.tag.push\", ...]. Optional parameter to filter specific audit log actions.")),
		mcp.WithString("name", mcp.Description("name. Optional parameter to filter audit log events to a specific name. For repository events, this is the name of the repository. For organization events, this is the name of the organization. For team member events, this is the username of the team member.")),
		mcp.WithString("actor", mcp.Description("actor name. Optional parameter to filter audit log events to the specific user who triggered the event.")),
		mcp.WithString("from", mcp.Description("Start of the time window you wish to query audit events for.")),
		mcp.WithString("to", mcp.Description("End of the time window you wish to query audit events for.")),
		mcp.WithNumber("page", mcp.Description("page - specify page number. Page number to get.")),
		mcp.WithNumber("page_size", mcp.Description("page_size - specify page size. Number of events to return per page.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Auditlogs_getauditlogsHandler(cfg),
	}
}
