package checker

import (
	dockercli "availability-checker/pkg/containeractions"
	"availability-checker/pkg/credentialprovider"
	"availability-checker/pkg/database"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	dct "github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLChecker struct {
	Server             string
	Port               string
	DBConnection       database.DBConnection
	CredentialProvider credentialprovider.CredentialProvider
	containerClient    dockercli.DockerClient
}

func (c *MySQLChecker) Name() string {
	return fmt.Sprintf("MySQL: %s:%s", c.Server, c.Port)
}

func (c *MySQLChecker) Check() (bool, error) {
	user, pwd, err := c.CredentialProvider.GetCredentials("mysql")
	if err != nil {
		return false, fmt.Errorf("error getting credentials: %v", err)
	}

	if user == "" || pwd == "" {
		return false, errors.New("empty username or password")
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, pwd, c.Server, c.Port)

	err = c.DBConnection.Open("mysql", connectionString)
	if err != nil {
		fmt.Printf("Error opening connection: %v\n", err)
		return false, err
	}
	defer c.DBConnection.Close()

	err = c.DBConnection.Ping()
	if err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		return false, err
	}

	return true, nil
}

func (c *MySQLChecker) Fix() error {
	ctx := context.TODO()
	if err := c.containerClient.NewClient(); err != nil {
		return err
	}
	defer c.containerClient.Close()

	// Check if the container exists
	containers, err := c.containerClient.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	var containerID string
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/test-mysql" {
				containerID = container.ID
				break
			}
		}
		if containerID != "" {
			break
		}
	}

	if containerID == "" {
		// Container doesn't exist, create it
		reader, err := c.containerClient.ImagePull(ctx, "mysql:latest", types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(ioutil.Discard, reader)
		reader.Close()

		hostConfig := &dct.HostConfig{
			PortBindings: nat.PortMap{
				"3306/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "3306"}},
			},
		}

		resp, err := c.containerClient.ContainerCreate(ctx, &dct.Config{
			Image: "mysql",
			Env:   []string{"MYSQL_ROOT_PASSWORD=superuser"},
			Tty:   true,
			ExposedPorts: nat.PortSet{
				"3306/tcp": struct{}{},
			},
		}, hostConfig, "test-mysql")
		if err != nil {
			return err
		}

		containerID = resp.ID
	}

	if err := c.containerClient.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (c *MySQLChecker) IsFixable() bool {
	return true
}
