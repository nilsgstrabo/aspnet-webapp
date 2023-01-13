using System;
using System.Linq;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using aspnet_webapp.Services;
using Microsoft.Extensions.Azure;
using Azure.Identity;

namespace aspnet_webapp
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
                .AddJwtBearer(o=>{
                    o.Audience="5e48ca1f-a2bf-4dec-b96d-bbf8ce69f9f6";
                    o.ClaimsIssuer="https://login.microsoftonline.com/3aa4a235-b6e2-48d5-9195-7fcf05b459b0/v2.0";
                    o.MetadataAddress="https://login.microsoftonline.com/3aa4a235-b6e2-48d5-9195-7fcf05b459b0/v2.0/.well-known/openid-configuration";
                    o.TokenValidationParameters=new Microsoft.IdentityModel.Tokens.TokenValidationParameters{
                        NameClaimType="name"
                    };
                });
            
            services.AddAuthorization();
            services.AddRazorPages(c=>{
                c.Conventions.AllowAnonymousToPage("/Index");
                c.Conventions.AllowAnonymousToPage("/Error");
            });

            services.AddScoped<IUserInfoService, UserInfoService>();
            
            services.AddAzureClients(builder=> {
                builder.AddSecretClient(new Uri(Configuration["KEY_VAULT_URL"]));
                builder.UseCredential(new ManagedIdentityCredential());
            });
        }
           
       // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                app.UseExceptionHandler("/Error");
                // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
            }

            app.UseStaticFiles();
            app.UseAuthentication();

            app.Use(async (ctx, next) => {
                var logfactory = app.ApplicationServices.GetService<ILoggerFactory>();
                var logger=logfactory.CreateLogger("middleware");
                foreach (var h in ctx.Request.Headers.Where(h=>h.Key.StartsWith("X-Auth") || h.Key.ToLower().StartsWith("auth")).AsEnumerable()) //.Where(h=>h.Key.StartsWith("X-Custom") || h.Key.ToLower().StartsWith("auth"))
                {
                    var value = h.Value.FirstOrDefault() ?? "";
                    logger.LogInformation("{0}:{1}", h.Key, value.Substring(0,value.Length>20 ? 20 : value.Length));
                }
                await next();
            });

            app.UseRouting();
            app.UseAuthorization();
            app.UseEndpoints(endpoints =>
            {
                endpoints.MapRazorPages().RequireAuthorization();
            });
        }
    }
}
