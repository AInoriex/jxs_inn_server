//@Author	AInoriex
//@Date		2025.04.23
//@Desc		此文件保存公共的结构
// 配置`config.yaml`需要读取的配置

package config

// 配置文件映射
// 绑定关系 `config.yaml` key → `utils/config/model.go` ConfigMap.key → `utils/config/model.go` Struct
var ConfigMap = map[string]interface{}{
	"common_config": CommonConfig,
	"db_config":     DbConfig,
}

var CommonConfig = new(CommonConf)
var DbConfig = new(DbConf)

type Mysql struct {
	Host   string `mapstructure:"host"`
	Db     string `mapstructure:"db"`
	MaxCon int    `mapstructure:"max_con"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

type LogConfig struct {
	SavePath string `mapstructure:"save_path"`
}

// 华为Obs
type HuaweiOBS struct {
	Url       string `mapstructure:"url"`
	SecretId  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
	Env       string `mapstructure:"env"`
	Cdn       string `mapstructure:"cdn"`
	Bucket    string `mapstructure:"bucket"`
}

// http服务端
type HttpServerConf struct {
	Addr    string `mapstructure:"addr"`
	Timeout int64  `mapstructure:"timeout"`
	Network string `mapstructure:"network"`
}

// 通用配置
type CommonConf struct {
	AppName    string                 `mapstructure:"app_name"`    // 应用名称
	Env        string                 `mapstructure:"env"`         // 环境
	Log        LogConfig              `mapstructure:"log"`         // 日志配置
	OpenDbLog  bool                   `mapstructure:"open_db_log"` // 是否开启数据库日志
	HuaweiOBS  HuaweiOBS              `mapstructure:"huawei_obs"`  // 华为obs配置
	ApiHost    string                 `mapstructure:"api_host"`    // api域名
	HttpServer HttpServerConf         `mapstructure:"http_server"` // http服务
	JwtSecret  string                 `mapstructure:"jwt_secret"`  // jwt密钥
	YltAccount map[string]interface{} `mapstructure:"ylt_account"` // ylt账号
}

// 数据库配置
type DbConf struct {
	Mysql Mysql `mapstructure:"mysql"`
	Redis Redis `mapstructure:"redis"`
}
