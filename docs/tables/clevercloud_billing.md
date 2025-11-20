# clevercloud_billing

Lists billing invoices for your Clever Cloud organization based on the actual API schema.

## Configuration

Configure via Steampipe config file (`~/.steampipe/config/clevercloud.spc`):

```hcl
connection "clevercloud" {
  plugin = "clevercloud"
  token = "your_api_token"
  organization_id = "your_organization_id"
}
```

## Examples

### Recent invoices

```sql
select invoice_number, emission_date, total_amount_excl_tax, currency, status
from clevercloud_billing 
order by emission_date desc 
limit 10;
```

### Total spending by month

```sql
select 
  date_trunc('month', emission_date) as month,
  sum(total_amount_excl_tax) as total_amount,
  currency
from clevercloud_billing 
where status = 'PAID'
group by date_trunc('month', emission_date), currency
order by month desc;
```

### Invoices by billing company  

```sql
select 
  billing_company,
  count(*) as invoice_count,
  sum(total_amount_excl_tax) as total_amount
from clevercloud_billing
where billing_company is not null
group by billing_company;
```

## Schema

| Column Name | Type | Description |
|-------------|------|-------------|
| invoice_number | text | Unique invoice identifier |
| kind | text | Type of invoice |
| owner_id | text | Organization owner ID |
| emission_date | timestamp | Date the invoice was issued |
| pay_date | timestamp | Date the invoice was paid |
| status | text | Payment status (PAID, PENDING, etc.) |
| currency | text | Currency code (EUR, USD, etc.) |
| category | text | Service category (PAAS, etc.) |
| total_amount_excl_tax | double | Total amount excluding tax |
| total_tax_amount | double | Tax amount |
| total_amount_display | text | Total amount formatted with currency |
| total_tax_display | text | Tax amount formatted with currency |
| vat_percent | integer | VAT percentage applied |
| discount | integer | Discount applied |
| price_factor | integer | Price factor applied |
| kpi_compute_months | integer | Number of months for KPI computation |
| payment_provider | text | Payment provider (stripe, etc.) |
| transaction_id | text | Payment transaction ID |
| customer_order_id | text | Customer order ID |
| billing_name | text | Name on billing address |
| billing_company | text | Company name on billing address |
| billing_address | text | Street address |
| billing_city | text | Billing city |
| billing_zipcode | text | Postal code |
| billing_country | text | Country code |
| billing_vat_number | text | VAT number |
