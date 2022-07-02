package status

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPeer(t *testing.T) {
	key := "abc"
	status := NewRecorder()
	err := status.AddPeer(key)
	assert.NoError(t, err, "shouldn't return error")

	_, exists := status.peers[key]
	assert.True(t, exists, "value was found")

	err = status.AddPeer(key)

	assert.Error(t, err, "should return error on duplicate")
}

func TestUpdatePeerState(t *testing.T) {
	key := "abc"
	ip := "10.10.10.10"
	status := NewRecorder()
	peerState := PeerState{
		PubKey: key,
	}

	status.peers[key] = peerState

	peerState.IP = ip

	err := status.UpdatePeerState(peerState)
	assert.NoError(t, err, "shouldn't return error")

	state, exists := status.peers[key]
	assert.True(t, exists, "state should be found")
	assert.Equal(t, ip, state.IP, "ip should be equal")
}

func TestRemovePeer(t *testing.T) {
	key := "abc"
	status := NewRecorder()
	peerState := PeerState{
		PubKey: key,
	}

	status.peers[key] = peerState

	err := status.RemovePeer(key)
	assert.NoError(t, err, "shouldn't return error")

	_, exists := status.peers[key]
	assert.False(t, exists, "state value shouldn't be found")

	err = status.RemovePeer("not existing")
	assert.Error(t, err, "should return error when peer doesn't exist")
}

func TestUpdateLocalPeerState(t *testing.T) {
	localPeerState := LocalPeerState{
		IP:              "10.10.10.10",
		PubKey:          "abc",
		KernelInterface: false,
	}
	status := NewRecorder()

	err := status.UpdateLocalPeerState(localPeerState)
	assert.NoError(t, err, "shouldn't return error")

	assert.Equal(t, localPeerState, status.localPeer, "local peer status should be equal")
}

func TestUpdateSignalState(t *testing.T) {
	url := "https://signal"
	var tests = []struct {
		name      string
		connected bool
		want      SignalState
	}{
		{"should mark as connected", true, SignalState{

			URL:       url,
			Connected: true,
		}},
		{"should mark as disconnected", false, SignalState{
			URL:       url,
			Connected: false,
		}},
	}

	status := NewRecorder()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.connected {
				status.MarkSignalConnected(url)
			} else {
				status.MarkSignalDisconnected(url)
			}
			assert.Equal(t, test.want, status.signal, "signal status should be equal")
		})
	}
}

func TestUpdateManagementState(t *testing.T) {
	url := "https://management"
	var tests = []struct {
		name      string
		connected bool
		want      ManagementState
	}{
		{"should mark as connected", true, ManagementState{

			URL:       url,
			Connected: true,
		}},
		{"should mark as disconnected", false, ManagementState{
			URL:       url,
			Connected: false,
		}},
	}

	status := NewRecorder()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.connected {
				status.MarkManagementConnected(url)
			} else {
				status.MarkManagementDisconnected(url)
			}
			assert.Equal(t, test.want, status.management, "signal status should be equal")
		})
	}
}

func TestGetFullStatus(t *testing.T) {
	key1 := "abc"
	key2 := "def"
	managementState := ManagementState{
		URL:       "https://signal",
		Connected: true,
	}
	signalState := SignalState{
		URL:       "https://signal",
		Connected: true,
	}
	peerState1 := PeerState{
		PubKey: key1,
	}

	peerState2 := PeerState{
		PubKey: key2,
	}

	status := NewRecorder()

	status.management = managementState
	status.signal = signalState
	status.peers[key1] = peerState1
	status.peers[key2] = peerState2

	fullStatus := status.GetFullStatus()

	assert.Equal(t, managementState, fullStatus.ManagementState, "management status should be equal")
	assert.Equal(t, signalState, fullStatus.SignalState, "signal status should be equal")
	assert.ElementsMatch(t, []PeerState{peerState1, peerState2}, fullStatus.Peers, "peers states should match")
}
