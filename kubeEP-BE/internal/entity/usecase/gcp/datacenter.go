package gcpUCEntity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
)

type DatacenterData struct {
	Credentials json.RawMessage
	Name        string
}

type DatacenterDetailedData struct {
	ID          uuid.UUID
	Name        string
	Credentials json.RawMessage
	Metadata    json.RawMessage
	Datacenter  constant.DatacenterType
}

type SAKeyCredentials struct {
	Type                    *string `json:"type" validate:"required"`
	ProjectId               *string `json:"project_id" validate:"required"`
	PrivateKeyId            *string `json:"private_key_id" validate:"required"`
	PrivateKey              *string `json:"private_key" validate:"required"`
	ClientEmail             *string `json:"client_email" validate:"required"`
	ClientId                *string `json:"client_id" validate:"required"`
	AuthUri                 *string `json:"auth_uri" validate:"required"`
	TokenUri                *string `json:"token_uri" validate:"required"`
	AuthProviderX509CertUrl *string `json:"auth_provider_x509_cert_url" validate:"required"`
	ClientX509CertUrl       *string `json:"client_x509_cert_url" validate:"required"`
}

type DatacenterMetaData struct {
	ProjectId string `json:"project_id"`
	SAEmail   string `json:"sa_email"`
}
