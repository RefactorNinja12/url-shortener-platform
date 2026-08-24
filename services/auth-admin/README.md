# auth-admin

C#/.NET Minimal API för autentisering. Utfärdar JWT:er (HMAC-SHA256) och har
en skyddad endpoint som exempel på hur andra tjänster kan verifiera dem.

Standalone i det här skedet — go-api kräver ännu ingen inloggning för att
skapa länkar (se "Öppna beslut" i projektplanen).

## Endpoints

- `POST /register` — `{ "email": "...", "password": "..." }` → skapar en användare
- `POST /login` — samma body → `{ "token": "<jwt>" }`
- `GET /me` — kräver `Authorization: Bearer <jwt>` → returnerar claims ur token
- `GET /healthz`

## Köra lokalt

```bash
dotnet run
```

Läser anslutningssträng och JWT-hemlighet från `appsettings.Development.json`
när `ASPNETCORE_ENVIRONMENT=Development`.

## Migrationer

EF Core-migrationer körs automatiskt vid uppstart (`db.Database.Migrate()`).
Ny migration efter en modelländring:

```bash
dotnet ef migrations add <Namn> -o Migrations
```
