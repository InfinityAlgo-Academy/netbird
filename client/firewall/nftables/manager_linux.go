package nftables

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	fw "github.com/netbirdio/netbird/client/firewall"
	"github.com/netbirdio/netbird/iface"
)

const (
	// FilterTableName is the name of the table that is used for filtering by the Netbird client
	FilterTableName = "netbird-acl"

	// FilterInputChainName is the name of the chain that is used for filtering incoming packets
	FilterInputChainName = "netbird-acl-input-filter"

	// FilterOutputChainName is the name of the chain that is used for filtering outgoing packets
	FilterOutputChainName = "netbird-acl-output-filter"
)

var anyIP = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// Manager of iptables firewall
type Manager struct {
	mutex sync.Mutex

	conn      *nftables.Conn
	tableIPv4 *nftables.Table
	tableIPv6 *nftables.Table

	filterInputChainIPv4  *nftables.Chain
	filterOutputChainIPv4 *nftables.Chain

	filterInputChainIPv6  *nftables.Chain
	filterOutputChainIPv6 *nftables.Chain

	rulesetManager *rulesetManager
	setRemovedIPs  map[string]struct{}

	wgIface iFaceMapper
}

// iFaceMapper defines subset methods of interface required for manager
type iFaceMapper interface {
	Name() string
	Address() iface.WGAddress
}

// Create nftables firewall manager
func Create(wgIface iFaceMapper) (*Manager, error) {
	m := &Manager{
		conn: &nftables.Conn{},

		rulesetManager: newRuleManager(),
		setRemovedIPs:  map[string]struct{}{},

		wgIface: wgIface,
	}

	if err := m.Reset(); err != nil {
		return nil, err
	}

	return m, nil
}

// AddFiltering rule to the firewall
//
// If comment argument is empty firewall manager should set
// rule ID as comment for the rule
func (m *Manager) AddFiltering(
	ip net.IP,
	proto fw.Protocol,
	sPort *fw.Port,
	dPort *fw.Port,
	direction fw.RuleDirection,
	action fw.Action,
	ipsetName string,
	comment string,
) (fw.Rule, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var (
		err   error
		ipset *nftables.Set
		table *nftables.Table
		chain *nftables.Chain
	)

	if direction == fw.RuleDirectionOUT {
		table, chain, err = m.chain(
			ip,
			FilterOutputChainName,
			nftables.ChainHookOutput,
			nftables.ChainPriorityFilter,
			nftables.ChainTypeFilter)
	} else {
		table, chain, err = m.chain(
			ip,
			FilterInputChainName,
			nftables.ChainHookInput,
			nftables.ChainPriorityFilter,
			nftables.ChainTypeFilter)
	}
	if err != nil {
		return nil, err
	}

	rawIP := ip.To4()
	if rawIP == nil {
		rawIP = ip.To16()
	}

	if ipsetName != "" {
		// if we already have set with given name, just add ip to the set
		// and return rule with new ID in other case let's create rule
		// with fresh created set and set element

		var isSetNew bool
		ipset, isSetNew, err = m.getOrCreateSet(table, rawIP, ipsetName)
		if err != nil {
			return nil, fmt.Errorf("get set name: %v", err)
		}

		if err := m.conn.SetAddElements(ipset, []nftables.SetElement{{Key: rawIP}}); err != nil {
			return nil, fmt.Errorf("add set element for the first time: %v", err)
		}

		if !isSetNew {
			// if we already have nftables rules with set
			// just add new rule to the ruleset and return new fw.Rule object

			// check ruleset exists for that ipset
			ruleset, ok := m.rulesetManager.getRulesetBySetID(ipset.ID)
			if !ok {
				return nil, fmt.Errorf("ipset exists in nftables but not in the manager, ipset is not synced")
			}

			return m.rulesetManager.addRule(ruleset, rawIP)
		}
	}

	ifaceKey := expr.MetaKeyIIFNAME
	if direction == fw.RuleDirectionOUT {
		ifaceKey = expr.MetaKeyOIFNAME
	}
	expressions := []expr.Any{
		&expr.Meta{Key: ifaceKey, Register: 1},
		&expr.Cmp{
			Op:       expr.CmpOpEq,
			Register: 1,
			Data:     ifname(m.wgIface.Name()),
		},
	}

	if proto != "all" {
		expressions = append(expressions, &expr.Payload{
			DestRegister: 1,
			Base:         expr.PayloadBaseNetworkHeader,
			Offset:       uint32(9),
			Len:          uint32(1),
		})

		var protoData []byte
		switch proto {
		case fw.ProtocolTCP:
			protoData = []byte{unix.IPPROTO_TCP}
		case fw.ProtocolUDP:
			protoData = []byte{unix.IPPROTO_UDP}
		case fw.ProtocolICMP:
			protoData = []byte{unix.IPPROTO_ICMP}
		default:
			return nil, fmt.Errorf("unsupported protocol: %s", proto)
		}
		expressions = append(expressions, &expr.Cmp{
			Register: 1,
			Op:       expr.CmpOpEq,
			Data:     protoData,
		})
	}

	// check if rawIP contains zeroed IPv4 0.0.0.0 or same IVv6 value
	// in that case not add IP match expression into the rule definition
	if !bytes.HasPrefix(anyIP, rawIP) {
		// source address position
		adrLen := uint32(len(rawIP))
		adrOffset := uint32(12)
		if adrLen == 16 {
			adrOffset = 8
		}

		// change to destination address position if need
		if direction == fw.RuleDirectionOUT {
			adrOffset += adrLen
		}

		expressions = append(expressions,
			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseNetworkHeader,
				Offset:       adrOffset,
				Len:          adrLen,
			},
		)
		// add individual IP for match if no ipset defined
		if ipset == nil {
			expressions = append(expressions,
				&expr.Cmp{
					Op:       expr.CmpOpEq,
					Register: 1,
					Data:     rawIP,
				},
			)
		} else {
			expressions = append(expressions,
				&expr.Lookup{
					SourceRegister: 1,
					SetName:        ipsetName,
					SetID:          ipset.ID,
				},
			)
		}
	}

	if sPort != nil && len(sPort.Values) != 0 {
		expressions = append(expressions,
			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseTransportHeader,
				Offset:       0,
				Len:          2,
			},
			&expr.Cmp{
				Op:       expr.CmpOpEq,
				Register: 1,
				Data:     encodePort(*sPort),
			},
		)
	}

	if dPort != nil && len(dPort.Values) != 0 {
		expressions = append(expressions,
			&expr.Payload{
				DestRegister: 1,
				Base:         expr.PayloadBaseTransportHeader,
				Offset:       2,
				Len:          2,
			},
			&expr.Cmp{
				Op:       expr.CmpOpEq,
				Register: 1,
				Data:     encodePort(*dPort),
			},
		)
	}

	if action == fw.ActionAccept {
		expressions = append(expressions, &expr.Verdict{Kind: expr.VerdictAccept})
	} else {
		expressions = append(expressions, &expr.Verdict{Kind: expr.VerdictDrop})
	}

	rulesetID := uuid.New().String()
	userData := []byte(strings.Join([]string{rulesetID, comment}, " "))

	rule := m.conn.InsertRule(&nftables.Rule{
		Table:    table,
		Chain:    chain,
		Position: 0,
		Exprs:    expressions,
		UserData: userData,
	})

	ruleset := m.rulesetManager.createRuleset(rulesetID, rule, ipset)
	return m.rulesetManager.addRule(ruleset, rawIP)
}

// getOrCreateSet in given table by name
//
// It tries to get set by name if fails creates new one by this name.
// If new set need to be created it calls firewall flush method.
// Second returned argument is a flag that indicates is set just created or not.
func (m *Manager) getOrCreateSet(
	table *nftables.Table,
	rawIP []byte,
	name string,
) (*nftables.Set, bool, error) {
	ipset, err := m.conn.GetSetByName(table, name)
	if err != nil {
		keyType := nftables.TypeIPAddr
		if len(rawIP) == 16 {
			keyType = nftables.TypeIP6Addr
		}
		// else we create new ipset and continue creating rule
		ipset = &nftables.Set{
			Name:    name,
			Table:   table,
			Dynamic: true,
			KeyType: keyType,
		}

		if err := m.conn.AddSet(ipset, nil); err != nil {
			return nil, false, fmt.Errorf("create set: %v", err)
		}

		if err := m.conn.Flush(); err != nil {
			return nil, false, fmt.Errorf("flush created set: %v", err)
		}

		ipset, err = m.conn.GetSetByName(table, name)
		if err != nil {
			return nil, false, fmt.Errorf("get created set: %v", err)
		}
		return ipset, true, nil
	}
	return ipset, false, nil
}

// chain returns the chain for the given IP address with specific settings
func (m *Manager) chain(
	ip net.IP,
	name string,
	hook nftables.ChainHook,
	priority nftables.ChainPriority,
	cType nftables.ChainType,
) (*nftables.Table, *nftables.Chain, error) {
	var err error

	getChain := func(c *nftables.Chain, tf nftables.TableFamily) (*nftables.Chain, error) {
		if c != nil {
			return c, nil
		}
		return m.createChainIfNotExists(tf, name, hook, priority, cType)
	}

	if ip.To4() != nil {
		if name == FilterInputChainName {
			m.filterInputChainIPv4, err = getChain(m.filterInputChainIPv4, nftables.TableFamilyIPv4)
			return m.tableIPv4, m.filterInputChainIPv4, err
		}
		m.filterOutputChainIPv4, err = getChain(m.filterOutputChainIPv4, nftables.TableFamilyIPv4)
		return m.tableIPv4, m.filterOutputChainIPv4, err
	}
	if name == FilterInputChainName {
		m.filterInputChainIPv6, err = getChain(m.filterInputChainIPv6, nftables.TableFamilyIPv6)
		return m.tableIPv4, m.filterInputChainIPv6, err
	}
	m.filterOutputChainIPv6, err = getChain(m.filterOutputChainIPv6, nftables.TableFamilyIPv6)
	return m.tableIPv4, m.filterOutputChainIPv6, err
}

// table returns the table for the given family of the IP address
func (m *Manager) table(family nftables.TableFamily) (*nftables.Table, error) {
	if family == nftables.TableFamilyIPv4 {
		if m.tableIPv4 != nil {
			return m.tableIPv4, nil
		}

		table, err := m.createTableIfNotExists(nftables.TableFamilyIPv4)
		if err != nil {
			return nil, err
		}
		m.tableIPv4 = table
		return m.tableIPv4, nil
	}

	if m.tableIPv6 != nil {
		return m.tableIPv6, nil
	}

	table, err := m.createTableIfNotExists(nftables.TableFamilyIPv6)
	if err != nil {
		return nil, err
	}
	m.tableIPv6 = table
	return m.tableIPv6, nil
}

func (m *Manager) createTableIfNotExists(family nftables.TableFamily) (*nftables.Table, error) {
	tables, err := m.conn.ListTablesOfFamily(family)
	if err != nil {
		return nil, fmt.Errorf("list of tables: %w", err)
	}

	for _, t := range tables {
		if t.Name == FilterTableName {
			return t, nil
		}
	}

	return m.conn.AddTable(&nftables.Table{Name: FilterTableName, Family: nftables.TableFamilyIPv4}), nil
}

func (m *Manager) createChainIfNotExists(
	family nftables.TableFamily,
	name string,
	hooknum nftables.ChainHook,
	priority nftables.ChainPriority,
	chainType nftables.ChainType,
) (*nftables.Chain, error) {
	table, err := m.table(family)
	if err != nil {
		return nil, err
	}

	chains, err := m.conn.ListChainsOfTableFamily(family)
	if err != nil {
		return nil, fmt.Errorf("list of chains: %w", err)
	}

	for _, c := range chains {
		if c.Name == name && c.Table.Name == table.Name {
			return c, nil
		}
	}

	polAccept := nftables.ChainPolicyAccept
	chain := &nftables.Chain{
		Name:     name,
		Table:    table,
		Hooknum:  hooknum,
		Priority: priority,
		Type:     chainType,
		Policy:   &polAccept,
	}

	chain = m.conn.AddChain(chain)

	ifaceKey := expr.MetaKeyIIFNAME
	shiftDSTAddr := 0
	if name == FilterOutputChainName {
		ifaceKey = expr.MetaKeyOIFNAME
		shiftDSTAddr = 1
	}

	expressions := []expr.Any{
		&expr.Meta{Key: ifaceKey, Register: 1},
		&expr.Cmp{
			Op:       expr.CmpOpEq,
			Register: 1,
			Data:     ifname(m.wgIface.Name()),
		},
	}

	mask, _ := netip.AddrFromSlice(m.wgIface.Address().Network.Mask)
	if m.wgIface.Address().IP.To4() == nil {
		ip, _ := netip.AddrFromSlice(m.wgIface.Address().Network.IP.To16())
		expressions = append(expressions,
			&expr.Payload{
				DestRegister: 2,
				Base:         expr.PayloadBaseNetworkHeader,
				Offset:       uint32(8 + (16 * shiftDSTAddr)),
				Len:          16,
			},
			&expr.Bitwise{
				SourceRegister: 2,
				DestRegister:   2,
				Len:            16,
				Xor:            []byte{0x0, 0x0, 0x0, 0x0},
				Mask:           mask.Unmap().AsSlice(),
			},
			&expr.Cmp{
				Op:       expr.CmpOpNeq,
				Register: 2,
				Data:     ip.Unmap().AsSlice(),
			},
			&expr.Verdict{Kind: expr.VerdictAccept},
		)
	} else {
		ip, _ := netip.AddrFromSlice(m.wgIface.Address().Network.IP.To4())
		expressions = append(expressions,
			&expr.Payload{
				DestRegister: 2,
				Base:         expr.PayloadBaseNetworkHeader,
				Offset:       uint32(12 + (4 * shiftDSTAddr)),
				Len:          4,
			},
			&expr.Bitwise{
				SourceRegister: 2,
				DestRegister:   2,
				Len:            4,
				Xor:            []byte{0x0, 0x0, 0x0, 0x0},
				Mask:           m.wgIface.Address().Network.Mask,
			},
			&expr.Cmp{
				Op:       expr.CmpOpNeq,
				Register: 2,
				Data:     ip.Unmap().AsSlice(),
			},
			&expr.Verdict{Kind: expr.VerdictAccept},
		)
	}

	_ = m.conn.AddRule(&nftables.Rule{
		Table: table,
		Chain: chain,
		Exprs: expressions,
	})

	expressions = []expr.Any{
		&expr.Meta{Key: ifaceKey, Register: 1},
		&expr.Cmp{
			Op:       expr.CmpOpEq,
			Register: 1,
			Data:     ifname(m.wgIface.Name()),
		},
		&expr.Verdict{Kind: expr.VerdictDrop},
	}
	_ = m.conn.AddRule(&nftables.Rule{
		Table: table,
		Chain: chain,
		Exprs: expressions,
	})

	if err := m.conn.Flush(); err != nil {
		return nil, err
	}

	return chain, nil
}

// DeleteRule from the firewall by rule definition
func (m *Manager) DeleteRule(rule fw.Rule) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	nativeRule, ok := rule.(*Rule)
	if !ok {
		return fmt.Errorf("invalid rule type")
	}

	if nativeRule.nftRule == nil {
		return nil
	}

	if m.rulesetManager.deleteRule(nativeRule) {
		// deleteRule indicates that we still have IP in the ruleset
		// it means we should not remove the nftables rule but need to update set
		// so we prepare IP to be removed from set on the next flush call
		if nativeRule.nftSet != nil {
			// call twice of delete set element raises error
			// so we need to check if element is already removed
			key := fmt.Sprintf("%s:%v", nativeRule.nftSet.Name, nativeRule.ip)
			if _, ok := m.setRemovedIPs[key]; !ok {
				err := m.conn.SetDeleteElements(nativeRule.nftSet, []nftables.SetElement{{Key: nativeRule.ip}})
				if err != nil {
					log.Errorf("delete elements for set %q: %v", nativeRule.nftSet.Name, err)
				}
				m.setRemovedIPs[key] = struct{}{}
			}
		}
		return nil
	}

	// ruleset doesn't contain IP anymore (or contains only one), remove nft rule
	if err := m.conn.DelRule(nativeRule.nftRule); err != nil {
		log.Errorf("failed to delete rule: %v", err)
	}
	nativeRule.nftRule = nil

	if nativeRule.nftSet != nil {
		m.conn.DelSet(nativeRule.nftSet)
		nativeRule.nftSet = nil
	}
	return nil
}

// Reset firewall to the default state
func (m *Manager) Reset() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	chains, err := m.conn.ListChains()
	if err != nil {
		return fmt.Errorf("list of chains: %w", err)
	}
	for _, c := range chains {
		if c.Name == FilterInputChainName || c.Name == FilterOutputChainName {
			m.conn.DelChain(c)
		}
	}

	tables, err := m.conn.ListTables()
	if err != nil {
		return fmt.Errorf("list of tables: %w", err)
	}
	for _, t := range tables {
		if t.Name == FilterTableName {
			m.conn.DelTable(t)
		}
	}

	return m.conn.Flush()
}

// Flush doesn't need to be implemented for this manager
func (m *Manager) Flush() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.conn.Flush(); err != nil {
		return err
	}

	m.setRemovedIPs = map[string]struct{}{}

	if err := m.refreshRuleHandles(m.tableIPv4, m.filterInputChainIPv4); err != nil {
		log.Errorf("failed to refresh rule handles ipv4 input chain: %v", err)
	}

	if err := m.refreshRuleHandles(m.tableIPv4, m.filterOutputChainIPv4); err != nil {
		log.Errorf("failed to refresh rule handles IPv4 output chain: %v", err)
	}

	if err := m.refreshRuleHandles(m.tableIPv6, m.filterInputChainIPv6); err != nil {
		log.Errorf("failed to refresh rule handles IPv6 input chain: %v", err)
	}

	if err := m.refreshRuleHandles(m.tableIPv6, m.filterOutputChainIPv6); err != nil {
		log.Errorf("failed to refresh rule handles IPv6 output chain: %v", err)
	}

	return nil
}

func (m *Manager) refreshRuleHandles(table *nftables.Table, chain *nftables.Chain) error {
	if table == nil || chain == nil {
		return nil
	}

	list, err := m.conn.GetRules(table, chain)
	if err != nil {
		return err
	}

	for _, rule := range list {
		if len(rule.UserData) != 0 {
			if err := m.rulesetManager.setNftRuleHandle(rule); err != nil {
				log.Errorf("failed to set rule handle: %v", err)
			}
		}
	}

	return nil
}

func encodePort(port fw.Port) []byte {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(port.Values[0]))
	return bs
}

func ifname(n string) []byte {
	b := make([]byte, 16)
	copy(b, []byte(n+"\x00"))
	return b
}
