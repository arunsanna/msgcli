# Azure App Registration Setup

This guide walks you through creating an Azure app registration for msgcli. This is a one-time setup that takes about 5 minutes.

## Prerequisites

- A Microsoft account (personal @outlook.com/@hotmail.com, or work/school Microsoft 365)
- Access to [Azure Portal](https://portal.azure.com)

## Step 1: Create App Registration

1. Go to [Azure Portal](https://portal.azure.com)
2. Search for "App registrations" in the top search bar
3. Click **"New registration"**

## Step 2: Configure Registration

Fill in the form:

| Field | Value |
|-------|-------|
| **Name** | `msgcli` (or any name you prefer) |
| **Supported account types** | Select: **"Accounts in any organizational directory and personal Microsoft accounts"** |
| **Redirect URI** | Leave blank (not needed for device code flow) |

Click **"Register"**

## Step 3: Copy Client ID

After registration, you'll see the app overview page.

1. Find **"Application (client) ID"** - it looks like a GUID: `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
2. Copy this value - you'll need it for msgcli

## Step 4: Enable Public Client Flow

1. In the left sidebar, click **"Authentication"**
2. Scroll down to **"Advanced settings"**
3. Find **"Allow public client flows"**
4. Set it to **"Yes"**
5. Click **"Save"** at the top

## Step 5: Configure msgcli

```bash
# Store your client ID
msgcli auth setup
# Paste the client ID when prompted

# Add your first account
msgcli auth add personal
# Follow the device code flow in your browser
```

## That's It!

No additional API permissions need to be configured. msgcli uses "dynamic consent" which requests permissions when you first authenticate.

## Permissions Requested

When you authenticate, you'll be asked to grant these permissions:

| Permission | Purpose |
|------------|---------|
| **Mail.ReadWrite** | Read, create, update, delete emails |
| **Mail.Send** | Send emails |
| **Calendars.ReadWrite** | Read, create, update, delete calendar events |
| **User.Read** | Read your basic profile info |
| **offline_access** | Stay signed in (refresh tokens) |

## Troubleshooting

### "AADSTS65001: The user or administrator has not consented to use the application"

This usually means:
- The app registration wasn't created correctly
- "Allow public client flows" is not enabled
- Your organization restricts third-party app access

**Solution**: Verify the setup steps above, or ask your IT admin to approve the app.

### "AADSTS7000218: The request body must contain the following parameter: 'client_assertion' or 'client_secret'"

This means "Allow public client flows" is not enabled.

**Solution**: Go to Authentication → Enable "Allow public client flows" → Save

### Device code flow times out

The code expires after 15 minutes. Make sure to complete the browser authentication promptly.

### Token refresh fails

Your refresh token may have expired (typically after 90 days of inactivity).

**Solution**: Re-authenticate with `msgcli auth add <alias>`

## Multiple Accounts

You can add multiple accounts with different aliases:

```bash
msgcli auth add personal   # Personal @outlook.com account
msgcli auth add work       # Work @company.com account

# Use specific account
msgcli mail list -a personal
msgcli calendar list -a work
```

## Security Notes

- Your tokens are stored in your OS keychain (macOS Keychain, Windows Credential Manager, or Linux Secret Service)
- The client ID is not a secret - it only identifies your app registration
- No client secret is used (public client flow)
- Tokens are refreshed automatically when they expire
