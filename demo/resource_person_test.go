package demo

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

func TestAccPerson_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	rn := "demo_person.person"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDemoPersonConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDemoPersonExists(rn),
					resource.TestCheckResourceAttr(rn, "age", "50"),
					resource.TestCheckResourceAttr(rn, "hobby", "demos"),
				),
			},
			{
				Config: testAccDemoPersonUpdatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDemoPersonExists(rn),
					resource.TestCheckResourceAttr(rn, "hobby", "demos-updated"),
				),
			},
		},
	})
}

func testAccCheckDemoPersonExists(demoPersonPath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[demoPersonPath]
		if !ok {
			return fmt.Errorf("not found: %s", demoPersonPath)
		}
		personName := rs.Primary.ID
		client := testAccProvider.Meta().(*elasticsearch.Client)
		response, e := client.Get(index, personName)
		if e != nil {
			return e
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("response for person [%s] is not %v, got: %v", personName, http.StatusOK, response.StatusCode)
		}
		return nil
	}
}

func testAccDemoPersonConfig(name string) string {
	return fmt.Sprintf(`
resource "demo_person" "person" {
  name  = "demo-%s"
  age   = 50
  hobby = "demos"
}`, name)
}

func testAccDemoPersonUpdatedConfig(name string) string {
	return fmt.Sprintf(`
resource "demo_person" "person" {
  name  = "demo-%s"
  age   = 50
  hobby = "demos-updated"
}`, name)
}
