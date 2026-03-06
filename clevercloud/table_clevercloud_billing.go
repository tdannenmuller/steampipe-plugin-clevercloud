package clevercloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// tableCleverCloudBilling defines the schema for the billing table
func tableCleverCloudBilling(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "clevercloud_billing",
		Description: "Billing invoices information from Clever Cloud.",
		List: &plugin.ListConfig{
			// Fetch the billing data from the Clever Cloud API
			Hydrate: listInvoices,
		},
		Columns: []*plugin.Column{
			{Name: "invoice_number", Type: proto.ColumnType_STRING, Description: "The unique invoice number.", Transform: transform.FromField("InvoiceNumber")},
			{Name: "kind", Type: proto.ColumnType_STRING, Description: "The type of invoice.", Transform: transform.FromField("Kind")},
			{Name: "owner_id", Type: proto.ColumnType_STRING, Description: "The organization owner ID.", Transform: transform.FromField("OwnerID")},
			{Name: "emission_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date when the invoice was issued.", Transform: transform.FromField("EmissionDate")},
			{Name: "pay_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date when the invoice was paid.", Transform: transform.FromField("PayDate")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The payment status of the invoice.", Transform: transform.FromField("Status")},
			{Name: "currency", Type: proto.ColumnType_STRING, Description: "The currency of the invoice amounts.", Transform: transform.FromField("Currency")},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "The service category (e.g., PAAS).", Transform: transform.FromField("Category")},

			// Amount fields
			{Name: "total_amount_excl_tax", Type: proto.ColumnType_DOUBLE, Description: "Total amount excluding tax.", Transform: transform.FromField("TotalTaxExcluded.Amount")},
			{Name: "total_tax_amount", Type: proto.ColumnType_DOUBLE, Description: "Total tax amount.", Transform: transform.FromField("TotalTax.Amount")},
			{Name: "total_amount_display", Type: proto.ColumnType_STRING, Description: "Total amount excluding tax (formatted).", Transform: transform.FromField("TotalTaxExcluded.DefaultDisplay")},
			{Name: "total_tax_display", Type: proto.ColumnType_STRING, Description: "Total tax amount (formatted).", Transform: transform.FromField("TotalTax.DefaultDisplay")},

			// Tax and discount info
			{Name: "vat_percent", Type: proto.ColumnType_DOUBLE, Description: "VAT percentage applied.", Transform: transform.FromField("VatPercent")},
			{Name: "discount", Type: proto.ColumnType_DOUBLE, Description: "Discount applied to the invoice.", Transform: transform.FromField("Discount")},
			{Name: "price_factor", Type: proto.ColumnType_DOUBLE, Description: "Price factor applied.", Transform: transform.FromField("PriceFactor")},

			// Payment info
			{Name: "payment_provider", Type: proto.ColumnType_STRING, Description: "Payment provider used (e.g., stripe).", Transform: transform.FromField("PaymentProvider")},
			{Name: "transaction_id", Type: proto.ColumnType_STRING, Description: "Payment transaction ID.", Transform: transform.FromField("TransactionID")},
			{Name: "customer_order_id", Type: proto.ColumnType_STRING, Description: "Customer order ID.", Transform: transform.FromField("CustomerOrderID")},

			// Address information (flattened for easier querying)
			{Name: "billing_name", Type: proto.ColumnType_STRING, Description: "Name on the billing address.", Transform: transform.FromField("Address.Name")},
			{Name: "billing_company", Type: proto.ColumnType_STRING, Description: "Company name on the billing address.", Transform: transform.FromField("Address.Company")},
			{Name: "billing_address", Type: proto.ColumnType_STRING, Description: "Billing street address.", Transform: transform.FromField("Address.Address")},
			{Name: "billing_city", Type: proto.ColumnType_STRING, Description: "Billing city.", Transform: transform.FromField("Address.City")},
			{Name: "billing_zipcode", Type: proto.ColumnType_STRING, Description: "Billing postal code.", Transform: transform.FromField("Address.Zipcode")},
			{Name: "billing_country", Type: proto.ColumnType_STRING, Description: "Billing country code.", Transform: transform.FromField("Address.CountryAlpha2")},
			{Name: "billing_vat_number", Type: proto.ColumnType_STRING, Description: "VAT number on the billing address.", Transform: transform.FromField("Address.VatNumber")},
		},
	}
}

// listInvoices fetches the invoice data from the Clever Cloud API
func listInvoices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get configuration from connection
	var config cleverCloudConfig
	if d.Connection != nil && d.Connection.Config != nil {
		config = d.Connection.Config.(cleverCloudConfig)
	}

	// Try to fetch real data if credentials are provided
	if config.Token != nil && config.OrganizationID != nil {
		plugin.Logger(ctx).Info("Fetching invoices from Clever Cloud API", "organization_id", *config.OrganizationID)

		endpoint := "https://api-bridge.clever-cloud.com"
		if config.APIEndpoint != nil {
			endpoint = *config.APIEndpoint
		}

		invoices, err := fetchInvoiceData(ctx, *config.Token, *config.OrganizationID, endpoint)
		if err != nil {
			plugin.Logger(ctx).Error("Failed to fetch invoices from API", "error", err)
			return nil, fmt.Errorf("failed to fetch invoices: %w", err)
		}

		plugin.Logger(ctx).Info("Successfully fetched invoices", "count", len(invoices))

		for _, invoice := range invoices {
			d.StreamListItem(ctx, invoice)
		}
		return nil, nil
	}

	// Log warning about missing credentials and return empty result
	plugin.Logger(ctx).Warn("Connection not configured with token and organization_id, no data will be returned")
	return nil, nil
}

// Address represents the billing address information
type Address struct {
	AddressID             string  `json:"address_id"`
	OwnerID               string  `json:"owner_id"`
	Name                  string  `json:"name"`
	Company               *string `json:"company"`
	Address               string  `json:"address"`
	City                  string  `json:"city"`
	Zipcode               string  `json:"zipcode"`
	CountryAlpha2         string  `json:"country_alpha2"`
	VatNumber             *string `json:"vat_number"`
	VatPercent            float64 `json:"vat_percent"`
	CustomerCostCenter    *string `json:"customer_cost_center"`
	CustomerPurchaseOrder *string `json:"customer_purchase_order"`
}

// MonetaryAmount represents an amount with currency information
type MonetaryAmount struct {
	Currency        string  `json:"currency"`
	Amount          float64 `json:"amount"`
	AmountFormatted string  `json:"amount_formatted"`
	DefaultDisplay  string  `json:"default_display"`
}

// Invoice represents a Clever Cloud invoice (renamed from Billing)
type Invoice struct {
	InvoiceNumber    string         `json:"invoice_number"`
	Kind             string         `json:"kind"`
	OwnerID          string         `json:"owner_id"`
	Address          Address        `json:"address"`
	EmissionDate     string         `json:"emission_date"`
	PayDate          *string        `json:"pay_date"`
	Status           string         `json:"status"`
	Currency         string         `json:"currency"`
	PriceFactor      float64        `json:"price_factor"`
	Discount         float64        `json:"discount"`
	VatPercent       float64        `json:"vat_percent"`
	TotalTaxExcluded MonetaryAmount `json:"total_tax_excluded"`
	TotalTax         MonetaryAmount `json:"total_tax"`
	Category         string         `json:"category"`
	PaymentProvider  *string        `json:"payment_provider"`
	TransactionID    *string        `json:"transaction_id"`
	CustomerOrderID  *string        `json:"customer_order_id"`
}

func fetchInvoiceData(ctx context.Context, token string, organizationID string, apiEndpoint string) ([]Invoice, error) {
	// Build the full URL for the invoices endpoint
	url := fmt.Sprintf("%s/v4/billing/organisations/%s/invoices", apiEndpoint, organizationID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	var invoiceData []Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoiceData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return invoiceData, nil
}
