# GlutenScout

Helps gluten-free diners find restaurant meals other gluten-free people have reviewed.

> Sample data is **fictional**. Always confirm with a restaurant before ordering.

The app is the React + TypeScript prototype in [`frontend/`](frontend/). It's live at
<https://gluten-scout-mvp.vercel.app>. See [frontend/README.md](frontend/README.md) to run it,
test it and deploy it.

## Older Go version

The first MVP (Go, HTMX and PostgreSQL) lives on the
[`go-mvp`](https://github.com/Master-004/glutenScoutMVP/tree/go-mvp) branch. To work on it in its
own folder next to this one:

```sh
git worktree add ../glutenScoutMVP-go go-mvp
```

## Repo layout

```
frontend/   the React app (deployed by Vercel)
docs/       design notes and the original slice specs
.github/    CI checks, Dependabot, PR template
```
