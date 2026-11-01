# Resume site

A minimal personal site. A small Go program reads `data/resume.json`, renders `templates/index.html`, and writes plain HTML, CSS and JS to `public/`. GitHub Actions builds and deploys it to GitHub Pages on every push to `main`.

## Layout

```
main.go                      generator (standard library only)
data/resume.json             all the content
templates/index.html         page template
static/style.css             styles (light and dark themes, print styles)
static/main.js               theme toggle
.github/workflows/deploy.yml build and deploy
```

## Before you publish

Open `data/resume.json` and replace the two placeholder links with your real profiles:

```json
{ "label": "LinkedIn", "url": "https://www.linkedin.com/in/your-linkedin-handle" },
{ "label": "GitHub",   "url": "https://github.com/your-github-username" }
```

The build prints a warning while either one still contains `your-`.

## Run locally

Requires Go 1.22 or newer.

```sh
go run .                          # writes public/
cd public && python3 -m http.server 8000
```

Then open http://localhost:8000.

## Deploy to GitHub Pages

1. Create a new repository on GitHub. Naming it `<your-username>.github.io` gives you the address `https://<your-username>.github.io`. Any other name works too and is served at `https://<your-username>.github.io/<repo-name>/`.
2. Push this folder to the `main` branch:

   ```sh
   git init -b main
   git add .
   git commit -m "Add resume site"
   git remote add origin git@github.com:<your-username>/<repo>.git
   git push -u origin main
   ```

3. In the repository, go to **Settings → Pages** and set **Source** to **GitHub Actions**.
4. The **Deploy to GitHub Pages** workflow runs on that push (or run it from the **Actions** tab). When it finishes, the site is live.

From then on, edit `data/resume.json` (or the template and styles), push, and the site updates.

## Custom domain

Add the domain under **Settings → Pages → Custom domain** and create the DNS records GitHub describes there.
