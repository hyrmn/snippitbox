using System.Collections.Generic;

using HtmlBuilders;

using LiteDB;

using Microsoft.AspNetCore.Html;

var builder = WebApplication.CreateBuilder();
builder.Services
    .AddSingleton<Render>()
    .AddSingleton<SnippitRepository>()
    .AddAntiforgery()
    .AddMemoryCache();

builder.Logging.AddConsole().SetMinimumLevel(LogLevel.Warning);

var app = builder.Build();

app.MapGet("/", (SnippitRepository repo, Render render) =>
{
    repo.Save(new Snippit { Id = ObjectId.NewObjectId(), Description = "Test", CreatedAt = DateTimeOffset.UtcNow });
    return repo.GetPaged();
    //return Results.Text(repo.DbPath, "text/html");
});

app.Run();

class Render
{

}

record Snippit
{
    public ObjectId Id { get; set; } = ObjectId.Empty;
    public string Description { get; set; } = string.Empty;
    public DateTimeOffset CreatedAt { get; set; }
}

record Summary
{
    public ObjectId Id { get; set; } = ObjectId.Empty;
    public string Description { get; set; } = string.Empty;
    public DateTimeOffset CreatedAt { get; set; }
}

record SummaryResult
{
    public int Total { get; set; }
    public int Page { get; set; }
    public Summary[] List { get; set; } = Array.Empty<Summary>();
}

class SnippitRepository
{
    private readonly string dbPath;
    private readonly ILogger<SnippitRepository> logger;

    public SnippitRepository(IWebHostEnvironment env, ILogger<SnippitRepository> logger)
    {
        this.dbPath = Path.Combine(env.ContentRootPath, "data/snippits.db");
        this.logger = logger;
    }

    public string DbPath => dbPath;

    public SummaryResult GetPaged()
    {
        using var db = new LiteDatabase(dbPath);
        var col = db.GetCollection<Snippit>(nameof(Snippit));
        col.EnsureIndex(s => s.CreatedAt);
        return new()
        {
            Total = col.Count(),
            List = col.Query()
                .OrderByDescending(s => s.CreatedAt)
                .Select(s => new Summary { Id = s.Id, Description = s.Description, CreatedAt = s.CreatedAt })
                .ToArray()
        };
    }

    public void Save(Snippit snippit)
    {
        using var db = new LiteDatabase(dbPath);
        var col = db.GetCollection<Snippit>(nameof(Snippit));
        col.EnsureIndex(s => s.CreatedAt);

        col.Upsert(snippit);
    }

    public Snippit Get(ObjectId id) => throw new NotImplementedException();
}