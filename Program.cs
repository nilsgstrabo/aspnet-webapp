using System;
using Azure.Identity;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;

namespace aspnet_webapp
{
    public class Program
    {
        public static void Main(string[] args)
        {
            var builder = CreateHostBuilder(args);
            
            builder.Build().Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureAppConfiguration(config => 
                {
                    config.AddAzureKeyVault(
                        new Uri(Environment.GetEnvironmentVariable("KEY_VAULT_URL")), 
                        new DefaultAzureCredential() // Try different strategies to acquire credentials
                    );
                })
                .ConfigureWebHostDefaults(webBuilder =>
                {
                    webBuilder.UseStartup<Startup>();
                });
    }
}
