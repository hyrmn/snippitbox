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
    if (result.Total == 0) return Results.LocalRedirect("/new");
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
    var snippit = new Snippit(ObjectId.NewObjectId(), form["Description"], form["Contents"], DateTimeOffset.UtcNow);
    repo.Save(snippit);
    return Results.Redirect($"/{snippit.Id}");
});

app.MapGet("/{snippitId}", async (string snippitId, SnippitRepository repo) =>
{
    var snippit = repo.Get(new ObjectId(snippitId));
    if(snippit is not null)
        return Results.Text(await RenderSnippit(snippit), "text/html");
    return Results.LocalRedirect("/new");
});

app.Run();

static async Task<string> RenderHome(SummaryResult result)
{
    var contentTemplate = Parse(@"
    {{- for summary in result.list -}}
        <div class=""an-glass"">
            <div class=""an-group"">
                <span class=""an-85p"">{{ summary.description }}</span>
                <a href=""/{{ summary.id }}"" class=""an-button an-15p"">View</a>
            </div>
        </div>
    {{- end -}}
    ");

    return await Layout().RenderAsync(new { Title = "Snippits", Content = await contentTemplate.RenderAsync(new { Result = result }) });
}

static async Task<string> RenderSnippit(Snippit snippit)
{
    var contentTemplate = Parse(@"
    <div class=""an-glass"">
        <h2 class=""an-header"">{{ snippit.description }}</h2>
        <div class=""an-box"">
            {{ snippit.contents }}
        </div>
    </div>
    ");

    return await Layout().RenderAsync(new { Title = "Snippits", Content = await contentTemplate.RenderAsync(new { Snippit = snippit }) });
}

static async Task<string> RenderNew(AntiforgeryTokenSet antiForgeryToken)
{
    var contentTemplate = Parse(@"
    <div class=""an-box an-form"">
        <form method=""post"">
            <input type=""hidden"" name=""{{token.form_field_name}}"" value=""{{token.request_token}}"">
            <div>
                <label for=""Description"">Description</label><input name=""Description"" type=""text"" />
            </div>
            <div class=""an-gen-spacer""></div>
            <div>
                <label for=""Contents"">Snippit</label>
                <div class=""an-gen-spacer""></div>
                <textarea name=""Contents"" rows=""20"" cols=""80""></textarea>
            </div>
            <div>
                <input type=""submit"" value=""Save"" />
            </div>
        </form>
    </div>
    ");

    return await Layout().RenderAsync(new { Title="Add Snippit", Content = await contentTemplate.RenderAsync(new { token = antiForgeryToken }) });
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
    <link href=""/css/customizations.css"" rel=""stylesheet"" type=""text/css"" />
    <title>{{ title }}</title>
    <body class=""anole"">
        <header class=""an-header an-header-full"">{{ title }}</header>
        <nav class=""an-head an-menu"">
            <a href=""/"" class=""an-home-button"">Snippits</a>
            <a href=""/new"" class=""an-button"">New snippit</a>
        </nav>
        <section class=""an-body"">
            {{~ content ~}}
        </section>
    </body>
    ");
}

record Snippit(ObjectId Id, string Description, string Contents, DateTimeOffset CreatedAt);

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
        col.EnsureIndex(s => s.Description);
        col.Upsert(snippit);
    }
}
