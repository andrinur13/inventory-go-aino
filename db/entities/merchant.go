package entities

type MerchantInfo struct {
	MerchantID   float64
	MerchantCode string
	MerchantName string
}

type Cluster struct {
	ClusterID   int    `json:"cluster_id"`
	ClusterMID  string `json:"cluster_mid"`
	ClusterName string `json:"cluster_name"`
	ClusterLogo string `json:"cluster_logo"`
	Site        []Site `json:"site"`
}

type Site struct {
	SiteID   int    `json:"site_id"`
	SiteMID  string `json:"site_mid"`
	SiteName string `json:"site_name"`
	SiteLogo string `json:"site_logo"`
}
