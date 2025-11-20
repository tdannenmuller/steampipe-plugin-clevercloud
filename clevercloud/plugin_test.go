package clevercloud

import (
	"context"
	"testing"
)

func TestPlugin(t *testing.T) {
	ctx := context.Background()
	p := Plugin(ctx)

	// Test plugin is not nil
	if p == nil {
		t.Fatal("Plugin should not be nil")
	}

	// Test plugin name
	expectedName := "clevercloud"
	if p.Name != expectedName {
		t.Errorf("Expected plugin name %s, got %s", expectedName, p.Name)
	}

	// Test table map
	if p.TableMap == nil {
		t.Fatal("TableMap should not be nil")
	}

	// Test billing table exists
	if _, exists := p.TableMap["clevercloud_billing"]; !exists {
		t.Error("clevercloud_billing table should exist in TableMap")
	}
}

func TestBillingTable(t *testing.T) {
	ctx := context.Background()
	table := tableCleverCloudBilling(ctx)

	// Test table is not nil
	if table == nil {
		t.Fatal("Billing table should not be nil")
	}

	// Test table name
	expectedName := "clevercloud_billing"
	if table.Name != expectedName {
		t.Errorf("Expected table name %s, got %s", expectedName, table.Name)
	}

	// Test description
	expectedDesc := "Billing invoices information from Clever Cloud."
	if table.Description != expectedDesc {
		t.Errorf("Expected table description '%s', got '%s'", expectedDesc, table.Description)
	}

	// Test columns
	if table.Columns == nil || len(table.Columns) == 0 {
		t.Fatal("Table should have columns")
	}

	// Test expected column names (updated for new schema)
	expectedColumns := []string{"invoice_number", "kind", "owner_id", "emission_date", "status", "currency", "total_amount_excl_tax"}
	columnNames := make([]string, len(table.Columns))
	for i, col := range table.Columns {
		columnNames[i] = col.Name
	}

	for _, expectedCol := range expectedColumns {
		found := false
		for _, actualCol := range columnNames {
			if actualCol == expectedCol {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected column %s not found in table", expectedCol)
		}
	}
}
