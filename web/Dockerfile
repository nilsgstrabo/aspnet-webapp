# https://hub.docker.com/_/microsoft-dotnet
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
# FROM rihagtest.azurecr.io/dotnet/sdk:6.0 as build

WORKDIR /source

# copy csproj and restore as distinct layers
COPY *.csproj .
RUN dotnet restore .

RUN --mount=type=secret,id=NUGET_EQUINOR_SECRET export BUILD_ARG=$(cat /run/secrets/NUGET_EQUINOR_SECRET) && \
    #instead of `echo...|wc` - use real command with $BUILD_ARG argument
    echo $BUILD_ARG

# copy everything else and build app
COPY . .
RUN dotnet publish -c release -o /app .

# final stage/image
FROM mcr.microsoft.com/dotnet/aspnet:8.0-noble-chiseled
# ENV DOTNET_EnableDiagnostics=0
WORKDIR /app

# RUN	apt-get update && apt-get -y install curl
# RUN	curl -sL https://aka.ms/InstallAzureCLIDeb | bash
# RUN az --version

COPY --from=build /app ./
COPY run.sh ./run.sh
# Add a new user "radix-non-root-user" with user id 1001
# RUN adduser -D --uid 1001 radix-non-root-user

# RUN useradd -m --uid 1001 radix-non-root-user

# RUN chown -R 1001 /opt/az
USER 1001

# CMD ["sh", "run.sh"]
ENTRYPOINT ["dotnet", "aspnet-webapp.dll"]
