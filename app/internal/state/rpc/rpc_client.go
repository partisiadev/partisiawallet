package rpc

//func NewChainRPCsFromChains(chains []Chain) []RpcClients {
//	rpcClientsArr := make([]RpcClients, 0)
//	for _, chain := range chains {
//		rpcClients := GetRPCChainClients(chain)
//		rpcClientsArr = append(rpcClientsArr, rpcClients)
//	}
//	sort.Slice(rpcClientsArr, func(i, j int) bool {
//		return strings.ToLower(rpcClientsArr[i].Chain.Name) < strings.ToLower(rpcClientsArr[j].Chain.Name)
//	})
//	return rpcClientsArr
//}

//func (c RpcClient) Connect() (*ethclient.Client, error) {
//	var err error
//	url, err := c.RPC.GetURL(c.ApiKey)
//	if err != nil {
//		return nil, err
//	}
//	conn, err := ethclient.Dial(url)
//	if err != nil {
//		if conn != nil {
//			conn.Close()
//		}
//		return nil, err
//	}
//	return conn, err
//}

//func (c *RpcClient) ShowBalance(acc wallet.Account) (string, error) {
//	if c.GetClient() == nil {
//		return "", ErrNotConnected
//	}
//	Client := c.GetClient()
//	bal, err := Client.BalanceAt(context.Background(), common2.HexToAddress(acc.EthAddress), nil)
//	if err != nil {
//		return "", err
//	}
//	fBalance := new(big.Float)
//	fBalance.SetString(bal.String())
//	val := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
//	return val.String(), nil
//}

// RpcClients wraps multiple RpcClient for a Chain
//type RpcClients struct {
//	Chain   Chain
//	Clients []RpcClient
//}

//type RpcClient struct {
//	ApiKey string
//	RPC    RPC
//}

//	func GetChainRPCClientsArr(chainsSlice []Chain) []*RpcClients {
//		rpcClients := make([]*RpcClients, len(chainsSlice))
//		for i, ch := range chainsSlice {
//			rpcClients[i] = GetRPCChainClients(ch)
//		}
//		return rpcClients
//	}
//func GetRPCChainClients(chain Chain) RpcClients {
//	rpcClients := RpcClients{
//		Chain:   chain,
//		Clients: make([]RpcClient, 0),
//	}
//	for _, rpc := range chain.RPC {
//		rpcClients.Clients = append(rpcClients.Clients, RpcClient{RPC: rpc})
//	}
//	return rpcClients
//}
