# https://hub.docker.com/_/microsoft-dotnet
FROM mcr.microsoft.com/dotnet/sdk:5.0-alpine AS build

WORKDIR /source

# copy csproj and restore as distinct layers
COPY *.csproj .
RUN dotnet restore .

# copy everything else and build app
COPY . .
RUN dotnet publish -c release -o /app .

# final stage/image
FROM mcr.microsoft.com/dotnet/aspnet:5.0-alpine
WORKDIR /app
COPY --from=build /app ./
# Add a new user "radix-non-root-user" with user id 1001
RUN adduser -D --uid 1001 radix-non-root-user
USER 1001

ENTRYPOINT ["dotnet", "aspnet-webapp.dll"]