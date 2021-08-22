using static HtmlBuilders.HtmlTags;

using LiteDB;

using Microsoft.AspNetCore.Html;
using HtmlBuilders;
using Microsoft.AspNetCore.Antiforgery;

var builder = WebApplication.CreateBuilder();
builder.Services
    .AddSingleton<SnippitRepository>()
    .AddAntiforgery()
    .AddMemoryCache();

builder.Logging.AddConsole().SetMinimumLevel(LogLevel.Warning);

var app = builder.Build();

app.MapGet("/", (HttpContext context, SnippitRepository repo, IAntiforgery antiForgery) =>
{
    var antiForgeryToken = antiForgery.GetAndStoreTokens(context);
    return Results.Text(RenderHome(antiForgeryToken), "text/html");
});

app.MapPost("/", async (HttpContext context, SnippitRepository repo, IAntiforgery antiForgery) => {
    await antiForgery.ValidateRequestAsync(context);
    var form = context.Request.Form;
    var snippit = new Snippit { Id = ObjectId.NewObjectId(), Description = form["Description"], CreatedAt = DateTimeOffset.UtcNow };
    repo.Save(snippit);
    return Results.Redirect("/");
});

app.Run();


static string RenderHome(AntiforgeryTokenSet antiForgeryToken)
{
    return Html
            .Attribute("lang", "en")
            .Append(
                Body.Append(H1.Append("Hi there"))
                    .Append(SnippitForm(antiForgeryToken))
            ).ToHtmlString();
}

static HtmlTag SnippitForm(AntiforgeryTokenSet antiForgeryToken)
{
    var descriptionField = Div
        .Append(Label.Append("Description"))
        .Append(Input.Text.Name("Description"));
    
    var submit = Div.Append(Button.Append("Submit"));

    var form = Form
           .Attribute("method", "post")
           .Attribute("action", "/")
             .Append(Input.Hidden.Name(antiForgeryToken.FormFieldName).Value(antiForgeryToken.RequestToken))
             .Append(descriptionField)
             .Append(submit);

    return form;
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

    public string DbPath => dbPath;

    public SummaryResult GetPaged(int start, int pageSize)
    {
        using var db = new LiteDatabase(dbPath);
        var col = db.GetCollection<Snippit>(nameof(Snippit));

        return new()
        {
            Total = col.Count(),
            Start = start,
            List = col.Find(Query.All(Query.Descending), skip: start, limit: pageSize)
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