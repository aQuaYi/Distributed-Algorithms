package raft

// LogEntry is log entry
type LogEntry struct {
	LogIndex int         // raft.logs 会被压缩裁剪，需要保存此 log 在原本的索引号
	LogTerm  int         // LEADER 在生成此 log 时的 LEADER.currentTerm
	Command  interface{} // 具体的命令内容
}
