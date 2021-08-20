using HtmlBuilders;
using Microsoft.AspNetCore.Html;
var builder = WebApplication.CreateBuilder();
builder.Services
    .AddSingleton<Render>()
    .AddAntiforgery()
    .AddMemoryCache();

builder.Logging.AddConsole().SetMinimumLevel(LogLevel.Warning);

var app = builder.Build();

app.MapGet("/", (Render render) => {
     return Results.Text("Hello World!", "text/html");
});

app.Run();

class Render
{

}