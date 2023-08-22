package config

type Config struct {
	WorksAddr       []string
	PomeloAddress   string
	CustomSendFiles []string // 自定义发送使用的数据配置目录
	Timeout         int      `json:",default=20"` // 默认连接超时 20s
}
