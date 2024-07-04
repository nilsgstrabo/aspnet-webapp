var builder = WebApplication.CreateBuilder(args);

builder.WebHost.ConfigureKestrel(k=>k.ListenAnyIP(8080));

var app = builder.Build();

app.MapGet("/", () => "Hello World!");

app.Run();
