using static Scriban.Template;

using LiteDB;
using Scriban;

using Microsoft.AspNetCore.Antiforgery;

var builder = WebApplication.CreateBuilder();
builder.Services
    .AddSingleton<SnippitRepository>()
    .AddAntiforgery()
    .AddMemoryCache();

builder.Logging.AddConsole().SetMinimumLevel(LogLevel.Warning);

var app = builder.Build();
app.UseStaticFiles();

app.MapGet("/", async (SnippitRepository repo) =>
{
    var result = repo.Find(0, 10);
    return Results.Text(await RenderHome(result), "text/html");
});

app.MapGet("/favicon.ico", () => Results.NotFound());

app.MapGet("/new", async (HttpContext context, SnippitRepository repo, IAntiforgery antiForgery) =>
{
    var antiForgeryToken = antiForgery.GetAndStoreTokens(context);
    return Results.Text(await RenderNew(antiForgeryToken), "text/html");
});

app.MapPost("/new", async (HttpContext context, SnippitRepository repo, IAntiforgery antiForgery) => {
    await antiForgery.ValidateRequestAsync(context);
    var form = context.Request.Form;
    var snippit = new Snippit { Id = ObjectId.NewObjectId(), Description = form["Description"], CreatedAt = DateTimeOffset.UtcNow };
    repo.Save(snippit);
    return Results.Redirect($"/{snippit.Id}");
});

app.MapGet("/{snippitId}", (HttpContext context, string snippitId, SnippitRepository repo) =>
{
    var snippit = repo.Get(new ObjectId(snippitId));
    return snippit;
});

app.Run();

static async Task<string> RenderHome(SummaryResult result)
{
    var contentTemplate = Parse(@"
    <div class=""an-box"">
        {{- for summary in result.list -}}
            {{ summary.description }}
        {{- end -}}
    </div>
    ");

    return await Layout().RenderAsync(new { Title = "Home", Content = await contentTemplate.RenderAsync(new { Result = result }) });
}

static async Task<string> RenderNew(AntiforgeryTokenSet antiForgeryToken)
{
    var contentTemplate = Parse(@"
        <form method=""post"">
            <input type=""hidden"" name=""{{token.form_field_name}}"" value=""{{token.request_token}}"">
            <div>
                <label for=""Description"">Description</label><input name=""Description"" type=""text"" />
            </div>
            <div>
                <button>Submit</button>
            </div>
        </form>
    ");

    return await Layout().RenderAsync(new { Content = await contentTemplate.RenderAsync(new { token = antiForgeryToken }) });
}

static Template Layout()
{
    return Parse(@"
    <!DOCTYPE html>
    <meta charset=""utf-8"">
    <meta name=""viewport"" content=""width=device-width, initial-scale=1"">
    <link href=""data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAA00lEQVQ4jZWSQRGEMAxFX6kBJFQCEpDAfS+c91QHSKgEJCChEioBCUhgD7RLtgsFMtML89/PTwI8qxpoAQ8EwD6B2wit4s134Q5ogCUzmO5EdlHsgUapH5P+Cp6yjl4kmaPmB2jRus/gEDt5YFUqmmg6CY5iPidir8AQdf1ZdCkOoF/8b9qyX8AdzWvibAal8lPJF/K5j8qewJ76Gh6iQW7i7nRuBJBMFuBNVVn2ZRa7y46w7SadszvGtqqBGfX9UeT3ZGpKBuZEkAx8MftFjRQW+AGoU1dleoYetgAAAABJRU5ErkJggg=="" rel=""icon"" type=""image/x-icon"" />
    <link rel=""preload"" href=""/css/anole.css"" as=""style"">
    <link href=""/css/anole.css"" rel=""stylesheet"" type=""text/css"" />
    <title>{{ title }}</title>
    <body class=""anole"">
        <div class=""an-body"">
            {{~ content ~}}
        </div>
    </body>
    ");
}

record Snippit
{
    public ObjectId Id { get; set; } = ObjectId.Empty;
    public string Description { get; set; } = string.Empty;
    public DateTimeOffset CreatedAt { get; set; }
    public File[] Files { get; set; } = Array.Empty<File>();
}

record File(string Name, string Contents, string FileType);

record Summary (string Id, string Description, DateTimeOffset CreatedAt);

record SummaryResult
{
    public int Total { get; set; }
    public int Start { get; set; }
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

    public SummaryResult Find(int start, int pageSize)
    {
        using var db = new LiteDatabase(dbPath);
        var col = db.GetCollection<Snippit>(nameof(Snippit));

        return new()
        {
            Total = col.Count(),
            Start = start,
            List = col.Find(Query.All(Query.Descending), skip: start, limit: pageSize)
                    .Select(s => new Summary(s.Id.ToString(),s.Description, s.CreatedAt))
                    .ToArray()
        };
    }

    public Snippit Get(ObjectId id)
    {
        using var db = new LiteDatabase(dbPath);
        var col = db.GetCollection<Snippit>(nameof(Snippit));
        return col.FindById(id);
    }

    public void Save(Snippit snippit)
    {
        using var db = new LiteDatabase(dbPath);
        var col = db.GetCollection<Snippit>(nameof(Snippit));
        col.EnsureIndex(s => s.CreatedAt);
        col.Upsert(snippit);
    }
}
