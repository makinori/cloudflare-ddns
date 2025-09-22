# Cloudflare DDNS

> 🌎 Automatically update your DNS records when your IP changes

## Installation

```bash
git clone https://github.com/makidoll/cloudflare-ddns.git
cd cloudflare-ddns
cp settings.example.jsonc settings.jsonc
```

Edit `settings.jsonc`. Interval is in minutes.

Highly recommend setting zones to `"hotmilk.space": ["hotmilk.space"]`<br>
then targeting a `*` wildcard, or any subdomains to `@`

You can find your Cloudflare API key here: https://dash.cloudflare.com/profile/api-tokens (Global API Key)

Then use `docker compose up -d` to build and run
