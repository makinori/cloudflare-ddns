# Cloudflare DDNS

> 🌎 Automatically update your DNS records when your IP changes

## Installation

Firstly make sure you have **Deno** and **Git** installed.

```bash
git clone https://github.com/makidoll/cloudflare-ddns.git
cd cloudflare-ddns
cp settings.example.jsonc settings.jsonc
```

Edit `settings.jsonc`. Interval is in minutes.

Highly recommend setting zones to `"hotmilk.space": ["hotmilk.space"]`<br>
then targeting a `*` wildcard, or any subdomains to `@`

You can find your Cloudflare API key here: https://dash.cloudflare.com/profile/api-tokens (Global API Key)

Run `deno run -A main.ts` indefinitely or use a process manager like `pm2`.

```bash
npm install pm2 -g
pm2 start ecosystem.config.js
pm2 ls
```
