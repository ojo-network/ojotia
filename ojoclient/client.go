package ojoclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetParams returns the current on-chain parameters of the x/oracle module.
func GetPrices(ctx context.Context, grpcEndpoint string) (*oracletypes.QueryExchangeRatesResponse, error) {
	grpcConn, err := grpc.Dial(
		grpcEndpoint,
		// the Cosmos SDK doesn't support any transport security mechanism
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialerFunc),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial Cosmos gRPC service: %w", err)
	}

	defer grpcConn.Close()
	queryClient := oracletypes.NewQueryClient(grpcConn)

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	queryResponse, err := queryClient.ExchangeRates(ctx, &oracletypes.QueryExchangeRates{})
	if err != nil {
		return nil, fmt.Errorf("failed to get x/oracle params: %w", err)
	}

	return queryResponse, nil
}

func EncodePrices(rates oracletypes.QueryExchangeRatesResponse) ([]byte, error) {
	bytes, err := json.Marshal(rates)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func DecodePrices(bytes []byte) (oracletypes.QueryExchangeRatesResponse, error) {
	rates := oracletypes.QueryExchangeRatesResponse{}
	err := json.Unmarshal(bytes, &rates)
	if err != nil {
		return oracletypes.QueryExchangeRatesResponse{}, err
	}
	return rates, nil
}

func ProtocolAndAddress(listenAddr string) (string, string) {
	protocol, address := "tcp", listenAddr

	parts := strings.SplitN(address, "://", 2)
	if len(parts) == 2 {
		protocol, address = parts[0], parts[1]
	}

	return protocol, address
}

// Connect dials the given address and returns a net.Conn. The protoAddr
// argument should be prefixed with the protocol,
// eg. "tcp://127.0.0.1:8080" or "unix:///tmp/test.sock".
func Connect(protoAddr string) (net.Conn, error) {
	proto, address := ProtocolAndAddress(protoAddr)
	conn, err := net.Dial(proto, address)
	return conn, err
}

func dialerFunc(_ context.Context, addr string) (net.Conn, error) {
	return Connect(addr)
}
