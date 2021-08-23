FROM mcr.microsoft.com/dotnet/sdk:6.0-alpine AS build
WORKDIR /src
COPY ["Snippitbox.csproj", "Program.cs", "./"]

RUN dotnet restore -r linux-musl-x64
RUN dotnet publish -c release -o /dist -r linux-musl-x64 --self-contained false --no-restore

FROM mcr.microsoft.com/dotnet/aspnet:6.0-alpine AS runtime
WORKDIR /app
COPY --from=build /dist .
COPY ./wwwroot ./wwwroot/

VOLUME /app/data

ENV ASPNETCORE_URLS http://+:80
EXPOSE 80

ENTRYPOINT ["./Snippitbox"]