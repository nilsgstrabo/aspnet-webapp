# https://hub.docker.com/_/microsoft-dotnet
FROM mcr.microsoft.com/dotnet/sdk:6.0-alpine AS build

ARG SECRET_1
ARG SECRET_2
ARG V
ARG W

RUN echo ${SECRET_1}
RUN A=$(echo ${SECRET_1} | base64 -d) && echo ${A}
RUN echo ${A}
RUN V=$(echo ${SECRET_1} | base64 -d) && echo ${V}
RUN echo ${V}
RUN W=$(echo ${SECRET_1} | base64 -d | sha256sum) && echo ${W}
RUN echo ${W}

WORKDIR /source

# copy csproj and restore as distinct layers
COPY *.csproj .
RUN dotnet restore .

# copy everything else and build app
COPY . .
RUN dotnet publish -c release -o /app .

# final stage/image
FROM mcr.microsoft.com/dotnet/aspnet:6.0-alpine
WORKDIR /app
COPY --from=build /app ./
# Add a new user "radix-non-root-user" with user id 1001
RUN adduser -D --uid 1001 radix-non-root-user
USER 1001

ENTRYPOINT ["dotnet", "aspnet-webapp.dll"]