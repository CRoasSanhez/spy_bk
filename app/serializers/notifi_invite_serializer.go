package serializers

import "spyc_backend/app/models"

type NotifyInviteSerializer struct {
	Notifications []models.Notification `json:"notifications"`
	Invitations []models.Invitation `json:"invitations"`
}

