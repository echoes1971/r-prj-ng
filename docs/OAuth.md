OAuth setup for ρBee (rhobee)
=================================

This document shows where to find the Google and GitHub OAuth parameters and how to configure them for production testing.

Files referenced
- Production compose: [docker-compose.yml](docker-compose.yml)
- Backend config (JSON): [be/config.json](be/config.json)

Callback endpoints (server):
- Google callback: `/oauth/google/callback`
- GitHub callback: `/oauth/github/callback`
- Telegram callback: `/oauth/telegram/callback`

1) Google: obtain `google_client_id` and `google_client_secret`

- Go to https://console.cloud.google.com/
- Create or select a Project.
- In the left menu go to "APIs & Services" → "Credentials".
- Create Credentials → OAuth 2.0 Client IDs → Web application.
- Set the Authorized redirect URI to the production callback URL, for example:

  https://your-domain.example.com/oauth/google/callback

- After creation note the `Client ID` and `Client secret`.

2) GitHub: obtain `github_client_id` and `github_client_secret`

- Go to https://github.com/settings/developers
- Click on "OAuth Apps" → "New OAuth App".
- Set the Application callback URL to the production callback URL, for example:

  https://your-domain.example.com/oauth/github/callback

- After creation note the `Client ID` and `Client secret`.

3) Telegram: obtain `telegram_bot_token`

Telegram uses a different flow (Login Widget) but is equally simple:

- Open Telegram and search for [@BotFather](https://t.me/botfather)
- Send `/newbot` and follow instructions to create a new bot
- BotFather will give you a **Bot Token** (e.g., `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`)
- Save this token — you'll use it as `telegram_bot_token` in backend config
- **Important**: Send `/setdomain` to @BotFather and set your production domain (e.g., `your-domain.example.com`) to allow the Login Widget to work

**How Telegram OAuth works:**
- No redirect to Telegram's site like Google/GitHub
- Frontend shows a Telegram Login Widget button
- User clicks → opens Telegram app/web to authorize
- Telegram redirects back to your callback URL with user data
- Backend verifies the data hash using HMAC-SHA256 with your bot token
- User is created/logged in and receives a JWT token

**Testing locally:**
- For local testing, you can use `localhost:8080` as domain with @BotFather
- The widget will still work but you may see warnings

4) Where to place the values

- For quick production deployment we added placeholders to the production compose file. Edit the `be` service environment values in [docker-compose.yml](docker-compose.yml) and replace the placeholder values:

  - GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID
  - GOOGLE_CLIENT_SECRET=YOUR_GOOGLE_CLIENT_SECRET
  - GOOGLE_REDIRECT_URL=https://your-domain.example.com/oauth/google/callback
  - GITHUB_CLIENT_ID=YOUR_GITHUB_CLIENT_ID
  - GITHUB_CLIENT_SECRET=YOUR_GITHUB_CLIENT_SECRET
  - GITHUB_REDIRECT_URL=https://your-domain.example.com/oauth/github/callback
  - TELEGRAM_BOT_TOKEN=YOUR_TELEGRAM_BOT_TOKEN

- The backend currently reads configuration from `be/config.json` (see [be/config.json](be/config.json)). You can also put the keys directly into that JSON file under the following fields:

```json
{
  "google_client_id": "...",
  "google_client_secret": "...",
  "google_redirect_url": "https://your-domain.example.com/oauth/google/callback",
  "github_client_id": "...",
  "github_client_secret": "...",
  "github_redirect_url": "https://your-domain.example.com/oauth/github/callback",
  "telegram_bot_token": "123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
}
```

**Note:** The Telegram Bot ID (numeric part before `:`) is automatically extracted from the token and exposed to the frontend via `REACT_APP_TELEGRAM_BOT_ID` environment variable. You don't need to configure it separately.
```

Note: the app currently loads `be/config.json` at startup. The compose environment placeholders are convenient helper values for production deployment; if you prefer to keep secrets out of the compose file, put them in a secure place and populate `be/config.json` before starting the service.

5) Frontend configuration

To enable OAuth buttons in the login page, set these environment variables for the frontend container:

- `REACT_APP_ENABLE_GOOGLE_OAUTH=true` (or `1`)
- `REACT_APP_ENABLE_GITHUB_OAUTH=true` (or `1`)
- `REACT_APP_ENABLE_TELEGRAM_OAUTH=true` (or `1`)

These flags control which OAuth buttons appear on the `/login` page.

6) Redirect / UX notes

- The server implements the OAuth callback endpoints at `/oauth/google/callback`, `/oauth/github/callback`, and `/oauth/telegram/callback`.
- All callbacks create a minimal user account (if not exists) and return an HTML page that:
  - Stores `access_token`, `expires_at`, `user_id`, `username` (login), and `groups` in `localStorage`
  - Redirects to the main app route (`/`)
- Telegram flow: user clicks button → opens Telegram app/web → authorizes → redirects back with query params → backend verifies hash → returns same HTML response

Redirect behavior implemented
- The server now redirects the browser back to the frontend root (`/`) and includes the token information inside the URL fragment (hash) to avoid sending tokens to the server in query parameters. Example fragment:

  #provider=google&access_token=<TOKEN>&expires_at=1670000000&user_id=123

- Frontend should read `window.location.hash`, parse the fragment, store the token (e.g., in `localStorage`) and then navigate to the main app route. Example (simplified):

```js
const hash = window.location.hash.substring(1);
const params = new URLSearchParams(hash);
const token = params.get('access_token');
if (token) {
  localStorage.setItem('token', token);
  // optional: store expires_at and user_id
  window.location.hash = '';
  window.location.href = '/';
}
```

7) Testing locally

- If you test locally and expose the service via a domain or use a local tunneling service (ngrok, localtunnel), register the tunnel URL as the authorized redirect URI in the provider settings.

8) Security

- Use HTTPS in production for OAuth redirect URIs.
- Protect client secrets; prefer secrets management or environment injection by your deployment tool instead of committing secrets to the repository.

If you want, I can also:
- add the same placeholders to `docker-compose.dev.yml` for your local hot-reload setup, or
- implement the frontend `/login` button which redirects back to the main page after successful login.

--
Generated instructions for testing and deployment.
