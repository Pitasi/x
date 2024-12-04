package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"golang.org/x/sync/errgroup"
)

func TestLocalnet(t *testing.T) {
	ctx := context.Background()

	fromDockerfile := testcontainers.FromDockerfile{
		Context:   "../",
		KeepImage: true,
		Repo:      "meteorkv",
		Tag:       "latest",
	}

	N := 3

	// start N validator containers, run "init" to generate their keys,
	// extract the public keys, and store them in a slice
	errg := errgroup.Group{}
	containers := make([]testcontainers.Container, N)
	validators := make([]*ValidatorInfo, N)
	var genesis []byte
	for i := 0; i < N; i++ {
		errg.Go(func() error {
			container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
				ContainerRequest: testcontainers.ContainerRequest{
					FromDockerfile: fromDockerfile,
					Entrypoint:     []string{"sh", "-c"},
					Cmd:            []string{"sleep infinity"},
					Name:           fmt.Sprintf("validator-%d", i),
				},
				Started: true,
				Reuse:   true,
			})
			if err != nil {
				return err
			}
			containers[i] = container

			validatorInfo, err := extractValidatorInfo(ctx, container)
			if err != nil {
				return err
			}
			validators[i] = validatorInfo

			if i == 0 {
				genesis, err = extractGenesis(ctx, container)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	err := errg.Wait()
	if err != nil {
		t.Fatal(err)
	}

	// generate genesis
	newGenesis, err := generateGenesis(genesis, validators)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Genesis: %s\n", string(newGenesis))

	// copy genesis to all validators
	for i := 0; i < N; i++ {
		err := containers[i].CopyToContainer(ctx, newGenesis, "/root/.meteorkv/config/genesis.json", 0777)
		if err != nil {
			t.Fatal(err)
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.44"))
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	for i, c := range containers {
		resp, err := cli.ContainerCommit(ctx, c.GetContainerID(), container.CommitOptions{})
		if err != nil {
			fmt.Printf("Error committing container: %v\n", err)
			return
		}

		if err := cli.ImageTag(ctx, resp.ID, fmt.Sprintf("validator-%d", i)); err != nil {
			fmt.Printf("Error tagging image: %v\n", err)
			return
		}
	}

	cli.NetworkRemove(ctx, "testnet")
	netResp, err := cli.NetworkCreate(ctx, "testnet", network.CreateOptions{
		Driver:     "bridge",
		Attachable: true,
		Internal:   true,
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{
				{
					Subnet:  "172.18.0.0/16",
					Gateway: "172.18.0.1",
				},
			},
		},
	})
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		fmt.Printf("Error creating network: %v\n", err)
		return
	}
	fmt.Printf("Created docker network: %s\n", netResp.ID)

	for i := 0; i < N; i++ {
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:      fmt.Sprintf("validator-%d", i),
				Entrypoint: []string{"/main"},
				Cmd:        []string{"start"},
				Name:       fmt.Sprintf("validator-run-%d", i),
			},
			Started: true,
			Reuse:   true,
		})
		if err != nil {
			fmt.Printf("Error creating container: %v\n", err)
			return
		}

		containerID := container.GetContainerID()
		networkConfig := network.EndpointSettings{
			IPAddress: fmt.Sprintf("192.168.1.%d", 100+i),
		}
		err = cli.NetworkConnect(ctx, "testnet", containerID, &networkConfig)
		if err != nil {
			fmt.Printf("Error connecting container to network with fixed IP: %v\n", err)
			return
		}
	}

	time.Sleep(time.Second * 600)
}

type ValidatorInfo struct {
	Address string `json:"address"`
	PubKey  PubKey `json:"pub_key"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func extractValidatorInfo(ctx context.Context, container testcontainers.Container) (*ValidatorInfo, error) {
	code, _, err := container.Exec(ctx, []string{"/main", "init"})
	if err != nil {
		return nil, fmt.Errorf("executing command: %v\n", err)
	}
	if code != 0 {
		return nil, fmt.Errorf("command failed with exit code %d\n", code)
	}

	f, err := container.CopyFileFromContainer(ctx, "/root/.meteorkv/config/priv_validator_key.json")
	if err != nil {
		return nil, fmt.Errorf("copying file from container: %v\n", err)
	}

	fdata, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading file: %v\n", err)
	}

	var validatorInfo ValidatorInfo
	err = json.Unmarshal(fdata, &validatorInfo)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling validator info: %v\n", err)
	}

	return &validatorInfo, nil
}

func extractGenesis(ctx context.Context, container testcontainers.Container) ([]byte, error) {
	f, err := container.CopyFileFromContainer(ctx, "/root/.meteorkv/config/genesis.json")
	if err != nil {
		return nil, fmt.Errorf("copying file from container: %v\n", err)
	}

	fdata, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading file: %v\n", err)
	}

	return fdata, nil
}

type GenesisValidator struct {
	Address string `json:"address"`
	PubKey  PubKey `json:"pub_key"`
	Power   string `json:"power"`
	Name    string `json:"name"`
}

func generateGenesis(genesisBz []byte, validators []*ValidatorInfo) ([]byte, error) {
	var genesis map[string]json.RawMessage
	err := json.Unmarshal(genesisBz, &genesis)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling genesis: %v\n", err)
	}

	genValidators := make([]GenesisValidator, len(validators))
	for i, validator := range validators {
		genValidators[i] = GenesisValidator{
			Address: validator.Address,
			PubKey:  validator.PubKey,
			Power:   "10",
			Name:    fmt.Sprintf("validator-%d", i),
		}
	}

	genValidatorsJson, err := json.Marshal(genValidators)
	if err != nil {
		return nil, fmt.Errorf("marshaling genesis validators: %v\n", err)
	}

	genesis["validators"] = genValidatorsJson

	genesisJson, err := json.Marshal(genesis)
	if err != nil {
		return nil, fmt.Errorf("marshaling genesis: %v\n", err)
	}

	return genesisJson, nil
}
