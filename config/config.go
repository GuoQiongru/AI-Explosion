package config

var Secret = "tiktok"

var OneDayOfHours = 60 * 60 * 24

// VideoCount 每次获取视频流的数量
const VideoCount = 5

// 存储的图片和视频的链接
const UrlPrefix = "http://47.113.148.197/" //服务器地址

// ConConfig ftp服务器地址
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
