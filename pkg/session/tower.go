package session

import (
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/haijianyang/cloudtower-go-sdk/client"
	"github.com/haijianyang/cloudtower-go-sdk/client/operations"
	"github.com/haijianyang/cloudtower-go-sdk/models"

	infrav1 "github.com/smartxworks/cluster-api-provider-elf/api/v1alpha4"
)

type TowerSession struct {
	operations.ClientService
}

func NewTowerSession(tower infrav1.Tower) (*TowerSession, error) {
	transport := httptransport.New(fmt.Sprint(tower.Server, ":", tower.Port), "/", []string{"http"})
	client := apiclient.New(transport, strfmt.Default)

	loginParams := operations.NewLoginParams()
	loginParams.RequestBody = &models.LoginInput{
		Username: &tower.Username,
		Password: &tower.Password,
		Source:   models.NewUserSource(models.UserSourceLOCAL),
	}

	loginResp, err := client.Operations.Login(loginParams)
	if err != nil {
		return nil, err
	}

	token := httptransport.BearerToken(*loginResp.Payload.Data.Token)
	transport.DefaultAuthentication = token
	client = apiclient.New(transport, strfmt.Default)

	return &TowerSession{client.Operations}, nil
}
