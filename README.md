# Snippitbox
I'm doing some very small and very focused playing around with ASP.NET CORE on .NET 6. This app was intented to provide some very basic, and not full-featured, functionality with the new Minimal API work. Snippitbox will take in some data and save it to a small local database. There is no validation or search. Heck, I don't even have paging wired up (although it is possible and the query supports it). I might come back to that, but so far I've gotten to where I wanted to be... tiny app that works and that can be built on Docker (or your favorite OCI-compatible host that can read Dockerfiles).

This is a single file app written in .NET 6 RC1 and using ASP.NET Core. If you don't have .NET 6.0 RC1, or later, you can install it from the nightly installer links near the bottom of https://github.com/dotnet/installer

Since you probably don't want beta bits on your local machine, you can execute `.\build-docker` to get a tagged image called `snippitbox`.
You can then execute `.\run-docker` to start a container with that tagged image on port 8000 in interactive mode (`ctrl+c` to kill it).

If you want a persistent database then you will want to pass in a volume mount for `/app/data`. My "prod" docker-compose does this and you can run `docker-compose -f .\docker-compose.dev.yaml up --build` 
to start a docker container on port 5020. The LiteDB database will be saved to the `db` host volume mount.

## Prerequisites

The first time you run an ASP.NET Core application locally, you'll want to install the dev cert. You can execute the following command from your favorite command prompt

```powershell
dotnet dev-certs https --trust
```

This will present you with a dialog asking if you're really sure that you want to do it. Yes, you do so click `yes`.

Then, once you run the app and navigate to the site, you'll likely get a browser warning asking you if you're really sure that you want to proceed because the certificate is self-signed. While this does raise
a philosophical question on if someone can ever truly trust themselves, that doesn't help us with moving forward so go into the appropriate section on the browser warning and tell it to trust the certificate.

## My Docker build won't work

(As of 2021-09-06), this code relies on the .NET 6.0 RC1 release. I've chosen to not pin to an image version. So, if you run into any issues then you should probably update your images like so

```powershell
> docker pull mcr.microsoft.com/dotnet/sdk:6.0-alpine
> docker pull mcr.microsoft.com/dotnet/aspnet:6.0-alpine
```

## Special Thanks

This little app is made possible by:
- [Scriban](https://github.com/scriban/scriban) for templating
- [LiteDB](https://www.litedb.org) for data storage
- [Roland Taylor](https://twitter.com/rolandixor/)'s [Anole](https://github.com/rolandixor/anole) CSS framework.
- Favicon is from [iconpacks.net](https://www.iconpacks.net)