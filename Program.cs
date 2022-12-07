using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Azure.Identity;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace aspnet_webapp
{
    public class Program
    {
        public static void Main(string[] args)
        {
            CreateHostBuilder(args).Build().Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureAppConfiguration(c=> 
                {
                    c.AddAzureKeyVault(new Uri(Environment.GetEnvironmentVariable("KEY_VAULT_URL")), new DefaultAzureCredential());
                }).ConfigureWebHostDefaults(webBuilder =>
                {
                    webBuilder.UseStartup<Startup>();
                });
    }
}
