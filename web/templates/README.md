# Template Organization

This directory contains all HTML templates organized in a hierarchical structure for better maintainability.

## Directory Structure

```
templates/
├── auth/                    # Authentication-related pages
│   └── login.html          # User login form
├── cloudflare/             # Cloudflare management pages
│   ├── accounts/           # Account management templates
│   │   └── CloudflareAccounts.html
│   ├── tunnels/            # Tunnel management templates
│   │   ├── CloudflareAllTunnels.html
│   │   ├── Cloudflare_CreateTunnel.html
│   │   └── Cloudflare_TunnelPublicHostname.html
│   └── zones/              # Zone management templates
│       ├── CloudflareZones.html
│       └── CloudflareZoneDetails.html
├── dashboard/              # Dashboard and main application pages
│   └── Dashboard.html      # Main dashboard page
└── layouts/                # Reusable layout components
    ├── header.html         # Common header template
    ├── footer.html         # Common footer template
    └── sidebar.html        # Navigation sidebar template

```

## Template Usage

### Layout Components
Layout components in the `layouts/` folder are included using Go template syntax with just the filename:
```html
{{template "header.html" .}}
{{template "sidebar.html" .}}
{{template "footer.html" .}}
```
**Note:** Gin loads templates by filename only, regardless of their directory structure.

### Template Loading
Templates are loaded using a glob pattern that includes all subdirectories:
```go
router.LoadHTMLGlob("web/templates/**/*.html")
```

### Go Code References
Templates are referenced in Go code using their relative path from the templates directory:
```go
c.HTML(http.StatusOK, "auth/login.html", gin.H{})
c.HTML(http.StatusOK, "dashboard/Dashboard.html", gin.H{})
c.HTML(http.StatusOK, "cloudflare/accounts/CloudflareAccounts.html", gin.H{})
c.HTML(http.StatusOK, "cloudflare/tunnels/CloudflareAllTunnels.html", gin.H{})
c.HTML(http.StatusOK, "cloudflare/zones/CloudflareZones.html", gin.H{})
```

## Important Gin Framework Limitation

**Gin Template Name Collision**: Gin framework loads templates by their filename only, not their full path. This means that even though we've organized templates into folders for better structure, all template names must still be unique across all directories. The folder organization helps with maintenance and logical grouping, but doesn't affect how templates are referenced in Go code or template includes.

## Benefits of This Organization

1. **Better Organization**: Related templates are grouped together
2. **Scalability**: Easy to add new categories as the application grows
3. **Maintainability**: Clear separation of concerns
4. **Reusability**: Layout components are centralized and easily shared
5. **Navigation**: Easier to find specific templates

## Adding New Templates

When adding new templates:
1. Place them in the appropriate category folder
2. Create new category folders if needed
3. Update the Go route handlers to use the correct path
4. Include layout components using the `layouts/` prefix
