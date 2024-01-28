/*
Released under YOLO licence. Idgaf what you do.
*/
package azrecon

import (
	"testing"
)

var (
	resources = []string{
		"vnet-prod-westus-001",
		"vnet-pord-westus-002",
		"vnet-prod-westus-003",
		"vnet-prod-westus-004",
		"vnet-prod-eastus-001",
		"vnet-prod-eastus-002",
		"vnet-prod-eastus-003",
		"vnet-prod-eastus-004",
		"nsg-prod-westus-001",
		"nsg-pord-westus-002",
		"nsg-prod-westus-003",
		"nsg-prod-westus-004",
		"nsg-prod-eastus-001",
		"nsg-prod-eastus-002",
		"nsg-prod-eastus-003",
		"nsg-prod-eastus-004",
		"vnet-dev-westus-001",
		"vnet-dev-westus-002",
		"vnet-dev-westus-003",
		"vnet-dev-westus-004",
		"sql-prod-westus-001",
		"sql-prod-westus-002",
		"sql-uat-westus-001",
		"sql-uat-westus-002",
		"sql-uat-westus-002-dev", // curveball entry
	}
)

func TestTestycoolsName(t *testing.T) {
	doubleSlices, options := getColumns(resources, Name)
	if len(doubleSlices) == 0 {
		t.Fatalf("doubleSlices = %q, optionals = %q", doubleSlices, options)
	} else {
		t.Logf("doubleSlices = %q, optionals = %q", doubleSlices, options)
	}
}

func TestTestycoolsMask(t *testing.T) {
	doubleSlices, options := getColumns(resources, Mask)
	if len(doubleSlices) == 0 {
		t.Fatalf("doubleSlices = %q, optionals = %q", doubleSlices, options)
	} else {
		t.Logf("doubleSlices = %q, optionals = %q", doubleSlices, options)
	}
}

func TestGetAllNameCombinations(t *testing.T) {
	combos := GetNameCombinations(resources)

	if len(combos) == 0 {
		t.Fatalf("combos = %q\n", combos)
	} else {
		t.Logf("combos = %q\n", combos)
	}

}

func TestGetAllMaskCombinations(t *testing.T) {
	combos := GetMaskCombinations(resources)

	if len(combos) == 0 {
		t.Fatalf("combos = %q\n", combos)
	} else {
		t.Logf("combos = %q\n", combos)
	}
}
