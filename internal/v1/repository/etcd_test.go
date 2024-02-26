package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/denizgursoy/clerk/internal/v1/config"
	"github.com/denizgursoy/clerk/internal/v1/usecases"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ETCDTestSuite struct {
	suite.Suite
	container  testcontainers.Container
	etcdClient *clientv3.Client
	r          *MemberETCDRepository
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ETCDTestSuite))
}

func (s *ETCDTestSuite) SetupSuite() {
	port, err := s.startTestContainer()
	require.NoError(s.T(), err)
	client, err := CreateETCDClient(config.Config{
		ETCDEndpoint: fmt.Sprintf("localhost:%s", port),
	})
	require.NoError(s.T(), err)
	require.NotNil(s.T(), client)
	s.etcdClient = client
	s.r = NewMemberETCDRepository(client)
}

func (s *ETCDTestSuite) TearDownSuite() {}

func (s *ETCDTestSuite) startTestContainer() (string, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "quay.io/coreos/etcd:v3.4.0",
		ExposedPorts: []string{"2379/tcp"},
		Cmd:          []string{"etcd", "--listen-client-urls", "http://0.0.0.0:2379", "--advertise-client-urls", "http://0.0.0.0:2379"},
		WaitingFor:   wait.ForLog("ready to serve client requests"),
	}
	etcdContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	p, err := etcdContainer.MappedPort(ctx, "2379")
	if err != nil {
		return "", err
	}
	time.Sleep(time.Second)
	s.container = etcdContainer

	return p.Port(), nil
}

func (s *ETCDTestSuite) SetupTest() {
	// clear all record before every session
	_, err := s.etcdClient.Delete(context.Background(), usecases.IDPrefix, clientv3.WithPrefix())
	require.NoError(s.T(), err)
}

func (s *ETCDTestSuite) TearDownTest() {
}
