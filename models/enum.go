package models

type Status string

const (
	Active          Status = "Active"
	Suspended       Status = "Suspended"
	Deleted         Status = "Deleted"
	Inactive        Status = "Inactive"
	PendindApproval Status = "Pending Approval"
	Rejected        Status = "Rejected"
	Approved        Status = "Approved"
	Submitted       Status = "Submitted"
	LoggedOut       Status = "Logged Out"
)

type NotificationTopics string

const (
	NotificationTopic NotificationTopics = "notifications"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypeEmail NotificationType = "email"
	NotificationTypePush  NotificationType = "push"
)
