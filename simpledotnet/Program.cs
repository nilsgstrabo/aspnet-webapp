using System.Net;
using Microsoft.AspNetCore.HttpOverrides;

var builder = WebApplication.CreateBuilder(args);

builder.Services.Configure<ForwardedHeadersOptions>(options =>
{
	options.ForwardedHeaders = ForwardedHeaders.XForwardedFor | ForwardedHeaders.XForwardedProto | ForwardedHeaders.XForwardedHost;

	// Trust Istio proxies in Radix.
	options.KnownIPNetworks.Add(new System.Net.IPNetwork(IPAddress.Parse("10.0.0.0"), 8));
	
});

var app = builder.Build();

var allowedNetworks = new[]
{
	new System.Net.IPNetwork(IPAddress.Parse("203.0.113.10"), 32),
	new System.Net.IPNetwork(IPAddress.Parse("143.97.110.1"), 24),
};

app.UseForwardedHeaders();

// app.Use(async (context, next) =>
// {
// 	var remoteIp = context.Connection.RemoteIpAddress;
// 	var isAllowed = remoteIp is not null && allowedNetworks.Any(network => network.Contains(remoteIp));

// 	if (!isAllowed)
// 	{
// 		context.Response.StatusCode = StatusCodes.Status403Forbidden;
// 		await context.Response.WriteAsync("Forbidden");
// 		return;
// 	}

// 	await next();
// });

// app.MapGet("/", () => "Hello world");

app.MapGet("/", (HttpContext context) => {

	
	return Results.Ok(new 
	{
		ClientIP = context.Connection.RemoteIpAddress?.ToString(),
		context.Request.Protocol,
		context.Request.Scheme,
		context.Request.Host.Host

	});
});

app.MapGet("/headers", (HttpContext context) =>
{
	
	var headers = context.Request.Headers
		.Select(header => new
		{
			Header = header.Key,
			Values = header.Value.ToArray(),
		});

	return Results.Ok(headers);
});



app.Run();