package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	nbgroup "github.com/netbirdio/netbird/management/server/group"
	nbpeer "github.com/netbirdio/netbird/management/server/peer"
	"github.com/netbirdio/netbird/management/server/posture"
)

func TestPeer_LoginExpired(t *testing.T) {
	tt := []struct {
		name              string
		expirationEnabled bool
		lastLogin         time.Time
		expected          bool
		accountSettings   *Settings
	}{
		{
			name:              "Peer Login Expiration Disabled. Peer Login Should Not Expire",
			expirationEnabled: false,
			lastLogin:         time.Now().UTC().Add(-25 * time.Hour),
			accountSettings: &Settings{
				PeerLoginExpirationEnabled: true,
				PeerLoginExpiration:        time.Hour,
			},
			expected: false,
		},
		{
			name:              "Peer Login Should Expire",
			expirationEnabled: true,
			lastLogin:         time.Now().UTC().Add(-25 * time.Hour),
			accountSettings: &Settings{
				PeerLoginExpirationEnabled: true,
				PeerLoginExpiration:        time.Hour,
			},
			expected: true,
		},
		{
			name:              "Peer Login Should Not Expire",
			expirationEnabled: true,
			lastLogin:         time.Now().UTC(),
			accountSettings: &Settings{
				PeerLoginExpirationEnabled: true,
				PeerLoginExpiration:        time.Hour,
			},
			expected: false,
		},
	}

	for _, c := range tt {
		t.Run(c.name, func(t *testing.T) {
			peer := &nbpeer.Peer{
				LoginExpirationEnabled: c.expirationEnabled,
				LastLogin:              c.lastLogin,
				UserID:                 userID,
			}

			expired, _ := peer.LoginExpired(c.accountSettings.PeerLoginExpiration)
			assert.Equal(t, expired, c.expected)
		})
	}
}

func TestAccountManager_GetNetworkMap(t *testing.T) {
	manager, err := createManager(t)
	if err != nil {
		t.Fatal(err)
		return
	}

	expectedId := "test_account"
	userId := "account_creator"
	account, err := createAccount(manager, expectedId, userId, "")
	if err != nil {
		t.Fatal(err)
	}

	setupKey, err := manager.CreateSetupKey(context.Background(), account.Id, "test-key", SetupKeyReusable, time.Hour, nil, 999, userId, false)
	if err != nil {
		t.Fatal("error creating setup key")
		return
	}

	peerKey1, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}

	peer1, _, _, err := manager.AddPeer(context.Background(), setupKey.Key, "", &nbpeer.Peer{
		Key:  peerKey1.PublicKey().String(),
		Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-1"},
	})
	if err != nil {
		t.Errorf("expecting peer to be added, got failure %v", err)
		return
	}

	peerKey2, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}
	_, _, _, err = manager.AddPeer(context.Background(), setupKey.Key, "", &nbpeer.Peer{
		Key:  peerKey2.PublicKey().String(),
		Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-2"},
	})

	if err != nil {
		t.Errorf("expecting peer to be added, got failure %v", err)
		return
	}

	networkMap, err := manager.GetNetworkMap(context.Background(), peer1.ID)
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(networkMap.Peers) != 1 {
		t.Errorf("expecting Account NetworkMap to have 1 peers, got %v", len(networkMap.Peers))
		return
	}

	if networkMap.Peers[0].Key != peerKey2.PublicKey().String() {
		t.Errorf(
			"expecting Account NetworkMap to have peer with a key %s, got %s",
			peerKey2.PublicKey().String(),
			networkMap.Peers[0].Key,
		)
	}
}

func TestAccountManager_GetNetworkMapWithPolicy(t *testing.T) {
	// TODO: disable until we start use policy again
	t.Skip()
	manager, err := createManager(t)
	if err != nil {
		t.Fatal(err)
		return
	}

	expectedID := "test_account"
	userID := "account_creator"
	account, err := createAccount(manager, expectedID, userID, "")
	if err != nil {
		t.Fatal(err)
	}

	var setupKey *SetupKey
	for _, key := range account.SetupKeys {
		if key.Type == SetupKeyReusable {
			setupKey = key
		}
	}

	peerKey1, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}

	peer1, _, _, err := manager.AddPeer(context.Background(), setupKey.Key, "", &nbpeer.Peer{
		Key:  peerKey1.PublicKey().String(),
		Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-1"},
	})
	if err != nil {
		t.Errorf("expecting peer to be added, got failure %v", err)
		return
	}

	peerKey2, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}
	peer2, _, _, err := manager.AddPeer(context.Background(), setupKey.Key, "", &nbpeer.Peer{
		Key:  peerKey2.PublicKey().String(),
		Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-2"},
	})
	if err != nil {
		t.Errorf("expecting peer to be added, got failure %v", err)
		return
	}

	policies, err := manager.ListPolicies(context.Background(), account.Id, userID)
	if err != nil {
		t.Errorf("expecting to get a list of rules, got failure %v", err)
		return
	}

	err = manager.DeletePolicy(context.Background(), account.Id, policies[0].ID, userID)
	if err != nil {
		t.Errorf("expecting to delete 1 group, got failure %v", err)
		return
	}
	var (
		group1 nbgroup.Group
		group2 nbgroup.Group
		policy Policy
	)

	group1.ID = xid.New().String()
	group2.ID = xid.New().String()
	group1.Name = "src"
	group2.Name = "dst"
	policy.ID = xid.New().String()
	group1.Peers = append(group1.Peers, peer1.ID)
	group2.Peers = append(group2.Peers, peer2.ID)

	err = manager.SaveGroup(context.Background(), account.Id, userID, &group1)
	if err != nil {
		t.Errorf("expecting group1 to be added, got failure %v", err)
		return
	}
	err = manager.SaveGroup(context.Background(), account.Id, userID, &group2)
	if err != nil {
		t.Errorf("expecting group2 to be added, got failure %v", err)
		return
	}

	policy.Name = "test"
	policy.Enabled = true
	policy.Rules = []*PolicyRule{
		{
			Enabled:       true,
			Sources:       []string{group1.ID},
			Destinations:  []string{group2.ID},
			Bidirectional: true,
			Action:        PolicyTrafficActionAccept,
		},
	}
	err = manager.SavePolicy(context.Background(), account.Id, userID, &policy)
	if err != nil {
		t.Errorf("expecting rule to be added, got failure %v", err)
		return
	}

	networkMap1, err := manager.GetNetworkMap(context.Background(), peer1.ID)
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(networkMap1.Peers) != 1 {
		t.Errorf(
			"expecting Account NetworkMap to have 1 peers, got %v: %v",
			len(networkMap1.Peers),
			networkMap1.Peers,
		)
		return
	}

	if networkMap1.Peers[0].Key != peerKey2.PublicKey().String() {
		t.Errorf(
			"expecting Account NetworkMap to have peer with a key %s, got %s",
			peerKey2.PublicKey().String(),
			networkMap1.Peers[0].Key,
		)
	}

	networkMap2, err := manager.GetNetworkMap(context.Background(), peer2.ID)
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(networkMap2.Peers) != 1 {
		t.Errorf("expecting Account NetworkMap to have 1 peers, got %v", len(networkMap2.Peers))
	}

	if len(networkMap2.Peers) > 0 && networkMap2.Peers[0].Key != peerKey1.PublicKey().String() {
		t.Errorf(
			"expecting Account NetworkMap to have peer with a key %s, got %s",
			peerKey1.PublicKey().String(),
			networkMap2.Peers[0].Key,
		)
	}

	policy.Enabled = false
	err = manager.SavePolicy(context.Background(), account.Id, userID, &policy)
	if err != nil {
		t.Errorf("expecting rule to be added, got failure %v", err)
		return
	}

	networkMap1, err = manager.GetNetworkMap(context.Background(), peer1.ID)
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(networkMap1.Peers) != 0 {
		t.Errorf(
			"expecting Account NetworkMap to have 0 peers, got %v: %v",
			len(networkMap1.Peers),
			networkMap1.Peers,
		)
		return
	}

	networkMap2, err = manager.GetNetworkMap(context.Background(), peer2.ID)
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(networkMap2.Peers) != 0 {
		t.Errorf("expecting Account NetworkMap to have 0 peers, got %v", len(networkMap2.Peers))
	}
}

func TestAccountManager_GetPeerNetwork(t *testing.T) {
	manager, err := createManager(t)
	if err != nil {
		t.Fatal(err)
		return
	}

	expectedId := "test_account"
	userId := "account_creator"
	account, err := createAccount(manager, expectedId, userId, "")
	if err != nil {
		t.Fatal(err)
	}

	setupKey, err := manager.CreateSetupKey(context.Background(), account.Id, "test-key", SetupKeyReusable, time.Hour, nil, 999, userId, false)
	if err != nil {
		t.Fatal("error creating setup key")
		return
	}

	peerKey1, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}

	peer1, _, _, err := manager.AddPeer(context.Background(), setupKey.Key, "", &nbpeer.Peer{
		Key:  peerKey1.PublicKey().String(),
		Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-1"},
	})
	if err != nil {
		t.Errorf("expecting peer to be added, got failure %v", err)
		return
	}

	peerKey2, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}
	_, _, _, err = manager.AddPeer(context.Background(), setupKey.Key, "", &nbpeer.Peer{
		Key:  peerKey2.PublicKey().String(),
		Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-2"},
	})

	if err != nil {
		t.Errorf("expecting peer to be added, got failure %v", err)
		return
	}

	network, err := manager.GetPeerNetwork(context.Background(), peer1.ID)
	if err != nil {
		t.Fatal(err)
		return
	}

	if account.Network.Identifier != network.Identifier {
		t.Errorf("expecting Account Networks ID to be equal, got %s expected %s", network.Identifier, account.Network.Identifier)
	}
}

func TestDefaultAccountManager_GetPeer(t *testing.T) {
	manager, err := createManager(t)
	if err != nil {
		t.Fatal(err)
		return
	}

	// account with an admin and a regular user
	accountID := "test_account"
	adminUser := "account_creator"
	someUser := "some_user"
	account := newAccountWithId(context.Background(), accountID, adminUser, "")
	account.Users[someUser] = &User{
		Id:   someUser,
		Role: UserRoleUser,
	}
	account.Settings.RegularUsersViewBlocked = false

	err = manager.Store.SaveAccount(context.Background(), account)
	if err != nil {
		t.Fatal(err)
		return
	}

	// two peers one added by a regular user and one with a setup key
	setupKey, err := manager.CreateSetupKey(context.Background(), account.Id, "test-key", SetupKeyReusable, time.Hour, nil, 999, adminUser, false)
	if err != nil {
		t.Fatal("error creating setup key")
		return
	}

	peerKey1, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}

	peer1, _, _, err := manager.AddPeer(context.Background(), "", someUser, &nbpeer.Peer{
		Key:  peerKey1.PublicKey().String(),
		Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-2"},
	})
	if err != nil {
		t.Errorf("expecting peer to be added, got failure %v", err)
		return
	}

	peerKey2, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}

	// the second peer added with a setup key
	peer2, _, _, err := manager.AddPeer(context.Background(), setupKey.Key, "", &nbpeer.Peer{
		Key:  peerKey2.PublicKey().String(),
		Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-2"},
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	// the user can see its own peer
	peer, err := manager.GetPeer(context.Background(), accountID, peer1.ID, someUser)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.NotNil(t, peer)

	// the user can see peer2 because peer1 of the user has access to peer2 due to the All group and the default rule 0 all-to-all access
	peer, err = manager.GetPeer(context.Background(), accountID, peer2.ID, someUser)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.NotNil(t, peer)

	// delete the all-to-all policy so that user's peer1 has no access to peer2
	for _, policy := range account.Policies {
		err = manager.DeletePolicy(context.Background(), accountID, policy.ID, adminUser)
		if err != nil {
			t.Fatal(err)
			return
		}
	}

	// at this point the user can't see the details of peer2
	peer, err = manager.GetPeer(context.Background(), accountID, peer2.ID, someUser) //nolint
	assert.Error(t, err)

	// admin users can always access all the peers
	peer, err = manager.GetPeer(context.Background(), accountID, peer1.ID, adminUser)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.NotNil(t, peer)

	peer, err = manager.GetPeer(context.Background(), accountID, peer2.ID, adminUser)
	if err != nil {
		t.Fatal(err)
		return
	}
	assert.NotNil(t, peer)
}

func TestDefaultAccountManager_GetPeers(t *testing.T) {
	testCases := []struct {
		name                string
		role                UserRole
		limitedViewSettings bool
		isServiceUser       bool
		expectedPeerCount   int
	}{
		{
			name:                "Regular user, no limited view settings, not a service user",
			role:                UserRoleUser,
			limitedViewSettings: false,
			isServiceUser:       false,
			expectedPeerCount:   1,
		},
		{
			name:                "Service user, no limited view settings",
			role:                UserRoleUser,
			limitedViewSettings: false,
			isServiceUser:       true,
			expectedPeerCount:   2,
		},
		{
			name:                "Regular user, limited view settings",
			role:                UserRoleUser,
			limitedViewSettings: true,
			isServiceUser:       false,
			expectedPeerCount:   0,
		},
		{
			name:                "Service user, limited view settings",
			role:                UserRoleUser,
			limitedViewSettings: true,
			isServiceUser:       true,
			expectedPeerCount:   2,
		},
		{
			name:                "Admin, no limited view settings, not a service user",
			role:                UserRoleAdmin,
			limitedViewSettings: false,
			isServiceUser:       false,
			expectedPeerCount:   2,
		},
		{
			name:                "Admin service user, no limited view settings",
			role:                UserRoleAdmin,
			limitedViewSettings: false,
			isServiceUser:       true,
			expectedPeerCount:   2,
		},
		{
			name:                "Admin, limited view settings",
			role:                UserRoleAdmin,
			limitedViewSettings: true,
			isServiceUser:       false,
			expectedPeerCount:   2,
		},
		{
			name:                "Admin Service user, limited view settings",
			role:                UserRoleAdmin,
			limitedViewSettings: true,
			isServiceUser:       true,
			expectedPeerCount:   2,
		},
		{
			name:                "Owner, no limited view settings",
			role:                UserRoleOwner,
			limitedViewSettings: true,
			isServiceUser:       false,
			expectedPeerCount:   2,
		},
		{
			name:                "Owner, limited view settings",
			role:                UserRoleOwner,
			limitedViewSettings: true,
			isServiceUser:       false,
			expectedPeerCount:   2,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			manager, err := createManager(t)
			if err != nil {
				t.Fatal(err)
				return
			}

			// account with an admin and a regular user
			accountID := "test_account"
			adminUser := "account_creator"
			someUser := "some_user"
			account := newAccountWithId(context.Background(), accountID, adminUser, "")
			account.Users[someUser] = &User{
				Id:            someUser,
				Role:          testCase.role,
				IsServiceUser: testCase.isServiceUser,
			}
			account.Policies = []*Policy{}
			account.Settings.RegularUsersViewBlocked = testCase.limitedViewSettings

			err = manager.Store.SaveAccount(context.Background(), account)
			if err != nil {
				t.Fatal(err)
				return
			}

			peerKey1, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				t.Fatal(err)
				return
			}

			peerKey2, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				t.Fatal(err)
				return
			}

			_, _, _, err = manager.AddPeer(context.Background(), "", someUser, &nbpeer.Peer{
				Key:  peerKey1.PublicKey().String(),
				Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-1"},
			})
			if err != nil {
				t.Errorf("expecting peer to be added, got failure %v", err)
				return
			}

			_, _, _, err = manager.AddPeer(context.Background(), "", adminUser, &nbpeer.Peer{
				Key:  peerKey2.PublicKey().String(),
				Meta: nbpeer.PeerSystemMeta{Hostname: "test-peer-2"},
			})
			if err != nil {
				t.Errorf("expecting peer to be added, got failure %v", err)
				return
			}

			peers, err := manager.GetPeers(context.Background(), accountID, someUser)
			if err != nil {
				t.Fatal(err)
				return
			}
			assert.NotNil(t, peers)

			assert.Len(t, peers, testCase.expectedPeerCount)

		})
	}

}

func setupTestAccountManager(b *testing.B, peers int, groups int) (*DefaultAccountManager, string, string, error) {
	b.Helper()

	manager, err := createManager(b)
	if err != nil {
		return nil, "", "", err
	}

	accountID := "test_account"
	adminUser := "account_creator"
	regularUser := "regular_user"

	account := newAccountWithId(context.Background(), accountID, adminUser, "")
	account.Users[regularUser] = &User{
		Id:   regularUser,
		Role: UserRoleUser,
	}

	// Create peers
	for i := 0; i < peers; i++ {
		peerKey, _ := wgtypes.GeneratePrivateKey()
		peer := &nbpeer.Peer{
			ID:       fmt.Sprintf("peer-%d", i),
			DNSLabel: fmt.Sprintf("peer-%d", i),
			Key:      peerKey.PublicKey().String(),
			IP:       net.ParseIP(fmt.Sprintf("100.64.%d.%d", i/256, i%256)),
			Status:   &nbpeer.PeerStatus{},
			UserID:   regularUser,
		}
		account.Peers[peer.ID] = peer
	}

	// Create groups and policies
	account.Policies = make([]*Policy, 0, groups)
	for i := 0; i < groups; i++ {
		groupID := fmt.Sprintf("group-%d", i)
		group := &nbgroup.Group{
			ID:   groupID,
			Name: fmt.Sprintf("Group %d", i),
		}
		for j := 0; j < peers/groups; j++ {
			peerIndex := i*(peers/groups) + j
			group.Peers = append(group.Peers, fmt.Sprintf("peer-%d", peerIndex))
		}
		account.Groups[groupID] = group

		// Create a policy for this group
		policy := &Policy{
			ID:      fmt.Sprintf("policy-%d", i),
			Name:    fmt.Sprintf("Policy for Group %d", i),
			Enabled: true,
			Rules: []*PolicyRule{
				{
					ID:            fmt.Sprintf("rule-%d", i),
					Name:          fmt.Sprintf("Rule for Group %d", i),
					Enabled:       true,
					Sources:       []string{groupID},
					Destinations:  []string{groupID},
					Bidirectional: true,
					Protocol:      PolicyRuleProtocolALL,
					Action:        PolicyTrafficActionAccept,
				},
			},
		}
		account.Policies = append(account.Policies, policy)
	}

	account.PostureChecks = []*posture.Checks{
		{
			ID:   "PostureChecksAll",
			Name: "All",
			Checks: posture.ChecksDefinition{
				NBVersionCheck: &posture.NBVersionCheck{
					MinVersion: "0.0.1",
				},
			},
		},
	}

	err = manager.Store.SaveAccount(context.Background(), account)
	if err != nil {
		return nil, "", "", err
	}

	return manager, accountID, regularUser, nil
}

func BenchmarkGetPeers(b *testing.B) {
	benchCases := []struct {
		name   string
		peers  int
		groups int
	}{
		{"Small", 50, 5},
		{"Medium", 500, 10},
		{"Large", 5000, 20},
		{"Small single", 50, 1},
		{"Medium single", 500, 1},
		{"Large 5", 5000, 5},
	}

	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)
	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			manager, accountID, userID, err := setupTestAccountManager(b, bc.peers, bc.groups)
			if err != nil {
				b.Fatalf("Failed to setup test account manager: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := manager.GetPeers(context.Background(), accountID, userID)
				if err != nil {
					b.Fatalf("GetPeers failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkUpdateAccountPeers(b *testing.B) {
	benchCases := []struct {
		name   string
		peers  int
		groups int
	}{
		{"Small", 50, 5},
		{"Medium", 500, 10},
		{"Large", 5000, 20},
		{"Small single", 50, 1},
		{"Medium single", 500, 1},
		{"Large 5", 5000, 5},
	}

	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			manager, accountID, _, err := setupTestAccountManager(b, bc.peers, bc.groups)
			if err != nil {
				b.Fatalf("Failed to setup test account manager: %v", err)
			}

			ctx := context.Background()

			account, err := manager.Store.GetAccount(ctx, accountID)
			if err != nil {
				b.Fatalf("Failed to get account: %v", err)
			}

			peerChannels := make(map[string]chan *UpdateMessage)

			for peerID := range account.Peers {
				peerChannels[peerID] = make(chan *UpdateMessage, channelBufferSize)
			}

			manager.peersUpdateManager.peerChannels = peerChannels

			b.ResetTimer()
			start := time.Now()

			for i := 0; i < b.N; i++ {
				manager.updateAccountPeers(ctx, account)
			}

			duration := time.Since(start)
			b.ReportMetric(float64(duration.Nanoseconds())/float64(b.N)/1e6, "ms/op")
			b.ReportMetric(0, "ns/op")
		})
	}
}
