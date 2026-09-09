package rules

// CommonTypeTokens maps Azure Terraform resource types to standard provider abbreviations
// and common shorthand synonyms used by DevOps teams.
var CommonTypeTokens = map[string][]string{
	// Compute & Containers
	"azurerm_virtual_machine":               {"vm", "virtualmachine"},
	"azurerm_linux_virtual_machine":         {"vm", "lvm", "linuxvm"},
	"azurerm_windows_virtual_machine":       {"vm", "wvm", "winvm"},
	"azurerm_virtual_machine_scale_set":     {"vmss", "scaleset"},
	"azurerm_kubernetes_cluster":            {"aks", "k8s", "kubernetes"},
	"azurerm_container_registry":            {"acr", "registry"},
	"azurerm_container_group":               {"aci", "containergroup"},
	"azurerm_container_app":                 {"ca", "containerapp"},
	"azurerm_container_app_environment":     {"cae", "containerappenv"},

	// Networking
	"azurerm_resource_group":                {"rg", "resourcegroup", "resgroup", "group"},
	"azurerm_virtual_network":               {"vnet", "vn", "virtualnetwork", "network"},
	"azurerm_subnet":                        {"snet", "subnet", "sub"},
	"azurerm_network_security_group":        {"nsg", "secgroup", "securitygroup"},
	"azurerm_network_security_rule":         {"nsgrule", "secrule"},
	"azurerm_route_table":                   {"rt", "routetable"},
	"azurerm_network_interface":             {"nic", "netint"},
	"azurerm_public_ip":                     {"pip", "ip", "publicip"},
	"azurerm_nat_gateway":                   {"ng", "natgw", "natgateway"},
	"azurerm_application_gateway":           {"agw", "appgw", "appgateway"},
	"azurerm_lb":                            {"lb", "loadbalancer"},
	"azurerm_firewall":                      {"fw", "afw", "firewall"},
	"azurerm_dns_zone":                      {"dns", "dnszone"},
	"azurerm_private_dns_zone":              {"pdns", "privatedns"},
	"azurerm_private_endpoint":              {"pe", "pep", "privateendpoint"},
	"azurerm_express_route_circuit":         {"erc", "expressroute"},

	// Storage & Databases
	"azurerm_storage_account":               {"st", "sa", "storage", "storageaccount"},
	"azurerm_storage_container":             {"container", "stcontainer"},
	"azurerm_storage_share":                 {"share", "fileshare"},
	"azurerm_mssql_server":                  {"sql", "sqlserver"},
	"azurerm_mssql_database":                {"sqldb", "db", "database"},
	"azurerm_postgresql_flexible_server":    {"psql", "postgres", "postgresql"},
	"azurerm_mysql_flexible_server":         {"mysql"},
	"azurerm_cosmosdb_account":              {"cosmos", "cosmosdb"},
	"azurerm_redis_cache":                   {"redis", "cache"},

	// Security & Identity
	"azurerm_key_vault":                     {"kv", "keyvault", "vault"},
	"azurerm_key_vault_secret":              {"secret", "kvsecret"},
	"azurerm_user_assigned_identity":        {"uai", "msi", "identity"},
	"azurerm_bastion_host":                  {"bas", "bastion"},

	// Web & App Services
	"azurerm_service_plan":                  {"asp", "appplan", "serviceplan"},
	"azurerm_linux_web_app":                 {"app", "webapp"},
	"azurerm_windows_web_app":               {"app", "webapp"},
	"azurerm_linux_function_app":            {"func", "funcapp", "functionapp"},
	"azurerm_windows_function_app":          {"func", "funcapp", "functionapp"},
	"azurerm_logic_app_workflow":            {"logic", "logicapp"},
	"azurerm_api_management":                {"apim", "apimanagement"},

	// Integration, Monitoring & Analytics
	"azurerm_log_analytics_workspace":       {"law", "workspace", "loganalytics"},
	"azurerm_application_insights":          {"appi", "appinsights", "insights"},
	"azurerm_eventhub_namespace":            {"evh", "eventhub"},
	"azurerm_servicebus_namespace":          {"sb", "servicebus"},
	"azurerm_data_factory":                  {"adf", "datafactory"},
	"azurerm_synapse_workspace":             {"syn", "synapse"},
	"azurerm_databricks_workspace":          {"dbw", "databricks"},
}