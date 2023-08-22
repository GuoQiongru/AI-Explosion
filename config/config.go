package config

var Secret = "tiktok"

var OneDayOfHours = 60 * 60 * 24
var OneMonth = 60 * 60 * 24 * 30

// VideoCount 每次获取视频流的数量
const VideoCount = 5

// 存储的图片和视频的链接
const UrlPrefix = "http://47.113.148.197/" //服务器地址

// ConConfig ftp服务器
const ConConfig = "47.113.148.197:21"
const FtpUser = "ftpuser"
const FtpPsw = "123456"
const HeartbeatTime = 2 * 60

// HostSSH SSH配置
const HostSSH = "47.113.148.197"
const UserSSH = "root"
const PasswordSSH = "AI-Explosion"
const TypeSSH = "password"
const PortSSH = 22
const MaxMsgCount = 100
const SSHHeartbeatTime = 10 * 60

const IsLike = 0     //点赞的状态
const Unlike = 1     //取消赞的状态
const LikeAction = 1 //点赞的行为
const Attempts = 3   //操作数据库的最大尝试次数

const RedisAddr = "47.113.148.197:6379"
const RedisPsw = "123456"
const DefaultRedisValue = -1 //防脏读
