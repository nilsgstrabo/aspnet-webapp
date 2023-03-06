# https://hub.docker.com/_/microsoft-dotnet
FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build

WORKDIR /source

# copy csproj and restore as distinct layers
COPY *.csproj .
RUN dotnet restore .

# copy everything else and build app
COPY . .
RUN dotnet publish -c release -o /app .

# final stage/image
FROM mcr.microsoft.com/dotnet/aspnet:6.0
WORKDIR /app


RUN apt-get update && apt-get install -y curl gnupg


# install ODBC 18
RUN curl https://packages.microsoft.com/keys/microsoft.asc | apt-key add -

#Debian 11
RUN curl https://packages.microsoft.com/config/debian/11/prod.list > /etc/apt/sources.list.d/mssql-release.list

RUN apt-get update
RUN ACCEPT_EULA=Y apt-get install -y msodbcsql18 mssql-tools18



RUN	apt-get update && apt-get -y install curl
RUN	curl -sL https://aka.ms/InstallAzureCLIDeb | bash
RUN az --version

COPY --from=build /app ./
# Add a new user "radix-non-root-user" with user id 1001
# RUN adduser -D --uid 1001 radix-non-root-user
RUN useradd -m --uid 1001 radix-non-root-user
# RUN chown -R 1001 /opt/az
USER 1001

ENTRYPOINT ["dotnet", "aspnet-webapp.dll"]