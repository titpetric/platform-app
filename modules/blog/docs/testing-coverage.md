# Testing Coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 70%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                                      | Coverage | Cognitive | Lines |
|--------|----------------------------------------------|----------|-----------|-------|
| ❌     | titpetric/platform-example/blog              | 11.78%   | 133       | 745   |
| ❌     | titpetric/platform-example/blog/cmd/blog     | 0.00%    | 2         | 22    |
| ❌     | titpetric/platform-example/blog/cmd/generate | 0.00%    | 5         | 51    |
| ✅     | titpetric/platform-example/blog/layout       | 96.97%   | 6         | 48    |
| ❌     | titpetric/platform-example/blog/markdown     | 72.00%   | 10        | 141   |
| ❌     | titpetric/platform-example/blog/model        | 24.00%   | 35        | 213   |
| ✅     | titpetric/platform-example/blog/schema       | 0.00%    | 0         | 0     |
| ✅     | titpetric/platform-example/blog/storage      | 90.06%   | 9         | 100   |
| ❌     | titpetric/platform-example/blog/view         | 0.00%    | 36        | 243   |

## Functions

| Status | Package                                      | Function                      | Coverage | Cognitive |
|--------|----------------------------------------------|-------------------------------|----------|-----------|
| ❌     | titpetric/platform-example/blog              | Generator.Generate            | 0.00%    | 12        |
| ❌     | titpetric/platform-example/blog              | Generator.copyAssets          | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog              | Generator.copyEmbeddedAssets  | 0.00%    | 8         |
| ❌     | titpetric/platform-example/blog              | Generator.copyLocalAssets     | 0.00%    | 9         |
| ❌     | titpetric/platform-example/blog              | Generator.generateArticlePage | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog              | Generator.generateFeed        | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog              | Generator.generateIndexPage   | 0.00%    | 2         |
| ✅     | titpetric/platform-example/blog              | Generator.generateStaticPages | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog              | Generator.walkPages           | 0.00%    | 37        |
| ❌     | titpetric/platform-example/blog              | Handlers.GetArticleHTML       | 0.00%    | 3         |
| ❌     | titpetric/platform-example/blog              | Handlers.GetArticleJSON       | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog              | Handlers.GetAtomFeed          | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog              | Handlers.IndexHTML            | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog              | Handlers.ListArticlesHTML     | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog              | Handlers.ListArticlesJSON     | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog              | Handlers.SearchArticlesJSON   | 0.00%    | 3         |
| ❌     | titpetric/platform-example/blog              | Module.Mount                  | 0.00%    | 1         |
| ✅     | titpetric/platform-example/blog              | Module.Name                   | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog              | Module.ScanMarkdownFiles      | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog              | Module.SetRepository          | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog              | Module.Start                  | 0.00%    | 4         |
| ✅     | titpetric/platform-example/blog              | Module.Stop                   | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog              | Module.parseMarkdownFile      | 0.00%    | 9         |
| ❌     | titpetric/platform-example/blog              | Module.scanMarkdownFiles      | 0.00%    | 9         |
| ✅     | titpetric/platform-example/blog              | NewGenerator                  | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog              | NewHandlers                   | 0.00%    | 1         |
| ❌     | titpetric/platform-example/blog              | NewModule                     | 0.00%    | 1         |
| ✅     | titpetric/platform-example/blog              | NewOverlayFS                  | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog              | OverlayFS.Glob                | 100.00%  | 5         |
| ✅     | titpetric/platform-example/blog              | OverlayFS.Open                | 85.70%   | 4         |
| ✅     | titpetric/platform-example/blog              | OverlayFS.ReadDir             | 91.30%   | 9         |
| ✅     | titpetric/platform-example/blog              | generateID                    | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog/cmd/blog     | main                          | 0.00%    | 1         |
| ❌     | titpetric/platform-example/blog/cmd/blog     | start                         | 0.00%    | 1         |
| ❌     | titpetric/platform-example/blog/cmd/generate | generate                      | 0.00%    | 4         |
| ❌     | titpetric/platform-example/blog/cmd/generate | main                          | 0.00%    | 1         |
| ✅     | titpetric/platform-example/blog/layout       | NewRenderer                   | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/layout       | Renderer.Render               | 90.90%   | 6         |
| ✅     | titpetric/platform-example/blog/layout       | Renderer.template             | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/markdown     | NewRenderer                   | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/markdown     | Renderer.Render               | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/markdown     | Renderer.highlightCodeBlocks  | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/markdown     | Renderer.processCodeBlock     | 87.50%   | 1         |
| ✅     | titpetric/platform-example/blog/markdown     | escapeHTML                    | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/markdown     | highlightCode                 | 88.50%   | 8         |
| ✅     | titpetric/platform-example/blog/markdown     | unescapeHTML                  | 100.00%  | 0         |
| ❌     | titpetric/platform-example/blog/markdown     | wrapCodePlain                 | 0.00%    | 1         |
| ❌     | titpetric/platform-example/blog/model        | Article.Delete                | 0.00%    | 1         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetCreatedAt          | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetDate               | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetDescription        | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetFilename           | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetID                 | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetLayout             | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetOgImage            | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetSlug               | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetSource             | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetTitle              | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetURL                | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.GetUpdatedAt          | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.Insert                | 80.00%   | 1         |
| ✅     | titpetric/platform-example/blog/model        | Article.Select                | 91.70%   | 4         |
| ✅     | titpetric/platform-example/blog/model        | Article.SetCreatedAt          | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.SetDate               | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/model        | Article.SetUpdatedAt          | 100.00%  | 0         |
| ❌     | titpetric/platform-example/blog/model        | Article.Update                | 0.00%    | 5         |
| ❌     | titpetric/platform-example/blog/model        | Migrations.Delete             | 0.00%    | 1         |
| ✅     | titpetric/platform-example/blog/model        | Migrations.GetFilename        | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Migrations.GetProject         | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Migrations.GetStatementIndex  | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | Migrations.GetStatus          | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog/model        | Migrations.Insert             | 0.00%    | 1         |
| ❌     | titpetric/platform-example/blog/model        | Migrations.Select             | 0.00%    | 4         |
| ❌     | titpetric/platform-example/blog/model        | Migrations.Update             | 0.00%    | 5         |
| ✅     | titpetric/platform-example/blog/model        | QueryConfig.Apply             | 88.20%   | 13        |
| ✅     | titpetric/platform-example/blog/model        | QueryConfig.WithColumns       | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | QueryConfig.WithLimit         | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | QueryConfig.WithOrderBy       | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | QueryConfig.WithStatement     | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | QueryConfig.WithTable         | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | QueryConfig.WithWhere         | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | WithColumns                   | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | WithLimit                     | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/model        | WithOrderBy                   | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/model        | WithStatement                 | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/model        | WithTable                     | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/model        | WithWhere                     | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/storage      | CountArticles                 | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/storage      | DB                            | 0.00%    | 0         |
| ✅     | titpetric/platform-example/blog/storage      | GetArticleBySlug              | 100.00%  | 1         |
| ✅     | titpetric/platform-example/blog/storage      | GetArticles                   | 83.30%   | 1         |
| ✅     | titpetric/platform-example/blog/storage      | InsertArticle                 | 100.00%  | 1         |
| ✅     | titpetric/platform-example/blog/storage      | NewStorage                    | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/storage      | SearchArticles                | 87.50%   | 1         |
| ✅     | titpetric/platform-example/blog/storage      | Storage.CountArticles         | 100.00%  | 1         |
| ✅     | titpetric/platform-example/blog/storage      | Storage.GetArticleBySlug      | 100.00%  | 1         |
| ✅     | titpetric/platform-example/blog/storage      | Storage.GetArticles           | 100.00%  | 1         |
| ✅     | titpetric/platform-example/blog/storage      | Storage.InitSchema            | 100.00%  | 0         |
| ✅     | titpetric/platform-example/blog/storage      | Storage.InsertArticle         | 100.00%  | 1         |
| ✅     | titpetric/platform-example/blog/storage      | Storage.SearchArticles        | 100.00%  | 1         |
| ✅     | titpetric/platform-example/blog/view         | IndexData.Map                 | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog/view         | NewViews                      | 0.00%    | 1         |
| ✅     | titpetric/platform-example/blog/view         | PostData.Map                  | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog/view         | StripFrontMatter              | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog/view         | Views.AtomFeed                | 0.00%    | 15        |
| ❌     | titpetric/platform-example/blog/view         | Views.Blog                    | 0.00%    | 4         |
| ❌     | titpetric/platform-example/blog/view         | Views.Index                   | 0.00%    | 3         |
| ✅     | titpetric/platform-example/blog/view         | Views.IndexFromArticles       | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog/view         | Views.Post                    | 0.00%    | 3         |
| ✅     | titpetric/platform-example/blog/view         | Views.PostFromArticle         | 0.00%    | 0         |
| ❌     | titpetric/platform-example/blog/view         | escapeXML                     | 0.00%    | 1         |
| ❌     | titpetric/platform-example/blog/view         | fillTemplateData              | 0.00%    | 3         |
| ❌     | titpetric/platform-example/blog/view         | loadFile                      | 0.00%    | 2         |
| ❌     | titpetric/platform-example/blog/view         | loadFileYaml                  | 0.00%    | 2         |
