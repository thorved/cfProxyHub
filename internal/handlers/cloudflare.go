package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cfPorxyHub/internal/services"
	"cfPorxyHub/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CloudflareHandler struct {
	cfService *services.CloudflareService
}

// NewCloudflareHandler creates a new Cloudflare handler
func NewCloudflareHandler(cfService *services.CloudflareService) *CloudflareHandler {
	return &CloudflareHandler{
		cfService: cfService,
	}
}

// GetAccounts handles the GET /api/cloudflare/accounts endpoint
func (h *CloudflareHandler) GetAccounts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accounts, err := h.cfService.GetCloudflareAccounts(ctx)
	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch accounts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":  "Accounts retrieved successfully",
		"accounts": accounts,
	})
}

// GetAccountByID handles the GET /api/cloudflare/accounts/:id endpoint
func (h *CloudflareHandler) GetAccountByID(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		utils.ErrorResponse(c, "Account ID is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	account, err := h.cfService.GetCloudflareAccountByID(ctx, accountID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Account retrieved successfully",
		"account": account,
	})
}

// GetAccountsHTML handles HTMX requests for accounts and returns HTML fragments
func (h *CloudflareHandler) GetAccountsHTML(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accounts, err := h.cfService.GetCloudflareAccounts(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, `
			<div class="col-12">
				<div class="alert alert-danger" role="alert">
					<i class="mdi mdi-alert-circle mr-2"></i>
					Failed to fetch accounts: %s
				</div>
			</div>
		`, err.Error())
		return
	}

	if len(accounts) == 0 {
		c.String(http.StatusOK, `
			<div class="col-12">
				<div class="card">
					<div class="card-body text-center py-5">
						<div class="preview-icon bg-light rounded-circle mx-auto mb-3" style="width: 80px; height: 80px; display: flex; align-items: center; justify-content: center;">
							<i class="mdi mdi-account-off text-muted" style="font-size: 32px;"></i>
						</div>
						<h5 class="text-muted mb-2">No Accounts Found</h5>
						<p class="text-muted">No Cloudflare accounts were found. Please check your API credentials.</p>
						<button class="btn btn-primary btn-sm" onclick="location.reload()">
							<i class="mdi mdi-refresh mr-2"></i>Reload Page
						</button>
					</div>
				</div>
			</div>
		`)
		return
	}

	var html string

	// Create a single card with a list of accounts
	html = `<div class="col-12">
		<div class="card">
			<div class="card-body">
				<h5 class="card-title mb-4">
					<i class="mdi mdi-cloud text-info mr-2"></i>
					Available Cloudflare Accounts
				</h5>
				<div class="table-responsive">
					<table class="table table-hover">
						<thead>
							<tr>
								<th>Account</th>
								<th>Type</th>
								<th>Created</th>
								<th>Account ID</th>
								<th>Action</th>
							</tr>
						</thead>
						<tbody>`

	for _, account := range accounts {
		createdDate := "N/A"
		if !account.CreatedOn.IsZero() {
			createdDate = account.CreatedOn.Format("Jan 2, 2006")
		}

		accountType := account.Type
		if accountType == "" {
			accountType = "Standard"
		}

		// Determine type badge color and fix type display
		badgeColor := "badge-secondary"
		typeDisplay := accountType
		switch accountType {
		case "Enterprise":
			badgeColor = "badge-dark"
			typeDisplay = "Enterprise"
		case "Business":
			badgeColor = "badge-success"
			typeDisplay = "Business"
		case "Pro":
			badgeColor = "badge-info"
			typeDisplay = "Pro"
		case "Standard":
			badgeColor = "badge-secondary"
			typeDisplay = "Standard"
		default:
			badgeColor = "badge-light text-dark"
			typeDisplay = "Standard"
		}

		html += fmt.Sprintf(`
							<tr class="account-row" data-account-id="%s" data-account-name="%s">
								<td>
									<div class="d-flex align-items-center">
										<div class="preview-icon bg-gradient-primary rounded-circle mr-3" style="width: 40px; height: 40px; display: flex; align-items: center; justify-content: center;">
											<i class="mdi mdi-account-circle text-white"></i>
										</div>
										<div>
											<h6 class="mb-1 font-weight-medium">%s</h6>
											<p class="text-muted small mb-0">Active Account</p>
										</div>
									</div>
								</td>
								<td>
									<span class="badge %s px-3 py-1">%s</span>
								</td>
								<td class="text-muted">%s</td>
								<td>
									<code class="small bg-light text-dark px-2 py-1 rounded">%s</code>
								</td>
								<td>
									<button class="btn btn-outline-success btn-sm select-account-btn" 
											data-account-id="%s" 
											data-account-name="%s">
										<i class="mdi mdi-check-circle mr-1"></i>
										Select
									</button>
								</td>
							</tr>`,
			account.ID, account.Name, account.Name, badgeColor, typeDisplay,
			createdDate, account.ID[:8]+"...", account.ID, account.Name)
	}

	html += `
						</tbody>
					</table>
				</div>
			</div>
		</div>
	</div>`

	c.String(http.StatusOK, html)
}
