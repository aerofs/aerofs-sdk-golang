package aerofssdk

// Structures used when communicating with an AeroFS Appliance

type ParentPath struct {
	Folders []Folder `json:"folders"`
}

type Children struct {
	Folders []Folder `json:"folders"`
	Files   []File   `json:"files"`
}

type SFGroupMember struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type SFPendingMember struct {
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name,omitempty"`
	LastName    string   `json:"last_name,omitempty"`
	Inviter     string   `json:"invited_by,omitempty"`
	Permissions []string `json:"permissions"`
	Note        string   `json:"note"`
}

type Invitee struct {
	EmailTo    string `json:"email_to"`
	EmailFrom  string `json:"email_from"`
	SignupCode string `json:"signup_code,omitempty"`
}

type Invitation struct {
	Id          string   `json:"shared_id"`
	Name        string   `json:"shared_name"`
	Inviter     string   `json:"invited_by"`
	Permissions []string `json:"permissions"`
}
