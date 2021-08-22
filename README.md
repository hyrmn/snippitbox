# snippitbox
A very small app to share very small things. Well, so far, a very small app that can get compiled.

This is a single file app written in .NET 6 beta something or other and using ASP.NET Core 6 beta somethingorother. 
Since you probably don't want beta bits on your local machine, you can execute `.\run-docker` to get a tagged image called `snippitbox`.
You can then execute `.\run-docker` to start a container with that tagged image on port 8000 in interactive mode (`ctrl+c` to kill it).

## Prerequisites

The first time you run an ASP.NET Core application locally, you'll want to install the dev cert. You can execute the following command from your favorite command prompt

```powershell
dotnet dev-certs https --trust
```

This will present you with a dialog asking if you're really sure that you want to do it. Yes, you do so click `yes`.

Then, once you run the app and navigate to the site, you'll likely get a browser warning asking you if you're really sure that you want to proceed because the certificate is self-signed. While this does raise
a philosophical question on if someone can ever truly trust themselves, that doesn't help us with moving forward so go into the appropriate section on the browser warning and tell it to trust the certificate.

## Special Thanks

This little app is made possible by:
- [Scriban](https://github.com/scriban/scriban) for templating
- [LiteDB](https://www.litedb.org) for data storage
- [Roland Taylor](https://twitter.com/rolandixor/)'s [Anole](https://github.com/rolandixor/anole) CSS framework.
