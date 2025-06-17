package models

type Subscriber struct {
	Uuid          string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"uuid,omitempty"`
	Email         string `json:"email"`
	PolicyUpdates bool   `json:"policyUpdates"`
	Incidents     bool   `json:"incidents"`
	NewFeatures   bool   `json:"newFeatures"`
	News          bool   `json:"news"`
	Others        bool   `json:"others"`
	CreatedAt     string `gorm:"type:datetime;default:now()" json:"createdAt,omitempty"`
}
