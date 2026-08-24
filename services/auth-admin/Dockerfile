# syntax=docker/dockerfile:1

FROM mcr.microsoft.com/dotnet/sdk:9.0-alpine AS build
WORKDIR /src

COPY AuthAdmin.csproj .
RUN dotnet restore

COPY . .
RUN dotnet publish -c Release -o /out --no-restore

FROM mcr.microsoft.com/dotnet/aspnet:9.0-alpine
WORKDIR /app
RUN adduser -D -u 10001 appuser
COPY --from=build /out .

USER appuser
EXPOSE 8080
ENTRYPOINT ["dotnet", "AuthAdmin.dll"]
