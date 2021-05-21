package entities

type InboxNotificationModel struct {
	Inbox_id        		int          	`json:"inbox_id" gorm:"primary_key"`
	Agent_name				string			`json:"agent_name" gorm:"-"`
	Agent_group_name		string			`json:"agent_group_name" gorm:"-"`
	Inbox_created_at        string  	    `json:"inbox_created_at" gorm:"-"`
	Inbox_image_url			string			`json:"inbox_image_url" gorm:"-"`
	Inbox_title				string			`json:"inbox_title" gorm:"-"`
	Inbox_short_desc		string			`json:"inbox_short_desc" gorm:"-"`
	Inbox_full_desc			string			`json:"inbox_full_desc" gorm:"-"`
	// InboxExtras				*InboxExtras	`json:"inbox_extras" gorm:"-"`
}

func (InboxNotificationModel) TableName() string {
	return "inbox_notification"
}

type InboxExtras struct {
	CsPhoneNo  		string `json:"cs_phone_no"`
	CsFaxNo    		string `json:"cs_fax_no"`
	CsWaNo     		string `json:"cs_wa_no"`
	CsEmail    		string `json:"cs_email"`
	PrivacyPolicy   string `json:"privacy_policy"`
	TermCondition   string `json:"term_condition"`
}
