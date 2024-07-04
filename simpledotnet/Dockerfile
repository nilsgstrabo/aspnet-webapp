FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine3.20 AS build
WORKDIR /source
COPY *.csproj .
RUN dotnet restore .
COPY . .
RUN dotnet publish -c release -o /app .

# final stage/image
FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine3.20
WORKDIR /app
COPY --from=build /app ./
USER 1001
CMD ["dotnet", "app.dll"]