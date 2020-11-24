package entities

type MerchantInfo struct {
	MerchantID   float64
	MerchantCode string
	MerchantName string
}

type Cluster struct {
	ClusterID          int    `json:"cluster_id"`
	ClusterMID         string `json:"cluster_mid"`
	ClusterName        string `json:"cluster_name"`
	ClusterLogo        string `json:"cluster_logo"`
	ClusterDescription string `json:"cluster_description"`
	Site               []Site `json:"site"`
}

type Site struct {
	SiteID        int     `json:"site_id"`
	SiteMID       string  `json:"site_mid"`
	SiteName      string  `json:"site_name"`
	SiteEstimated string  `json:"site_estimated"`
	SiteLogo      string  `json:"site_logo"`
	SiteLat       string  `json:"site_latitude"`
	SiteLong      string  `json:"site_longitude"`
	Trf           SiteTrf `json:"tariff"`
}

type SiteTrf struct {
	Adult []SiteTrfModel `json:"adult"`
	Child []SiteTrfModel `json:"child"`
}

type SiteTrfModel struct {
	Trf_id            int     `json:"trf_id"`
	Trf_code          string  `json:"trf_code"`
	Trfftype_name     string  `json:"trf_type"`
	Trf_name          string  `json:"trf_name"`
	Trf_value         float32 `json:"trf_value"`
	Trf_currency_code string  `json:"trf_currency_code"`
	Mtick_name        string  `json:"tick_name"`
}

func (SiteTrfModel) TableName() string {
	return "master_tariff"
}
