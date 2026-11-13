// Package buildinfo 保存编译时注入的版本信息。
package buildinfo

var (
	// Version 是发布版本，由编译参数注入。
	Version = "dev"
	// Commit 是 Git 提交号，由编译参数注入。
	Commit = "none"
	// Time 是构建时间，由编译参数注入。
	Time = "unknown"
)
