package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakUserRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeycloakUserRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	username := data.Get("username").(string)

	user, err := keycloakClient.GetUserByUsername(realmId, username)
	if err != nil {
		return err
	}

	mapFromUserToData(data, user)

	return nil
}
