var builder = WebApplication.CreateBuilder();
builder.Services
  .AddAntiforgery()
  .AddMemoryCache();

builder.Logging.AddConsole().SetMinimumLevel(LogLevel.Warning);

var app = builder.Build();

app.MapGet("/", () => "Hello World!");

app.Run();