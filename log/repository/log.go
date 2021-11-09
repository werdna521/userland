package repository

type AuditLog struct {
	ID        string
	UserID    string
	RemoteIP  string
	AuditType string
}
