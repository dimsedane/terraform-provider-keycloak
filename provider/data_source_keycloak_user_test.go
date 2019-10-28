package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakDataSourceUser_basic(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	username := "username-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakUser_basic(realm, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists("keycloak_user.user"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "id", "data.keycloak_user.user", "id"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "realm_id", "data.keycloak_user.user", "realm_id"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "username", "data.keycloak_user.user", "username"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "description", "data.keycloak_user.user", "description"),
					testAccCheckDataKeycloakUser("data.keycloak_user.user"),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakUser(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]
		username := rs.Primary.Attributes["username"]

		user, err := keycloakClient.GetUser(realmId, id)
		if err != nil {
			return err
		}

		if user.Username != username {
			return fmt.Errorf("expected user with ID %s to have username %s, but got %s", id, username, user.Username)
		}

		return nil
	}
}

func testDataSourceKeycloakUser_basic(realm, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id    = "${keycloak_realm.realm.id}"
	username    = "%s"
}

data "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username     = "${keycloak_user.user.username}"
}
	`, realm, username)
}
