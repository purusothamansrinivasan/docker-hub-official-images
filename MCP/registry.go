package main

import (
	"github.com/docker-hub-api/mcp-server/config"
	"github.com/docker-hub-api/mcp-server/models"
	tools_audit_logs "github.com/docker-hub-api/mcp-server/tools/audit_logs"
	tools_images "github.com/docker-hub-api/mcp-server/tools/images"
	tools_authentication "github.com/docker-hub-api/mcp-server/tools/authentication"
	tools_repositories "github.com/docker-hub-api/mcp-server/tools/repositories"
	tools_access_tokens "github.com/docker-hub-api/mcp-server/tools/access_tokens"
	tools_org_settings "github.com/docker-hub-api/mcp-server/tools/org_settings"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_audit_logs.CreateAuditlogs_getauditlogsTool(cfg),
		tools_images.CreateGetnamespacesrepositoriesimagesTool(cfg),
		tools_authentication.CreatePostusers2faloginTool(cfg),
		tools_repositories.CreateHead_v2_namespaces_namespace_repositories_repository_tags_tagTool(cfg),
		tools_repositories.CreateGet_v2_namespaces_namespace_repositories_repository_tags_tagTool(cfg),
		tools_audit_logs.CreateAuditlogs_getauditactionsTool(cfg),
		tools_images.CreateGetnamespacesrepositoriesimagestagsTool(cfg),
		tools_access_tokens.CreatePatch_v2_access_tokens_uuidTool(cfg),
		tools_access_tokens.CreateDelete_v2_access_tokens_uuidTool(cfg),
		tools_access_tokens.CreateGet_v2_access_tokens_uuidTool(cfg),
		tools_repositories.CreateGet_v2_namespaces_namespace_repositories_repository_tagsTool(cfg),
		tools_repositories.CreateHead_v2_namespaces_namespace_repositories_repository_tagsTool(cfg),
		tools_org_settings.CreateGet_v2_orgs_name_settingsTool(cfg),
		tools_org_settings.CreatePut_v2_orgs_name_settingsTool(cfg),
		tools_authentication.CreatePostusersloginTool(cfg),
		tools_images.CreateGetnamespacesrepositoriesimagessummaryTool(cfg),
		tools_access_tokens.CreateGet_v2_access_tokensTool(cfg),
		tools_access_tokens.CreatePost_v2_access_tokensTool(cfg),
		tools_images.CreatePostnamespacesdeleteimagesTool(cfg),
	}
}
