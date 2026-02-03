package friendlycaptcha

import "github.com/guregu/null/v6"

// RiskScore represents a risk score value ranging from 1 to 5.
//   - 0: Unknown or missing
//   - 1: Very low risk
//   - 2: Low risk
//   - 3: Medium risk
//   - 4: High risk
//   - 5: Very high risk
type RiskScore uint8

const (
	// RiskScoreUnknown represents an unknown or missing risk score.
	RiskScoreUnknown RiskScore = 0
	// RiskScoreVeryLow represents a very low risk score (1/5).
	RiskScoreVeryLow RiskScore = 1
	// RiskScoreLow represents a low risk score (2/5).
	RiskScoreLow RiskScore = 2
	// RiskScoreMedium represents a medium risk score (3/5).
	RiskScoreMedium RiskScore = 3
	// RiskScoreHigh represents a high risk score (4/5).
	RiskScoreHigh RiskScore = 4
	// RiskScoreVeryHigh represents a very high risk score (5/5).
	RiskScoreVeryHigh RiskScore = 5
)

// RiskIntelligenceData contains all risk intelligence information.
//
// Field availability depends on enabled modules.
type RiskIntelligenceData struct {
	// RiskScores from various signals, these summarize the risk intelligence assessment.
	//
	// Available when the Risk Scores module is enabled.
	// Null when the Risk Scores module is not enabled.
	RiskScores null.Value[RiskScoresData] `json:"risk_scores"`

	// Network contains network-related risk intelligence.
	Network NetworkData `json:"network"`

	// Client contains client/device risk intelligence.
	Client ClientData `json:"client"`
}

// RiskScoresData summarizes the entire risk intelligence assessment into scores per category.
//
// Available when the Risk Scores module is enabled for your account.
// Null when the Risk Scores module is not enabled for your account.
type RiskScoresData struct {
	// Overall risk score combining all signals.
	Overall RiskScore `json:"overall"`

	// Network-related risk score. Captures likelihood of automation/malicious activity based on
	// IP address, ASN, reputation, geolocation, past abuse from this network, and other network signals.
	Network RiskScore `json:"network"`

	// Browser-related risk score. Captures likelihood of automation, malicious activity or browser spoofing based on
	// user agent consistency, automation traces, past abuse, and browser characteristics.
	Browser RiskScore `json:"browser"`
}

// NetworkAutonomousSystemData contains information about the AS that owns the IP.
//
// Available when the IP Intelligence module is enabled for your account.
// Null when the IP Intelligence module is not enabled for your account.
type NetworkAutonomousSystemData struct {
	// Number is the Autonomous System Number (ASN) identifier.
	// Example: 3209 for Vodafone GmbH
	Number int `json:"number"`

	// Name of the autonomous system. This is usually a short name or handle.
	// Example: "VODANET"
	Name string `json:"name"`

	// Company is the organization name that owns the ASN.
	// Example: "Vodafone GmbH"
	Company string `json:"company"`

	// Description of the company that owns the ASN.
	// Example: "Provides mobile and fixed broadband and telecommunication services to consumers and businesses."
	Description string `json:"description"`

	// Domain name associated with the ASN.
	// Example: "vodafone.de"
	Domain string `json:"domain"`

	// Country is the two-letter ISO 3166-1 alpha-2 country code where the ASN is registered.
	// Example: "DE"
	Country string `json:"country"`

	// RIR is the Regional Internet Registry that allocated the ASN.
	// Example: "RIPE"
	RIR string `json:"rir"`

	// Route is the IP route associated with the ASN in CIDR notation.
	// Example: "88.64.0.0/12"
	Route string `json:"route"`

	// Type of the autonomous system.
	// Example: "isp"
	Type string `json:"type"`
}

// NetworkGeolocationCountryData contains detailed country data.
type NetworkGeolocationCountryData struct {
	// ISO2 is the two-letter ISO 3166-1 alpha-2 country code.
	// Example: "DE"
	ISO2 string `json:"iso2"`

	// ISO3 is the three-letter ISO 3166-1 alpha-3 country code.
	// Example: "DEU"
	ISO3 string `json:"iso3"`

	// Name is the English name of the country.
	// Example: "Germany"
	Name string `json:"name"`

	// NameNative is the native name of the country.
	// Example: "Deutschland"
	NameNative string `json:"name_native"`

	// Region is the major world region.
	// Example: "Europe"
	Region string `json:"region"`

	// Subregion is the more specific world region.
	// Example: "Western Europe"
	Subregion string `json:"subregion"`

	// Currency is the ISO 4217 currency code.
	// Example: "EUR"
	Currency string `json:"currency"`

	// CurrencyName is the full name of the currency.
	// Example: "Euro"
	CurrencyName string `json:"currency_name"`

	// PhoneCode is the international dialing code.
	// Example: "49"
	PhoneCode string `json:"phone_code"`

	// Capital is the name of the capital city.
	// Example: "Berlin"
	Capital string `json:"capital"`
}

// NetworkGeolocationData contains geographic location of the IP address.
//
// Available when the IP Intelligence module is enabled.
// Null when the IP Intelligence module is not enabled.
type NetworkGeolocationData struct {
	// Country information.
	Country NetworkGeolocationCountryData `json:"country"`

	// City name. Empty string if unknown.
	// Example: "Eschborn"
	City string `json:"city"`

	// State, region, or province. Empty string if unknown.
	// Example: "Hessen"
	State string `json:"state"`
}

// NetworkAbuseContactData contains contact details for reporting abuse.
//
// Available when the IP Intelligence module is enabled.
// Null when the IP Intelligence module is not enabled.
type NetworkAbuseContactData struct {
	// Address is the postal address of the abuse contact.
	// Example: "Vodafone GmbH, Campus Eschborn, Duesseldorfer Strasse 15, D-65760 Eschborn, Germany"
	Address string `json:"address"`

	// Name of the abuse contact person or team.
	// Example: "Vodafone Germany IP Core Backbone"
	Name string `json:"name"`

	// Email is the abuse contact email address.
	// Example: "abuse.de@vodafone.com"
	Email string `json:"email"`

	// Phone is the abuse contact phone number.
	// Example: "+49 6196 52352105"
	Phone string `json:"phone"`
}

// NetworkAnonymizationData contains detection of VPNs, proxies, and anonymization services.
//
// Available when the Anonymization Detection module is enabled.
// Null when the Anonymization Detection module is not enabled.
type NetworkAnonymizationData struct {
	// VPNScore is the likelihood that the IP is from a VPN service.
	VPNScore RiskScore `json:"vpn_score"`

	// ProxyScore is the likelihood that the IP is from a proxy service.
	ProxyScore RiskScore `json:"proxy_score"`

	// Tor indicates whether the IP is a Tor exit node.
	Tor bool `json:"tor"`

	// ICloudPrivateRelay indicates whether the IP is from iCloud Private Relay.
	ICloudPrivateRelay bool `json:"icloud_private_relay"`
}

// NetworkData contains information about the network.
type NetworkData struct {
	// IP is the IP address used when requesting the challenge.
	// Example: "88.64.4.22"
	IP string `json:"ip"`

	// AS contains Autonomous System information.
	//
	// Available when the IP Intelligence module is enabled.
	// Null when the IP Intelligence module is not enabled.
	AS null.Value[NetworkAutonomousSystemData] `json:"as"`

	// Geolocation information.
	//
	// Available when the IP Intelligence module is enabled.
	// Null when the IP Intelligence module is not enabled.
	Geolocation null.Value[NetworkGeolocationData] `json:"geolocation"`

	// AbuseContact is the abuse contact information.
	//
	// Available when the IP Intelligence module is enabled.
	// Null when the IP Intelligence module is not enabled.
	AbuseContact null.Value[NetworkAbuseContactData] `json:"abuse_contact"`

	// Anonymization contains IP masking/anonymization information.
	//
	// Available when the Anonymization Detection module is enabled.
	// Null when the Anonymization Detection module is not enabled.
	Anonymization null.Value[NetworkAnonymizationData] `json:"anonymization"`
}

// ClientTimeZoneData contains IANA time zone data.
//
// Available when the Browser Identification module is enabled.
// Null when the Browser Identification module is not enabled.
type ClientTimeZoneData struct {
	// Name is the IANA time zone name reported by the browser.
	// Example: "America/New_York" or "Europe/Berlin"
	Name string `json:"name"`

	// CountryISO2 is the two-letter ISO 3166-1 alpha-2 country code derived from the time zone.
	// "XU" if timezone is missing or cannot be mapped to a country (e.g., "Etc/UTC").
	// Example: "US" or "DE"
	CountryISO2 string `json:"country_iso2"`
}

// ClientBrowserData contains detected browser details.
//
// Available when the Browser Identification module is enabled.
// Null when the Browser Identification module is not enabled.
type ClientBrowserData struct {
	// ID is the unique browser identifier. Empty string if browser could not be identified.
	// Example: "firefox", "chrome", "chrome_android", "edge", "safari", "safari_ios", "webview_ios"
	ID string `json:"id"`

	// Name is the human-readable browser name. Empty string if browser could not be identified.
	// Example: "Firefox", "Chrome", "Edge", "Safari", "Safari on iOS", "WebView on iOS"
	Name string `json:"name"`

	// Version is the browser version name. Assumed to be the most recent release matching the signature if exact version unknown. Empty if unknown.
	// Example: "146.0" or "16.5"
	Version string `json:"version"`

	// ReleaseDate is the release date of the browser version in "YYYY-MM-DD" format. Empty string if unknown.
	// Example: "2026-01-28"
	ReleaseDate string `json:"release_date"`
}

// ClientBrowserEngineData contains detected rendering engine details.
//
// Available when the Browser Identification module is enabled.
// Null when the Browser Identification module is not enabled.
type ClientBrowserEngineData struct {
	// ID is the unique rendering engine identifier. Empty string if engine could not be identified.
	// Example: "gecko", "blink", "webkit"
	ID string `json:"id"`

	// Name is the human-readable engine name. Empty string if engine could not be identified.
	// Example: "Gecko", "Blink", "WebKit"
	Name string `json:"name"`

	// Version is the rendering engine version. Assumed to be the most recent release matching the signature if exact version unknown. Empty if unknown.
	// Example: "146.0" or "16.5"
	Version string `json:"version"`
}

// ClientDeviceData contains detected device details.
//
// Available when the Browser Identification module is enabled.
// Null when the Browser Identification module is not enabled.
type ClientDeviceData struct {
	// Type is the device type.
	// Example: "desktop", "mobile", "tablet"
	Type string `json:"type"`

	// Brand is the device brand.
	// Example: "Apple", "Samsung", "Google"
	Brand string `json:"brand"`

	// Model is the device model name.
	// Example: "iPhone 17", "Galaxy S21 (SM-G991B)", "Pixel 10"
	Model string `json:"model"`
}

// ClientOSData contains detected OS details.
//
// Available when the Browser Identification module is enabled.
// Null when the Browser Identification module is not enabled.
type ClientOSData struct {
	// ID is the unique operating system identifier. Empty string if OS could not be identified.
	// Example: "windows", "macos", "ios", "android", "linux"
	ID string `json:"id"`

	// Name is the human-readable operating system name. Empty string if OS could not be identified.
	// Example: "Windows", "macOS", "iOS", "Android", "Linux"
	Name string `json:"name"`

	// Version is the operating system version.
	// Example: "10", "11.2.3", "14.4"
	Version string `json:"version"`
}

// TLSSignatureData contains TLS client hello signatures.
//
// Available when the Bot Detection module is enabled.
// Null when the Bot Detection module is not enabled.
type TLSSignatureData struct {
	// JA3 is the JA3 hash.
	// Example: "d87a30a5782a73a83c1544bb06332780"
	JA3 string `json:"ja3"`

	// JA3N is the JA3N hash.
	// Example: "28ecc2d2875b345cecbb632b12d8c1e0"
	JA3N string `json:"ja3n"`

	// JA4 is the JA4 signature.
	// Example: "t13d1516h2_8daaf6152771_02713d6af862"
	JA4 string `json:"ja4"`
}

// ClientAutomationKnownBotData contains detected known bot details.
type ClientAutomationKnownBotData struct {
	// Detected indicates whether a known bot was detected.
	Detected bool `json:"detected"`

	// ID is the bot identifier. Empty if no bot detected.
	// Example: "googlebot", "bingbot", "chatgpt"
	ID string `json:"id"`

	// Name is the human-readable bot name. Empty if no bot detected.
	// Example: "Googlebot", "Bingbot", "ChatGPT"
	Name string `json:"name"`

	// Type is the bot type classification. Empty if no bot detected.
	Type string `json:"type"`

	// URL is the link to bot documentation. Empty if no bot detected.
	// Example: "https://developers.google.com/search/docs/crawling-indexing/googlebot"
	URL string `json:"url"`
}

// ClientAutomationToolData contains detected automation tool details.
type ClientAutomationToolData struct {
	// Detected indicates whether an automation tool was detected.
	Detected bool `json:"detected"`

	// ID is the automation tool identifier. Empty if no tool detected.
	// Example: "puppeteer", "selenium", "playwright"
	ID string `json:"id"`

	// Name is the human-readable tool name. Empty if no tool detected.
	// Example: "Puppeteer", "Selenium WebDriver", "Playwright"
	Name string `json:"name"`

	// Type is the automation tool type. Empty if no tool detected.
	Type string `json:"type"`
}

// ClientAutomationData contains information about detected automation.
//
// Available when the Bot Detection module is enabled.
// Null when the Bot Detection module is not enabled.
type ClientAutomationData struct {
	// Headless indicates whether the browser was detected as running in headless mode.
	Headless bool `json:"headless"`

	// AutomationTool contains detected automation tool information.
	AutomationTool ClientAutomationToolData `json:"automation_tool"`

	// KnownBot contains detected known bot information.
	KnownBot ClientAutomationKnownBotData `json:"known_bot"`
}

// ClientData contains information about the user agent and device.
type ClientData struct {
	// HeaderUserAgent is the User-Agent HTTP header value.
	// Example: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:146.0) Gecko/20100101 Firefox/146.0"
	HeaderUserAgent string `json:"header_user_agent"`

	// TimeZone contains time zone information.
	//
	// Available when the Browser Identification module is enabled.
	// Null when the Browser Identification module is not enabled.
	TimeZone null.Value[ClientTimeZoneData] `json:"time_zone"`

	// Browser information.
	//
	// Available when the Browser Identification module is enabled.
	// Null when the Browser Identification module is not enabled.
	Browser null.Value[ClientBrowserData] `json:"browser"`

	// BrowserEngine information.
	//
	// Available when the Browser Identification module is enabled.
	// Null when the Browser Identification module is not enabled.
	BrowserEngine null.Value[ClientBrowserEngineData] `json:"browser_engine"`

	// Device information.
	//
	// Available when the Browser Identification module is enabled.
	// Null when the Browser Identification module is not enabled.
	Device null.Value[ClientDeviceData] `json:"device"`

	// OS information.
	//
	// Available when the Browser Identification module is enabled.
	// Null when the Browser Identification module is not enabled.
	OS null.Value[ClientOSData] `json:"os"`

	// TLSSignature contains TLS signatures.
	//
	// Available when the Bot Detection module is enabled.
	// Null when the Bot Detection module is not enabled.
	TLSSignature null.Value[TLSSignatureData] `json:"tls_signature"`

	// Automation contains automation detection data.
	//
	// Available when the Bot Detection module is enabled.
	// Null when the Bot Detection module is not enabled.
	Automation null.Value[ClientAutomationData] `json:"automation"`
}
