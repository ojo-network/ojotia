package tia

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/ojo-network/ojotia/ojoclient"
	openrpc "github.com/rollkit/celestia-openrpc"
	"github.com/rollkit/celestia-openrpc/types/blob"
	"github.com/rollkit/celestia-openrpc/types/share"
)

const (
	namespace = "ojotia"
)

type DAConfig struct {
	Rpc         string `koanf:"rpc"`
	NamespaceId string `koanf:"namespace-id"`
	AuthToken   string `koanf:"auth-token"`
}

type CelestiaDA struct {
	cfg       DAConfig
	client    *openrpc.Client
	namespace share.Namespace
}

func NewCelestiaDA(cfg DAConfig) (*CelestiaDA, error) {
	daClient, err := openrpc.NewClient(context.Background(), cfg.Rpc, cfg.AuthToken)
	if err != nil {
		return nil, err
	}

	if cfg.NamespaceId == "" {
		return nil, errors.New("namespace id cannot be blank")
	}
	nsBytes := []byte(cfg.NamespaceId)

	namespace, err := share.NewBlobNamespaceV0(nsBytes)
	if err != nil {
		return nil, err
	}

	return &CelestiaDA{
		cfg:       cfg,
		client:    daClient,
		namespace: namespace,
	}, nil
}

func (c *CelestiaDA) Store(ctx context.Context, message []byte) ([]byte, uint64, error) {
	dataBlob, err := blob.NewBlobV0(c.namespace, message)
	if err != nil {
		return nil, 0, err
	}
	commitment, err := blob.CreateCommitment(dataBlob)
	if err != nil {
		return nil, 0, err
	}
	height, err := c.client.Blob.Submit(ctx, []*blob.Blob{dataBlob}, openrpc.DefaultSubmitOptions())
	if err != nil {
		return nil, 0, err
	}
	if height == 0 {
		return nil, 0, errors.New("unexpected response code")
	}

	return commitment, height, nil
}

func (c *CelestiaDA) Read(ctx context.Context, commitment string, height uint64) ([]byte, error) {
	fmt.Println("Requesting data from Celestia", "namespace", c.cfg.NamespaceId, "commitment", commitment, "height", height)

	blob, err := c.client.Blob.Get(ctx, height, c.namespace, []byte(commitment))
	if err != nil {
		return nil, err
	}

	fmt.Println("Succesfully fetched data from Celestia", "namespace", c.cfg.NamespaceId, "height", height, "commitment", commitment)

	return blob.Data, nil
}

func Submit(auth, celestiaRPC, ojoRPC string, ctx context.Context) error {
	// Check if filename is provided
	if auth == "" {
		return fmt.Errorf("Please supply auth token")
	}
	if celestiaRPC == "" {
		return fmt.Errorf("Please supply celestia RPC")
	}
	if ojoRPC == "" {
		return fmt.Errorf("Please supply ojo RPC")
	}

	// Start Celestia DA
	daConfig := DAConfig{
		Rpc:         celestiaRPC, //"http://localhost:26658",
		NamespaceId: namespace,
		AuthToken:   auth,
	}

	celestiaDA, err := NewCelestiaDA(daConfig)
	if err != nil {
		return err
	}

	prices, err := ojoclient.GetPrices(ctx, ojoRPC)
	if err != nil {
		return err
	}
	data, err := ojoclient.EncodePrices(*prices)
	if err != nil {
		return err
	}

	commitment, height, err := celestiaDA.Store(context.Background(), data)
	if err != nil {
		return err
	}
	fmt.Println("Succesfully submitted blob to Celestia")
	fmt.Println("Height: ", height)
	fmt.Println("Commitment string: ", hex.EncodeToString(commitment))
	return nil
}

func Query(auth, celestiaRPC, commitment, height string) error {
	if commitment == "" {
		return fmt.Errorf("Please provide commitment")
	}
	if auth == "" {
		return fmt.Errorf("Please supply auth token")
	}
	if celestiaRPC == "" {
		return fmt.Errorf("Please supply celestia RPC")
	}

	heightInt, err := strconv.ParseUint(string(height), 10, 64)
	if err != nil {
		return err
	}
	if heightInt == 0 {
		return fmt.Errorf("Please provide height")
	}

	// Start Celestia DA
	daConfig := DAConfig{
		Rpc:         celestiaRPC,
		NamespaceId: namespace,
		AuthToken:   auth,
	}
	celestiaDA, err := NewCelestiaDA(daConfig)
	if err != nil {
		return err
	}

	commitmentBytes, err := hex.DecodeString(commitment)
	if err != nil {
		return err
	}
	data, err := celestiaDA.Read(context.Background(), string(commitmentBytes), heightInt)
	if err != nil {
		return err
	}

	priceData, err := ojoclient.DecodePrices(data)
	if err != nil {
		return err
	}

	fmt.Println("Celestia queried successfully!")
	fmt.Println("Ojo Price Data:")
	fmt.Println(priceData)
	return nil
}
