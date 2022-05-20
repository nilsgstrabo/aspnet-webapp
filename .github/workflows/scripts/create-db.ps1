$targetSqlServerFQDN = "$(az sql server show -n sql-radix-cost-allocation-dev -g cost-allocation | jq -r .fullyQualifiedDomainName)"
Invoke-Sqlcmd -InputFile ${env:GITHUB_WORKSPACE}/.github/workflows/scripts/sql/init.sql -ServerInstance $targetSqlServerFQDN -Database sqldb-radix-cost-allocation