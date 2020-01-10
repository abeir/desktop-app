package config

type LoadState int

const (
	Unload LoadState = -1 // 未加载
	Loading LoadState = 0 // 加载中
	Loaded LoadState = 1  // 已加载
)

const (
	//环境变量中配置文件路径
	ApplicationEnvVar = "TRAN_TICKET_APP"
	// api.yml文件路径在环境中的变量名
	ApiEnvVar = "TRAN_TICKET_API"
)