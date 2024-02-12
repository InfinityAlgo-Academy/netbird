// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.1-0.20220912230023-4a1477f6a8ba DO NOT EDIT.
package api

import (
	"time"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
	TokenAuthScopes  = "TokenAuth.Scopes"
)

// Defines values for EventActivityCode.
const (
	EventActivityCodeAccountCreate                            EventActivityCode = "account.create"
	EventActivityCodeAccountSettingPeerLoginExpirationDisable EventActivityCode = "account.setting.peer.login.expiration.disable"
	EventActivityCodeAccountSettingPeerLoginExpirationEnable  EventActivityCode = "account.setting.peer.login.expiration.enable"
	EventActivityCodeAccountSettingPeerLoginExpirationUpdate  EventActivityCode = "account.setting.peer.login.expiration.update"
	EventActivityCodeDnsSettingDisabledManagementGroupAdd     EventActivityCode = "dns.setting.disabled.management.group.add"
	EventActivityCodeDnsSettingDisabledManagementGroupDelete  EventActivityCode = "dns.setting.disabled.management.group.delete"
	EventActivityCodeGroupAdd                                 EventActivityCode = "group.add"
	EventActivityCodeGroupUpdate                              EventActivityCode = "group.update"
	EventActivityCodeNameserverGroupAdd                       EventActivityCode = "nameserver.group.add"
	EventActivityCodeNameserverGroupDelete                    EventActivityCode = "nameserver.group.delete"
	EventActivityCodeNameserverGroupUpdate                    EventActivityCode = "nameserver.group.update"
	EventActivityCodePeerLoginExpirationDisable               EventActivityCode = "peer.login.expiration.disable"
	EventActivityCodePeerLoginExpirationEnable                EventActivityCode = "peer.login.expiration.enable"
	EventActivityCodePeerLoginExpire                          EventActivityCode = "peer.login.expire"
	EventActivityCodePeerRename                               EventActivityCode = "peer.rename"
	EventActivityCodePeerSshDisable                           EventActivityCode = "peer.ssh.disable"
	EventActivityCodePeerSshEnable                            EventActivityCode = "peer.ssh.enable"
	EventActivityCodePersonalAccessTokenCreate                EventActivityCode = "personal.access.token.create"
	EventActivityCodePersonalAccessTokenDelete                EventActivityCode = "personal.access.token.delete"
	EventActivityCodePolicyAdd                                EventActivityCode = "policy.add"
	EventActivityCodePolicyDelete                             EventActivityCode = "policy.delete"
	EventActivityCodePolicyUpdate                             EventActivityCode = "policy.update"
	EventActivityCodeRouteAdd                                 EventActivityCode = "route.add"
	EventActivityCodeRouteDelete                              EventActivityCode = "route.delete"
	EventActivityCodeRouteUpdate                              EventActivityCode = "route.update"
	EventActivityCodeRuleAdd                                  EventActivityCode = "rule.add"
	EventActivityCodeRuleDelete                               EventActivityCode = "rule.delete"
	EventActivityCodeRuleUpdate                               EventActivityCode = "rule.update"
	EventActivityCodeServiceUserCreate                        EventActivityCode = "service.user.create"
	EventActivityCodeServiceUserDelete                        EventActivityCode = "service.user.delete"
	EventActivityCodeSetupkeyAdd                              EventActivityCode = "setupkey.add"
	EventActivityCodeSetupkeyGroupAdd                         EventActivityCode = "setupkey.group.add"
	EventActivityCodeSetupkeyGroupDelete                      EventActivityCode = "setupkey.group.delete"
	EventActivityCodeSetupkeyOveruse                          EventActivityCode = "setupkey.overuse"
	EventActivityCodeSetupkeyPeerAdd                          EventActivityCode = "setupkey.peer.add"
	EventActivityCodeSetupkeyRevoke                           EventActivityCode = "setupkey.revoke"
	EventActivityCodeSetupkeyUpdate                           EventActivityCode = "setupkey.update"
	EventActivityCodeUserBlock                                EventActivityCode = "user.block"
	EventActivityCodeUserGroupAdd                             EventActivityCode = "user.group.add"
	EventActivityCodeUserGroupDelete                          EventActivityCode = "user.group.delete"
	EventActivityCodeUserInvite                               EventActivityCode = "user.invite"
	EventActivityCodeUserJoin                                 EventActivityCode = "user.join"
	EventActivityCodeUserPeerAdd                              EventActivityCode = "user.peer.add"
	EventActivityCodeUserPeerDelete                           EventActivityCode = "user.peer.delete"
	EventActivityCodeUserPeerLogin                            EventActivityCode = "user.peer.login"
	EventActivityCodeUserRoleUpdate                           EventActivityCode = "user.role.update"
	EventActivityCodeUserUnblock                              EventActivityCode = "user.unblock"
)

// Defines values for NameserverNsType.
const (
	NameserverNsTypeUdp NameserverNsType = "udp"
)

// Defines values for PolicyRuleAction.
const (
	PolicyRuleActionAccept PolicyRuleAction = "accept"
	PolicyRuleActionDrop   PolicyRuleAction = "drop"
)

// Defines values for PolicyRuleProtocol.
const (
	PolicyRuleProtocolAll  PolicyRuleProtocol = "all"
	PolicyRuleProtocolIcmp PolicyRuleProtocol = "icmp"
	PolicyRuleProtocolTcp  PolicyRuleProtocol = "tcp"
	PolicyRuleProtocolUdp  PolicyRuleProtocol = "udp"
)

// Defines values for PolicyRuleMinimumAction.
const (
	PolicyRuleMinimumActionAccept PolicyRuleMinimumAction = "accept"
	PolicyRuleMinimumActionDrop   PolicyRuleMinimumAction = "drop"
)

// Defines values for PolicyRuleMinimumProtocol.
const (
	PolicyRuleMinimumProtocolAll  PolicyRuleMinimumProtocol = "all"
	PolicyRuleMinimumProtocolIcmp PolicyRuleMinimumProtocol = "icmp"
	PolicyRuleMinimumProtocolTcp  PolicyRuleMinimumProtocol = "tcp"
	PolicyRuleMinimumProtocolUdp  PolicyRuleMinimumProtocol = "udp"
)

// Defines values for PolicyRuleUpdateAction.
const (
	PolicyRuleUpdateActionAccept PolicyRuleUpdateAction = "accept"
	PolicyRuleUpdateActionDrop   PolicyRuleUpdateAction = "drop"
)

// Defines values for PolicyRuleUpdateProtocol.
const (
	PolicyRuleUpdateProtocolAll  PolicyRuleUpdateProtocol = "all"
	PolicyRuleUpdateProtocolIcmp PolicyRuleUpdateProtocol = "icmp"
	PolicyRuleUpdateProtocolTcp  PolicyRuleUpdateProtocol = "tcp"
	PolicyRuleUpdateProtocolUdp  PolicyRuleUpdateProtocol = "udp"
)

// Defines values for UserStatus.
const (
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"
	UserStatusInvited UserStatus = "invited"
)

// AccessiblePeer defines model for AccessiblePeer.
type AccessiblePeer struct {
	// DnsLabel Peer's DNS label is the parsed peer name for domain resolution. It is used to form an FQDN by appending the account's domain to the peer label. e.g. peer-dns-label.netbird.cloud
	DnsLabel string `json:"dns_label"`

	// Id Peer ID
	Id string `json:"id"`

	// Ip Peer's IP address
	Ip string `json:"ip"`

	// Name Peer's hostname
	Name string `json:"name"`

	// UserId User ID of the user that enrolled this peer
	UserId string `json:"user_id"`
}

// Account defines model for Account.
type Account struct {
	// Id Account ID
	Id       string          `json:"id"`
	Settings AccountSettings `json:"settings"`
}

// AccountExtraSettings defines model for AccountExtraSettings.
type AccountExtraSettings struct {
	// PeerApprovalEnabled (Cloud only) Enables or disables peer approval globally. If enabled, all peers added will be in pending state until approved by an admin.
	PeerApprovalEnabled *bool `json:"peer_approval_enabled,omitempty"`
}

// AccountRequest defines model for AccountRequest.
type AccountRequest struct {
	Settings AccountSettings `json:"settings"`
}

// AccountSettings defines model for AccountSettings.
type AccountSettings struct {
	Extra *AccountExtraSettings `json:"extra,omitempty"`

	// GroupsPropagationEnabled Allows propagate the new user auto groups to peers that belongs to the user
	GroupsPropagationEnabled *bool `json:"groups_propagation_enabled,omitempty"`

	// JwtAllowGroups List of groups to which users are allowed access
	JwtAllowGroups *[]string `json:"jwt_allow_groups,omitempty"`

	// JwtGroupsClaimName Name of the claim from which we extract groups names to add it to account groups.
	JwtGroupsClaimName *string `json:"jwt_groups_claim_name,omitempty"`

	// JwtGroupsEnabled Allows extract groups from JWT claim and add it to account groups.
	JwtGroupsEnabled *bool `json:"jwt_groups_enabled,omitempty"`

	// PeerLoginExpiration Period of time after which peer login expires (seconds).
	PeerLoginExpiration int `json:"peer_login_expiration"`

	// PeerLoginExpirationEnabled Enables or disables peer login expiration globally. After peer's login has expired the user has to log in (authenticate). Applies only to peers that were added by a user (interactive SSO login).
	PeerLoginExpirationEnabled bool `json:"peer_login_expiration_enabled"`
}

// AccountUsageStats defines model for AccountUsageStats.
type AccountUsageStats struct {
	// ActivePeers The number of active peers in the account.
	ActivePeers int `json:"active_peers"`

	// ActiveUsers The number of active users in the account.
	ActiveUsers int `json:"active_users"`

	// TotalPeers The total number of peers in the account.
	TotalPeers int `json:"total_peers"`

	// TotalUsers The total number of users in the account.
	TotalUsers int `json:"total_users"`
}

// DNSSettings defines model for DNSSettings.
type DNSSettings struct {
	// DisabledManagementGroups Groups whose DNS management is disabled
	DisabledManagementGroups []string `json:"disabled_management_groups"`
}

// Event defines model for Event.
type Event struct {
	// Activity The activity that occurred during the event
	Activity string `json:"activity"`

	// ActivityCode The string code of the activity that occurred during the event
	ActivityCode EventActivityCode `json:"activity_code"`

	// Id Event unique identifier
	Id string `json:"id"`

	// InitiatorEmail The e-mail address of the initiator of the event. E.g., an e-mail of a user that triggered the event.
	InitiatorEmail string `json:"initiator_email"`

	// InitiatorId The ID of the initiator of the event. E.g., an ID of a user that triggered the event.
	InitiatorId string `json:"initiator_id"`

	// InitiatorName The name of the initiator of the event.
	InitiatorName string `json:"initiator_name"`

	// Meta The metadata of the event
	Meta map[string]string `json:"meta"`

	// TargetId The ID of the target of the event. E.g., an ID of the peer that a user removed.
	TargetId string `json:"target_id"`

	// Timestamp The date and time when the event occurred
	Timestamp time.Time `json:"timestamp"`
}

// EventActivityCode The string code of the activity that occurred during the event
type EventActivityCode string

// Group defines model for Group.
type Group struct {
	// Id Group ID
	Id string `json:"id"`

	// Issued How group was issued by API or from JWT token
	Issued *string `json:"issued,omitempty"`

	// Name Group Name identifier
	Name string `json:"name"`

	// Peers List of peers object
	Peers []PeerMinimum `json:"peers"`

	// PeersCount Count of peers associated to the group
	PeersCount int `json:"peers_count"`
}

// GroupMinimum defines model for GroupMinimum.
type GroupMinimum struct {
	// Id Group ID
	Id string `json:"id"`

	// Issued How group was issued by API or from JWT token
	Issued *string `json:"issued,omitempty"`

	// Name Group Name identifier
	Name string `json:"name"`

	// PeersCount Count of peers associated to the group
	PeersCount int `json:"peers_count"`
}

// GroupRequest defines model for GroupRequest.
type GroupRequest struct {
	// Name Group name identifier
	Name string `json:"name"`

	// Peers List of peers ids
	Peers *[]string `json:"peers,omitempty"`
}

// Nameserver defines model for Nameserver.
type Nameserver struct {
	// Ip Nameserver IP
	Ip string `json:"ip"`

	// NsType Nameserver Type
	NsType NameserverNsType `json:"ns_type"`

	// Port Nameserver Port
	Port int `json:"port"`
}

// NameserverNsType Nameserver Type
type NameserverNsType string

// NameserverGroup defines model for NameserverGroup.
type NameserverGroup struct {
	// Description Description of the nameserver group
	Description string `json:"description"`

	// Domains Match domain list. It should be empty only if primary is true.
	Domains []string `json:"domains"`

	// Enabled Nameserver group status
	Enabled bool `json:"enabled"`

	// Groups Distribution group IDs that defines group of peers that will use this nameserver group
	Groups []string `json:"groups"`

	// Id Nameserver group ID
	Id string `json:"id"`

	// Name Name of nameserver group name
	Name string `json:"name"`

	// Nameservers Nameserver list
	Nameservers []Nameserver `json:"nameservers"`

	// Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty.
	Primary bool `json:"primary"`

	// SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty.
	SearchDomainsEnabled bool `json:"search_domains_enabled"`
}

// NameserverGroupRequest defines model for NameserverGroupRequest.
type NameserverGroupRequest struct {
	// Description Description of the nameserver group
	Description string `json:"description"`

	// Domains Match domain list. It should be empty only if primary is true.
	Domains []string `json:"domains"`

	// Enabled Nameserver group status
	Enabled bool `json:"enabled"`

	// Groups Distribution group IDs that defines group of peers that will use this nameserver group
	Groups []string `json:"groups"`

	// Name Name of nameserver group name
	Name string `json:"name"`

	// Nameservers Nameserver list
	Nameservers []Nameserver `json:"nameservers"`

	// Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty.
	Primary bool `json:"primary"`

	// SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty.
	SearchDomainsEnabled bool `json:"search_domains_enabled"`
}

// Peer defines model for Peer.
type Peer struct {
	// AccessiblePeers List of accessible peers
	AccessiblePeers []AccessiblePeer `json:"accessible_peers"`

	// ApprovalRequired (Cloud only) Indicates whether peer needs approval
	ApprovalRequired *bool `json:"approval_required,omitempty"`

	// Connected Peer to Management connection status
	Connected bool `json:"connected"`

	// DnsLabel Peer's DNS label is the parsed peer name for domain resolution. It is used to form an FQDN by appending the account's domain to the peer label. e.g. peer-dns-label.netbird.cloud
	DnsLabel string `json:"dns_label"`

	// Groups Groups that the peer belongs to
	Groups []GroupMinimum `json:"groups"`

	// Hostname Hostname of the machine
	Hostname string `json:"hostname"`

	// Id Peer ID
	Id string `json:"id"`

	// Ip Peer's IP address
	Ip string `json:"ip"`

	// LastLogin Last time this peer performed log in (authentication). E.g., user authenticated.
	LastLogin time.Time `json:"last_login"`

	// LastSeen Last time peer connected to Netbird's management service
	LastSeen time.Time `json:"last_seen"`

	// LoginExpirationEnabled Indicates whether peer login expiration has been enabled or not
	LoginExpirationEnabled bool `json:"login_expiration_enabled"`

	// LoginExpired Indicates whether peer's login expired or not
	LoginExpired bool `json:"login_expired"`

	// Name Peer's hostname
	Name string `json:"name"`

	// Os Peer's operating system and version
	Os string `json:"os"`

	// SshEnabled Indicates whether SSH server is enabled on this peer
	SshEnabled bool `json:"ssh_enabled"`

	// UiVersion Peer's desktop UI version
	UiVersion *string `json:"ui_version,omitempty"`

	// UserId User ID of the user that enrolled this peer
	UserId *string `json:"user_id,omitempty"`

	// Version Peer's daemon or cli version
	Version string `json:"version"`
}

// PeerBase defines model for PeerBase.
type PeerBase struct {
	// ApprovalRequired (Cloud only) Indicates whether peer needs approval
	ApprovalRequired *bool `json:"approval_required,omitempty"`

	// Connected Peer to Management connection status
	Connected bool `json:"connected"`

	// DnsLabel Peer's DNS label is the parsed peer name for domain resolution. It is used to form an FQDN by appending the account's domain to the peer label. e.g. peer-dns-label.netbird.cloud
	DnsLabel string `json:"dns_label"`

	// Groups Groups that the peer belongs to
	Groups []GroupMinimum `json:"groups"`

	// Hostname Hostname of the machine
	Hostname string `json:"hostname"`

	// Id Peer ID
	Id string `json:"id"`

	// Ip Peer's IP address
	Ip string `json:"ip"`

	// LastLogin Last time this peer performed log in (authentication). E.g., user authenticated.
	LastLogin time.Time `json:"last_login"`

	// LastSeen Last time peer connected to Netbird's management service
	LastSeen time.Time `json:"last_seen"`

	// LoginExpirationEnabled Indicates whether peer login expiration has been enabled or not
	LoginExpirationEnabled bool `json:"login_expiration_enabled"`

	// LoginExpired Indicates whether peer's login expired or not
	LoginExpired bool `json:"login_expired"`

	// Name Peer's hostname
	Name string `json:"name"`

	// Os Peer's operating system and version
	Os string `json:"os"`

	// SshEnabled Indicates whether SSH server is enabled on this peer
	SshEnabled bool `json:"ssh_enabled"`

	// UiVersion Peer's desktop UI version
	UiVersion *string `json:"ui_version,omitempty"`

	// UserId User ID of the user that enrolled this peer
	UserId *string `json:"user_id,omitempty"`

	// Version Peer's daemon or cli version
	Version string `json:"version"`
}

// PeerBatch defines model for PeerBatch.
type PeerBatch struct {
	// AccessiblePeersCount Number of accessible peers
	AccessiblePeersCount int `json:"accessible_peers_count"`

	// ApprovalRequired (Cloud only) Indicates whether peer needs approval
	ApprovalRequired *bool `json:"approval_required,omitempty"`

	// Connected Peer to Management connection status
	Connected bool `json:"connected"`

	// DnsLabel Peer's DNS label is the parsed peer name for domain resolution. It is used to form an FQDN by appending the account's domain to the peer label. e.g. peer-dns-label.netbird.cloud
	DnsLabel string `json:"dns_label"`

	// Groups Groups that the peer belongs to
	Groups []GroupMinimum `json:"groups"`

	// Hostname Hostname of the machine
	Hostname string `json:"hostname"`

	// Id Peer ID
	Id string `json:"id"`

	// Ip Peer's IP address
	Ip string `json:"ip"`

	// LastLogin Last time this peer performed log in (authentication). E.g., user authenticated.
	LastLogin time.Time `json:"last_login"`

	// LastSeen Last time peer connected to Netbird's management service
	LastSeen time.Time `json:"last_seen"`

	// LoginExpirationEnabled Indicates whether peer login expiration has been enabled or not
	LoginExpirationEnabled bool `json:"login_expiration_enabled"`

	// LoginExpired Indicates whether peer's login expired or not
	LoginExpired bool `json:"login_expired"`

	// Name Peer's hostname
	Name string `json:"name"`

	// Os Peer's operating system and version
	Os string `json:"os"`

	// SshEnabled Indicates whether SSH server is enabled on this peer
	SshEnabled bool `json:"ssh_enabled"`

	// UiVersion Peer's desktop UI version
	UiVersion *string `json:"ui_version,omitempty"`

	// UserId User ID of the user that enrolled this peer
	UserId *string `json:"user_id,omitempty"`

	// Version Peer's daemon or cli version
	Version string `json:"version"`
}

// PeerMinimum defines model for PeerMinimum.
type PeerMinimum struct {
	// Id Peer ID
	Id string `json:"id"`

	// Name Peer's hostname
	Name string `json:"name"`
}

// PeerRequest defines model for PeerRequest.
type PeerRequest struct {
	// ApprovalRequired (Cloud only) Indicates whether peer needs approval
	ApprovalRequired       *bool  `json:"approval_required,omitempty"`
	LoginExpirationEnabled bool   `json:"login_expiration_enabled"`
	Name                   string `json:"name"`
	SshEnabled             bool   `json:"ssh_enabled"`
}

// PersonalAccessToken defines model for PersonalAccessToken.
type PersonalAccessToken struct {
	// CreatedAt Date the token was created
	CreatedAt time.Time `json:"created_at"`

	// CreatedBy User ID of the user who created the token
	CreatedBy string `json:"created_by"`

	// ExpirationDate Date the token expires
	ExpirationDate time.Time `json:"expiration_date"`

	// Id ID of a token
	Id string `json:"id"`

	// LastUsed Date the token was last used
	LastUsed *time.Time `json:"last_used,omitempty"`

	// Name Name of the token
	Name string `json:"name"`
}

// PersonalAccessTokenGenerated defines model for PersonalAccessTokenGenerated.
type PersonalAccessTokenGenerated struct {
	PersonalAccessToken PersonalAccessToken `json:"personal_access_token"`

	// PlainToken Plain text representation of the generated token
	PlainToken string `json:"plain_token"`
}

// PersonalAccessTokenRequest defines model for PersonalAccessTokenRequest.
type PersonalAccessTokenRequest struct {
	// ExpiresIn Expiration in days
	ExpiresIn int `json:"expires_in"`

	// Name Name of the token
	Name string `json:"name"`
}

// Policy defines model for Policy.
type Policy struct {
	// Description Policy friendly description
	Description string `json:"description"`

	// Enabled Policy status
	Enabled bool `json:"enabled"`

	// Id Policy ID
	Id *string `json:"id,omitempty"`

	// Name Policy name identifier
	Name string `json:"name"`

	// Rules Policy rule object for policy UI editor
	Rules []PolicyRule `json:"rules"`
}

// PolicyMinimum defines model for PolicyMinimum.
type PolicyMinimum struct {
	// Description Policy friendly description
	Description string `json:"description"`

	// Enabled Policy status
	Enabled bool `json:"enabled"`

	// Id Policy ID
	Id *string `json:"id,omitempty"`

	// Name Policy name identifier
	Name string `json:"name"`
}

// PolicyRule defines model for PolicyRule.
type PolicyRule struct {
	// Action Policy rule accept or drops packets
	Action PolicyRuleAction `json:"action"`

	// Bidirectional Define if the rule is applicable in both directions, sources, and destinations.
	Bidirectional bool `json:"bidirectional"`

	// Description Policy rule friendly description
	Description *string `json:"description,omitempty"`

	// Destinations Policy rule destination group IDs
	Destinations []GroupMinimum `json:"destinations"`

	// Enabled Policy rule status
	Enabled bool `json:"enabled"`

	// Id Policy rule ID
	Id *string `json:"id,omitempty"`

	// Name Policy rule name identifier
	Name string `json:"name"`

	// Ports Policy rule affected ports or it ranges list
	Ports *[]string `json:"ports,omitempty"`

	// Protocol Policy rule type of the traffic
	Protocol PolicyRuleProtocol `json:"protocol"`

	// Sources Policy rule source group IDs
	Sources []GroupMinimum `json:"sources"`
}

// PolicyRuleAction Policy rule accept or drops packets
type PolicyRuleAction string

// PolicyRuleProtocol Policy rule type of the traffic
type PolicyRuleProtocol string

// PolicyRuleMinimum defines model for PolicyRuleMinimum.
type PolicyRuleMinimum struct {
	// Action Policy rule accept or drops packets
	Action PolicyRuleMinimumAction `json:"action"`

	// Bidirectional Define if the rule is applicable in both directions, sources, and destinations.
	Bidirectional bool `json:"bidirectional"`

	// Description Policy rule friendly description
	Description *string `json:"description,omitempty"`

	// Enabled Policy rule status
	Enabled bool `json:"enabled"`

	// Id Policy rule ID
	Id *string `json:"id,omitempty"`

	// Name Policy rule name identifier
	Name string `json:"name"`

	// Ports Policy rule affected ports or it ranges list
	Ports *[]string `json:"ports,omitempty"`

	// Protocol Policy rule type of the traffic
	Protocol PolicyRuleMinimumProtocol `json:"protocol"`
}

// PolicyRuleMinimumAction Policy rule accept or drops packets
type PolicyRuleMinimumAction string

// PolicyRuleMinimumProtocol Policy rule type of the traffic
type PolicyRuleMinimumProtocol string

// PolicyRuleUpdate defines model for PolicyRuleUpdate.
type PolicyRuleUpdate struct {
	// Action Policy rule accept or drops packets
	Action PolicyRuleUpdateAction `json:"action"`

	// Bidirectional Define if the rule is applicable in both directions, sources, and destinations.
	Bidirectional bool `json:"bidirectional"`

	// Description Policy rule friendly description
	Description *string `json:"description,omitempty"`

	// Destinations Policy rule destination group IDs
	Destinations []string `json:"destinations"`

	// Enabled Policy rule status
	Enabled bool `json:"enabled"`

	// Id Policy rule ID
	Id *string `json:"id,omitempty"`

	// Name Policy rule name identifier
	Name string `json:"name"`

	// Ports Policy rule affected ports or it ranges list
	Ports *[]string `json:"ports,omitempty"`

	// Protocol Policy rule type of the traffic
	Protocol PolicyRuleUpdateProtocol `json:"protocol"`

	// Sources Policy rule source group IDs
	Sources []string `json:"sources"`
}

// PolicyRuleUpdateAction Policy rule accept or drops packets
type PolicyRuleUpdateAction string

// PolicyRuleUpdateProtocol Policy rule type of the traffic
type PolicyRuleUpdateProtocol string

// PolicyUpdate defines model for PolicyUpdate.
type PolicyUpdate struct {
	// Description Policy friendly description
	Description string `json:"description"`

	// Enabled Policy status
	Enabled bool `json:"enabled"`

	// Id Policy ID
	Id *string `json:"id,omitempty"`

	// Name Policy name identifier
	Name string `json:"name"`

	// Rules Policy rule object for policy UI editor
	Rules []PolicyRuleUpdate `json:"rules"`
}

// Route defines model for Route.
type Route struct {
	// Description Route description
	Description string `json:"description"`

	// Enabled Route status
	Enabled bool `json:"enabled"`

	// Groups Group IDs containing routing peers
	Groups []string `json:"groups"`

	// Id Route Id
	Id string `json:"id"`

	// Masquerade Indicate if peer should masquerade traffic to this route's prefix
	Masquerade bool `json:"masquerade"`

	// Metric Route metric number. Lowest number has higher priority
	Metric int `json:"metric"`

	// Network Network range in CIDR format
	Network string `json:"network"`

	// NetworkId Route network identifier, to group HA routes
	NetworkId string `json:"network_id"`

	// NetworkType Network type indicating if it is IPv4 or IPv6
	NetworkType string `json:"network_type"`

	// Peer Peer Identifier associated with route. This property can not be set together with `peer_groups`
	Peer *string `json:"peer,omitempty"`

	// PeerGroups Peers Group Identifier associated with route. This property can not be set together with `peer`
	PeerGroups *[]string `json:"peer_groups,omitempty"`
}

// RouteRequest defines model for RouteRequest.
type RouteRequest struct {
	// Description Route description
	Description string `json:"description"`

	// Enabled Route status
	Enabled bool `json:"enabled"`

	// Groups Group IDs containing routing peers
	Groups []string `json:"groups"`

	// Masquerade Indicate if peer should masquerade traffic to this route's prefix
	Masquerade bool `json:"masquerade"`

	// Metric Route metric number. Lowest number has higher priority
	Metric int `json:"metric"`

	// Network Network range in CIDR format
	Network string `json:"network"`

	// NetworkId Route network identifier, to group HA routes
	NetworkId string `json:"network_id"`

	// Peer Peer Identifier associated with route. This property can not be set together with `peer_groups`
	Peer *string `json:"peer,omitempty"`

	// PeerGroups Peers Group Identifier associated with route. This property can not be set together with `peer`
	PeerGroups *[]string `json:"peer_groups,omitempty"`
}

// Rule defines model for Rule.
type Rule struct {
	// Description Rule friendly description
	Description string `json:"description"`

	// Destinations Rule destination group IDs
	Destinations []GroupMinimum `json:"destinations"`

	// Disabled Rules status
	Disabled bool `json:"disabled"`

	// Flow Rule flow, currently, only "bidirect" for bi-directional traffic is accepted
	Flow string `json:"flow"`

	// Id Rule ID
	Id string `json:"id"`

	// Name Rule name identifier
	Name string `json:"name"`

	// Sources Rule source group IDs
	Sources []GroupMinimum `json:"sources"`
}

// RuleMinimum defines model for RuleMinimum.
type RuleMinimum struct {
	// Description Rule friendly description
	Description string `json:"description"`

	// Disabled Rules status
	Disabled bool `json:"disabled"`

	// Flow Rule flow, currently, only "bidirect" for bi-directional traffic is accepted
	Flow string `json:"flow"`

	// Name Rule name identifier
	Name string `json:"name"`
}

// RuleRequest defines model for RuleRequest.
type RuleRequest struct {
	// Description Rule friendly description
	Description string `json:"description"`

	// Destinations List of destination group IDs
	Destinations *[]string `json:"destinations,omitempty"`

	// Disabled Rules status
	Disabled bool `json:"disabled"`

	// Flow Rule flow, currently, only "bidirect" for bi-directional traffic is accepted
	Flow string `json:"flow"`

	// Name Rule name identifier
	Name string `json:"name"`

	// Sources List of source group IDs
	Sources *[]string `json:"sources,omitempty"`
}

// SetupKey defines model for SetupKey.
type SetupKey struct {
	// AutoGroups List of group IDs to auto-assign to peers registered with this key
	AutoGroups []string `json:"auto_groups"`

	// Ephemeral Indicate that the peer will be ephemeral or not
	Ephemeral bool `json:"ephemeral"`

	// Expires Setup Key expiration date
	Expires time.Time `json:"expires"`

	// Id Setup Key ID
	Id string `json:"id"`

	// Key Setup Key value
	Key string `json:"key"`

	// LastUsed Setup key last usage date
	LastUsed time.Time `json:"last_used"`

	// Name Setup key name identifier
	Name string `json:"name"`

	// Revoked Setup key revocation status
	Revoked bool `json:"revoked"`

	// State Setup key status, "valid", "overused","expired" or "revoked"
	State string `json:"state"`

	// Type Setup key type, one-off for single time usage and reusable
	Type string `json:"type"`

	// UpdatedAt Setup key last update date
	UpdatedAt time.Time `json:"updated_at"`

	// UsageLimit A number of times this key can be used. The value of 0 indicates the unlimited usage.
	UsageLimit int `json:"usage_limit"`

	// UsedTimes Usage count of setup key
	UsedTimes int `json:"used_times"`

	// Valid Setup key validity status
	Valid bool `json:"valid"`
}

// SetupKeyRequest defines model for SetupKeyRequest.
type SetupKeyRequest struct {
	// AutoGroups List of group IDs to auto-assign to peers registered with this key
	AutoGroups []string `json:"auto_groups"`

	// Ephemeral Indicate that the peer will be ephemeral or not
	Ephemeral *bool `json:"ephemeral,omitempty"`

	// ExpiresIn Expiration time in seconds
	ExpiresIn int `json:"expires_in"`

	// Name Setup Key name
	Name string `json:"name"`

	// Revoked Setup key revocation status
	Revoked bool `json:"revoked"`

	// Type Setup key type, one-off for single time usage and reusable
	Type string `json:"type"`

	// UsageLimit A number of times this key can be used. The value of 0 indicates the unlimited usage.
	UsageLimit int `json:"usage_limit"`
}

// User defines model for User.
type User struct {
	// AutoGroups Group IDs to auto-assign to peers registered by this user
	AutoGroups []string `json:"auto_groups"`

	// Email User's email address
	Email string `json:"email"`

	// Id User ID
	Id string `json:"id"`

	// IsBlocked Is true if this user is blocked. Blocked users can't use the system
	IsBlocked bool `json:"is_blocked"`

	// IsCurrent Is true if authenticated user is the same as this user
	IsCurrent *bool `json:"is_current,omitempty"`

	// IsServiceUser Is true if this user is a service user
	IsServiceUser *bool `json:"is_service_user,omitempty"`

	// Issued How user was issued by API or Integration
	Issued *string `json:"issued,omitempty"`

	// LastLogin Last time this user performed a login to the dashboard
	LastLogin *time.Time `json:"last_login,omitempty"`

	// Name User's name from idp provider
	Name string `json:"name"`

	// Role User's NetBird account role
	Role string `json:"role"`

	// Status User's status
	Status UserStatus `json:"status"`
}

// UserStatus User's status
type UserStatus string

// UserCreateRequest defines model for UserCreateRequest.
type UserCreateRequest struct {
	// AutoGroups Group IDs to auto-assign to peers registered by this user
	AutoGroups []string `json:"auto_groups"`

	// Email User's Email to send invite to
	Email *string `json:"email,omitempty"`

	// IsServiceUser Is true if this user is a service user
	IsServiceUser bool `json:"is_service_user"`

	// Name User's full name
	Name *string `json:"name,omitempty"`

	// Role User's NetBird account role
	Role string `json:"role"`
}

// UserRequest defines model for UserRequest.
type UserRequest struct {
	// AutoGroups Group IDs to auto-assign to peers registered by this user
	AutoGroups []string `json:"auto_groups"`

	// IsBlocked If set to true then user is blocked and can't use the system
	IsBlocked bool `json:"is_blocked"`

	// Role User's NetBird account role
	Role string `json:"role"`
}

// GetApiUsersParams defines parameters for GetApiUsers.
type GetApiUsersParams struct {
	// ServiceUser Filters users and returns either regular users or service users
	ServiceUser *bool `form:"service_user,omitempty" json:"service_user,omitempty"`
}

// PutApiAccountsAccountIdJSONRequestBody defines body for PutApiAccountsAccountId for application/json ContentType.
type PutApiAccountsAccountIdJSONRequestBody = AccountRequest

// PostApiDnsNameserversJSONRequestBody defines body for PostApiDnsNameservers for application/json ContentType.
type PostApiDnsNameserversJSONRequestBody = NameserverGroupRequest

// PutApiDnsNameserversNsgroupIdJSONRequestBody defines body for PutApiDnsNameserversNsgroupId for application/json ContentType.
type PutApiDnsNameserversNsgroupIdJSONRequestBody = NameserverGroupRequest

// PutApiDnsSettingsJSONRequestBody defines body for PutApiDnsSettings for application/json ContentType.
type PutApiDnsSettingsJSONRequestBody = DNSSettings

// PostApiGroupsJSONRequestBody defines body for PostApiGroups for application/json ContentType.
type PostApiGroupsJSONRequestBody = GroupRequest

// PutApiGroupsGroupIdJSONRequestBody defines body for PutApiGroupsGroupId for application/json ContentType.
type PutApiGroupsGroupIdJSONRequestBody = GroupRequest

// PutApiPeersPeerIdJSONRequestBody defines body for PutApiPeersPeerId for application/json ContentType.
type PutApiPeersPeerIdJSONRequestBody = PeerRequest

// PostApiPoliciesJSONRequestBody defines body for PostApiPolicies for application/json ContentType.
type PostApiPoliciesJSONRequestBody = PolicyUpdate

// PutApiPoliciesPolicyIdJSONRequestBody defines body for PutApiPoliciesPolicyId for application/json ContentType.
type PutApiPoliciesPolicyIdJSONRequestBody = PolicyUpdate

// PostApiRoutesJSONRequestBody defines body for PostApiRoutes for application/json ContentType.
type PostApiRoutesJSONRequestBody = RouteRequest

// PutApiRoutesRouteIdJSONRequestBody defines body for PutApiRoutesRouteId for application/json ContentType.
type PutApiRoutesRouteIdJSONRequestBody = RouteRequest

// PostApiRulesJSONRequestBody defines body for PostApiRules for application/json ContentType.
type PostApiRulesJSONRequestBody = RuleRequest

// PutApiRulesRuleIdJSONRequestBody defines body for PutApiRulesRuleId for application/json ContentType.
type PutApiRulesRuleIdJSONRequestBody = RuleRequest

// PostApiSetupKeysJSONRequestBody defines body for PostApiSetupKeys for application/json ContentType.
type PostApiSetupKeysJSONRequestBody = SetupKeyRequest

// PutApiSetupKeysKeyIdJSONRequestBody defines body for PutApiSetupKeysKeyId for application/json ContentType.
type PutApiSetupKeysKeyIdJSONRequestBody = SetupKeyRequest

// PostApiUsersJSONRequestBody defines body for PostApiUsers for application/json ContentType.
type PostApiUsersJSONRequestBody = UserCreateRequest

// PutApiUsersUserIdJSONRequestBody defines body for PutApiUsersUserId for application/json ContentType.
type PutApiUsersUserIdJSONRequestBody = UserRequest

// PostApiUsersUserIdTokensJSONRequestBody defines body for PostApiUsersUserIdTokens for application/json ContentType.
type PostApiUsersUserIdTokensJSONRequestBody = PersonalAccessTokenRequest
