package models

type Subscriber struct {
	Uuid          string `json:"uuid" gorm:"primary_key"`
	Email         string `json:"email"`
	PolicyUpdates bool   `json:"policyUpdates"`
	Incidents     bool   `json:"incidents"`
	NewFeatures   bool   `json:"newFeatures"`
	News          bool   `json:"news"`
	Other         bool   `json:"other"`
}
