$targetSqlServerFQDN = "$(az sql server show -n sql-radix-cost-allocation-dev -g cost-allocation | jq -r .fullyQualifiedDomainName)"

$access_token = (Get-AzAccessToken -ResourceUrl https://database.windows.net).Token

Invoke-Sqlcmd -InputFile ${env:GITHUB_WORKSPACE}/.github/workflows/scripts/sql/init.sql -AccessToken $access_token -ServerInstance $targetSqlServerFQDN -Database sqldb-radix-cost-allocation