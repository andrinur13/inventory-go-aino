package entities

type MconfigModel struct {
	// Mconfig_id        int          	`json:"mconfig_id" gorm:"primary_key"`
	// Mconfig_src_type	 int		  	`json:"mconfig_src_type"`
	Cs_phone_no			 string			`json:"-" gorm:"-"`
	Cs_fax_no			 string			`json:"-" gorm:"-"`
	Cs_wa_no			 string			`json:"-" gorm:"-"`
	Cs_email			 string			`json:"-" gorm:"-"`
	Mconfig_value        string			`json:"-"`
	MconfigValue         *MconfigValue 	`json:"mconfig_value" gorm:"-"`
	Created_at           string       	`json:"-"`
}

func (MconfigModel) TableName() string {
	return "mobile_config"
}

type MconfigValue struct {
	CsPhoneNo  		string `json:"cs_phone_no"`
	CsFaxNo    		string `json:"cs_fax_no"`
	CsWaNo     		string `json:"cs_wa_no"`
	CsEmail    		string `json:"cs_email"`
	PrivacyPolicy   string `json:"privacy_policy"`
	TermCondition   string `json:"term_condition"`
	ChildAge  		string `json:"child_age"`
	AdultAge  		string `json:"adult_age"`
}

func (MconfigValue) TableName() string {
	return "mobile_config"
}