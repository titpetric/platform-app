# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be covered by invocation.

The storage packages have integration tests.

## Packages

| Status | Package                    | Coverage | Cognitive | Lines |
|--------|----------------------------|----------|-----------|-------|
| ✅    | .                          | 0.00%    | 0         | 0     |
| ✅    | blog                       | 0.00%    | 0         | 3     |
| ❌    | blog/cmd/blog              | 0.00%    | 7         | 37    |
| ❌    | blog/cmd/blog/server       | 0.00%    | 1         | 26    |
| ❌    | blog/cmd/blog/version      | 0.00%    | 31        | 79    |
| ❌    | blog/cmd/generate          | 0.00%    | 5         | 50    |
| ✅    | blog/config                | 0.00%    | 0         | 3     |
| ✅    | blog/markdown              | 89.02%   | 10        | 131   |
| ❌    | blog/model                 | 37.55%   | 52        | 330   |
| ❌    | blog/schema                | 0.00%    | 5         | 19    |
| ❌    | blog/service               | 0.00%    | 104       | 490   |
| ❌    | blog/service/admin         | 43.35%   | 144       | 900   |
| ✅    | blog/service/api           | 83.01%   | 31        | 228   |
| ❌    | blog/service/web           | 14.69%   | 34        | 306   |
| ❌    | blog/storage               | 77.23%   | 56        | 461   |
| ❌    | blog/view                  | 37.69%   | 38        | 414   |
| ❌    | daily                      | 0.00%    | 16        | 76    |
| ❌    | daily/model                | 0.00%    | 35        | 197   |
| ✅    | daily/schema               | 0.00%    | 0         | 0     |
| ❌    | daily/storage              | 0.00%    | 29        | 171   |
| ✅    | daily/view                 | 50.00%   | 0         | 14    |
| ❌    | email                      | 0.00%    | 41        | 246   |
| ❌    | email/model                | 0.00%    | 57        | 338   |
| ✅    | email/schema               | 0.00%    | 0         | 3     |
| ❌    | email/smtp                 | 0.00%    | 7         | 46    |
| ❌    | email/storage              | 0.00%    | 15        | 122   |
| ❌    | internal/cmd/structurelint | 0.00%    | 57        | 169   |
| ✅    | maillist                   | 0.00%    | 0         | 23    |
| ✅    | maillist/model             | 0.00%    | 0         | 40    |
| ❌    | maillist/storage           | 0.00%    | 5         | 40    |
| ✅    | pulse                      | 0.00%    | 0         | 3     |
| ❌    | pulse/client               | 34.29%   | 30        | 163   |
| ❌    | pulse/cmd/pulse            | 0.00%    | 6         | 36    |
| ❌    | pulse/cmd/pulse/login      | 0.00%    | 13        | 96    |
| ❌    | pulse/cmd/pulse/record     | 0.00%    | 21        | 94    |
| ❌    | pulse/cmd/pulse/register   | 0.00%    | 18        | 128   |
| ❌    | pulse/cmd/pulse/server     | 0.00%    | 1         | 26    |
| ❌    | pulse/cmd/pulse/version    | 0.00%    | 31        | 79    |
| ✅    | pulse/config               | 0.00%    | 0         | 3     |
| ❌    | pulse/model                | 0.00%    | 57        | 311   |
| ✅    | pulse/schema               | 100.00%  | 0         | 3     |
| ❌    | pulse/service              | 0.00%    | 59        | 343   |
| ❌    | pulse/service/keycounter   | 0.00%    | 29        | 100   |
| ❌    | pulse/storage              | 26.32%   | 15        | 120   |
| ✅    | pulse/view                 | 0.00%    | 0         | 3     |
| ❌    | user                       | 1.79%    | 39        | 201   |
| ❌    | user/cmd/generate          | 0.00%    | 25        | 141   |
| ❌    | user/model                 | 8.80%    | 124       | 715   |
| ✅    | user/schema                | 0.00%    | 0         | 3     |
| ❌    | user/service               | 0.00%    | 7         | 76    |
| ❌    | user/service/api           | 28.38%   | 64        | 443   |
| ✅    | user/service/auth          | 87.93%   | 17        | 86    |
| ❌    | user/service/passkey       | 52.83%   | 16        | 147   |
| ❌    | user/service/web           | 27.78%   | 26        | 249   |
| ❌    | user/storage               | 0.62%    | 79        | 538   |
| ✅    | user/view                  | 100.00%  | 0         | 3     |

## Functions

| Status | Package                    | Function                                    | Coverage | Cognitive |
|--------|----------------------------|---------------------------------------------|----------|-----------|
| ✅    | blog                       | NewModule                                   | 0.00%    | 0         |
| ❌    | blog/cmd/blog              | main                                        | 0.00%    | 1         |
| ❌    |                            | run                                         | 0.00%    | 6         |
| ✅    | blog/cmd/blog/server       | NewCommand                                  | 0.00%    | 0         |
| ❌    |                            | Run                                         | 0.00%    | 1         |
| ✅    | blog/cmd/blog/version      | NewCommand                                  | 0.00%    | 0         |
| ❌    |                            | Run                                         | 0.00%    | 31        |
| ❌    | blog/cmd/generate          | generate                                    | 0.00%    | 4         |
| ❌    |                            | main                                        | 0.00%    | 1         |
| ✅    | blog/config                | ConfigFS                                    | 0.00%    | 0         |
| ✅    | blog/markdown              | NewRenderer                                 | 100.00%  | 0         |
| ✅    |                            | Renderer.Render                             | 100.00%  | 0         |
| ✅    |                            | Renderer.highlightCodeBlocks                | 100.00%  | 0         |
| ✅    |                            | Renderer.processCodeBlock                   | 95.65%   | 1         |
| ✅    |                            | escapeHTML                                  | 0.00%    | 0         |
| ✅    |                            | highlightCode                               | 91.89%   | 8         |
| ✅    |                            | unescapeHTML                                | 100.00%  | 0         |
| ❌    |                            | wrapCodePlain                               | 0.00%    | 1         |
| ✅    | blog/model                 | Article.Delete                              | 100.00%  | 1         |
| ✅    |                            | Article.GetCreatedAt                        | 0.00%    | 0         |
| ✅    |                            | Article.GetDate                             | 0.00%    | 0         |
| ✅    |                            | Article.GetDescription                      | 0.00%    | 0         |
| ✅    |                            | Article.GetDraft                            | 0.00%    | 0         |
| ✅    |                            | Article.GetFilename                         | 0.00%    | 0         |
| ✅    |                            | Article.GetID                               | 0.00%    | 0         |
| ✅    |                            | Article.GetLayout                           | 0.00%    | 0         |
| ✅    |                            | Article.GetOgImage                          | 0.00%    | 0         |
| ✅    |                            | Article.GetSlug                             | 0.00%    | 0         |
| ✅    |                            | Article.GetSource                           | 0.00%    | 0         |
| ✅    |                            | Article.GetTitle                            | 0.00%    | 0         |
| ✅    |                            | Article.GetURL                              | 0.00%    | 0         |
| ✅    |                            | Article.GetUpdatedAt                        | 0.00%    | 0         |
| ✅    |                            | Article.Insert                              | 80.00%   | 1         |
| ✅    |                            | Article.IsDraft                             | 100.00%  | 0         |
| ✅    |                            | Article.IsPublished                         | 100.00%  | 2         |
| ✅    |                            | Article.IsScheduled                         | 100.00%  | 2         |
| ✅    |                            | Article.Select                              | 91.67%   | 4         |
| ✅    |                            | Article.SetCreatedAt                        | 100.00%  | 0         |
| ✅    |                            | Article.SetDate                             | 100.00%  | 0         |
| ✅    |                            | Article.SetDescription                      | 0.00%    | 0         |
| ✅    |                            | Article.SetDraft                            | 0.00%    | 0         |
| ✅    |                            | Article.SetFilename                         | 0.00%    | 0         |
| ✅    |                            | Article.SetID                               | 0.00%    | 0         |
| ✅    |                            | Article.SetLayout                           | 0.00%    | 0         |
| ✅    |                            | Article.SetOgImage                          | 0.00%    | 0         |
| ✅    |                            | Article.SetSlug                             | 0.00%    | 0         |
| ✅    |                            | Article.SetSource                           | 0.00%    | 0         |
| ✅    |                            | Article.SetTitle                            | 0.00%    | 0         |
| ✅    |                            | Article.SetURL                              | 0.00%    | 0         |
| ✅    |                            | Article.SetUpdatedAt                        | 100.00%  | 0         |
| ✅    |                            | Article.Status                              | 100.00%  | 2         |
| ✅    |                            | Article.Update                              | 92.31%   | 5         |
| ❌    |                            | Migrations.Delete                           | 0.00%    | 1         |
| ✅    |                            | Migrations.GetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.GetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Insert                           | 0.00%    | 1         |
| ❌    |                            | Migrations.Select                           | 0.00%    | 4         |
| ✅    |                            | Migrations.SetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.SetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.SetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.SetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Update                           | 0.00%    | 5         |
| ✅    |                            | QueryConfig.Apply                           | 88.24%   | 13        |
| ✅    |                            | QueryConfig.WithColumns                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithLimit                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithOrderBy                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithStatement                   | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithTable                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithWhere                       | 0.00%    | 0         |
| ❌    |                            | Setting.Delete                              | 0.00%    | 1         |
| ✅    |                            | Setting.GetAnalyticsID                      | 0.00%    | 0         |
| ✅    |                            | Setting.GetCreatedAt                        | 0.00%    | 0         |
| ✅    |                            | Setting.GetFeatureComments                  | 0.00%    | 0         |
| ✅    |                            | Setting.GetFeaturePingback                  | 0.00%    | 0         |
| ✅    |                            | Setting.GetFeatureRss                       | 0.00%    | 0         |
| ✅    |                            | Setting.GetFeatureWebmention                | 0.00%    | 0         |
| ✅    |                            | Setting.GetMetaAuthorName                   | 0.00%    | 0         |
| ✅    |                            | Setting.GetMetaHeadshot                     | 0.00%    | 0         |
| ✅    |                            | Setting.GetMetaLang                         | 0.00%    | 0         |
| ✅    |                            | Setting.GetMetaSubtitle                     | 0.00%    | 0         |
| ✅    |                            | Setting.GetMetaURL                          | 0.00%    | 0         |
| ✅    |                            | Setting.GetPostsPerPage                     | 0.00%    | 0         |
| ✅    |                            | Setting.GetSeoDefaultImage                  | 0.00%    | 0         |
| ✅    |                            | Setting.GetSeoTitleSuffix                   | 0.00%    | 0         |
| ✅    |                            | Setting.GetSocialGithub                     | 0.00%    | 0         |
| ✅    |                            | Setting.GetSocialLinkedin                   | 0.00%    | 0         |
| ✅    |                            | Setting.GetSocialTwitter                    | 0.00%    | 0         |
| ✅    |                            | Setting.GetUpdatedAt                        | 0.00%    | 0         |
| ✅    |                            | Setting.GetUserID                           | 0.00%    | 0         |
| ✅    |                            | Setting.Insert                              | 80.00%   | 1         |
| ✅    |                            | Setting.Select                              | 83.33%   | 4         |
| ✅    |                            | Setting.SetAnalyticsID                      | 0.00%    | 0         |
| ✅    |                            | Setting.SetCreatedAt                        | 100.00%  | 0         |
| ✅    |                            | Setting.SetFeatureComments                  | 0.00%    | 0         |
| ✅    |                            | Setting.SetFeaturePingback                  | 0.00%    | 0         |
| ✅    |                            | Setting.SetFeatureRss                       | 0.00%    | 0         |
| ✅    |                            | Setting.SetFeatureWebmention                | 0.00%    | 0         |
| ✅    |                            | Setting.SetMetaAuthorName                   | 0.00%    | 0         |
| ✅    |                            | Setting.SetMetaHeadshot                     | 0.00%    | 0         |
| ✅    |                            | Setting.SetMetaLang                         | 0.00%    | 0         |
| ✅    |                            | Setting.SetMetaSubtitle                     | 0.00%    | 0         |
| ✅    |                            | Setting.SetMetaURL                          | 0.00%    | 0         |
| ✅    |                            | Setting.SetPostsPerPage                     | 0.00%    | 0         |
| ✅    |                            | Setting.SetSeoDefaultImage                  | 0.00%    | 0         |
| ✅    |                            | Setting.SetSeoTitleSuffix                   | 0.00%    | 0         |
| ✅    |                            | Setting.SetSocialGithub                     | 0.00%    | 0         |
| ✅    |                            | Setting.SetSocialLinkedin                   | 0.00%    | 0         |
| ✅    |                            | Setting.SetSocialTwitter                    | 0.00%    | 0         |
| ✅    |                            | Setting.SetUpdatedAt                        | 100.00%  | 0         |
| ✅    |                            | Setting.SetUserID                           | 0.00%    | 0         |
| ❌    |                            | Setting.Update                              | 0.00%    | 5         |
| ✅    |                            | WithColumns                                 | 0.00%    | 0         |
| ✅    |                            | WithLimit                                   | 100.00%  | 0         |
| ✅    |                            | WithOrderBy                                 | 100.00%  | 0         |
| ✅    |                            | WithStatement                               | 100.00%  | 0         |
| ✅    |                            | WithTable                                   | 0.00%    | 0         |
| ✅    |                            | WithWhere                                   | 100.00%  | 0         |
| ❌    | blog/schema                | GetTable                                    | 0.00%    | 4         |
| ❌    |                            | LoadSchema                                  | 0.00%    | 1         |
| ❌    | blog/service               | AdminFS                                     | 0.00%    | 1         |
| ❌    |                            | BlogModule.Mount                            | 0.00%    | 1         |
| ✅    |                            | BlogModule.Name                             | 0.00%    | 0         |
| ✅    |                            | BlogModule.ScanMarkdownFiles                | 0.00%    | 0         |
| ✅    |                            | BlogModule.SetRepository                    | 0.00%    | 0         |
| ❌    |                            | BlogModule.Start                            | 0.00%    | 4         |
| ✅    |                            | BlogModule.Stop                             | 0.00%    | 0         |
| ❌    |                            | BlogModule.initHandlers                     | 0.00%    | 2         |
| ❌    |                            | BlogModule.parseMarkdownFile                | 0.00%    | 10        |
| ❌    |                            | BlogModule.scanMarkdownFiles                | 0.00%    | 11        |
| ✅    |                            | FS                                          | 0.00%    | 0         |
| ❌    |                            | Generator.Generate                          | 0.00%    | 11        |
| ❌    |                            | Generator.copyAssets                        | 0.00%    | 2         |
| ❌    |                            | Generator.copyEmbeddedAssets                | 0.00%    | 8         |
| ❌    |                            | Generator.copyLocalAssets                   | 0.00%    | 9         |
| ❌    |                            | Generator.generateArticlePage               | 0.00%    | 2         |
| ❌    |                            | Generator.generateFeed                      | 0.00%    | 2         |
| ❌    |                            | Generator.generateIndexPage                 | 0.00%    | 2         |
| ✅    |                            | Generator.generateStaticPages               | 0.00%    | 0         |
| ❌    |                            | Generator.walkPages                         | 0.00%    | 37        |
| ✅    |                            | NewBlogModule                               | 0.00%    | 0         |
| ❌    |                            | NewGenerator                                | 0.00%    | 2         |
| ✅    |                            | generateID                                  | 0.00%    | 0         |
| ✅    | blog/service/admin         | ArticleRequest.BuildMarkdownContent         | 80.65%   | 11        |
| ✅    |                            | ArticleRequest.ToArticle                    | 100.00%  | 5         |
| ❌    |                            | ArticleRequest.UpdateArticle                | 71.43%   | 6         |
| ✅    |                            | ArticleRequest.Validate                     | 90.00%   | 15        |
| ✅    |                            | ArticleRequest.formatDateTimeForFrontmatter | 80.00%   | 2         |
| ✅    |                            | ArticleRequest.parseCombinedDateTime        | 66.67%   | 4         |
| ✅    |                            | ErrBadRequest                               | 100.00%  | 0         |
| ✅    |                            | ErrConflict                                 | 100.00%  | 0         |
| ✅    |                            | ErrForbidden                                | 0.00%    | 0         |
| ✅    |                            | ErrInternal                                 | 0.00%    | 0         |
| ✅    |                            | ErrNotFound                                 | 100.00%  | 0         |
| ✅    |                            | ErrUnauthorized                             | 0.00%    | 0         |
| ✅    |                            | Error.Error                                 | 0.00%    | 0         |
| ✅    |                            | Error.Unwrap                                | 0.00%    | 0         |
| ✅    |                            | Handlers.CheckSlugJSON                      | 100.00%  | 0         |
| ✅    |                            | Handlers.CreateArticleJSON                  | 100.00%  | 0         |
| ✅    |                            | Handlers.DashboardHTML                      | 0.00%    | 0         |
| ✅    |                            | Handlers.DeleteArticleJSON                  | 100.00%  | 0         |
| ✅    |                            | Handlers.EditArticleHTML                    | 100.00%  | 0         |
| ✅    |                            | Handlers.GetArticleJSON                     | 0.00%    | 0         |
| ✅    |                            | Handlers.GetSettingsJSON                    | 0.00%    | 0         |
| ✅    |                            | Handlers.GetSettingsSchemaJSON              | 0.00%    | 0         |
| ✅    |                            | Handlers.ListDraftsHTML                     | 0.00%    | 0         |
| ✅    |                            | Handlers.ListDraftsJSON                     | 0.00%    | 0         |
| ✅    |                            | Handlers.ListPublishedHTML                  | 0.00%    | 0         |
| ✅    |                            | Handlers.ListPublishedJSON                  | 0.00%    | 0         |
| ✅    |                            | Handlers.ListScheduledHTML                  | 0.00%    | 0         |
| ✅    |                            | Handlers.ListScheduledJSON                  | 0.00%    | 0         |
| ✅    |                            | Handlers.Mount                              | 0.00%    | 0         |
| ✅    |                            | Handlers.NewArticleHTML                     | 0.00%    | 0         |
| ✅    |                            | Handlers.PublishArticleJSON                 | 100.00%  | 0         |
| ✅    |                            | Handlers.SaveSettingsJSON                   | 100.00%  | 0         |
| ✅    |                            | Handlers.SettingsHTML                       | 0.00%    | 0         |
| ✅    |                            | Handlers.UpdateArticleJSON                  | 100.00%  | 0         |
| ✅    |                            | Handlers.checkSlugJSON                      | 85.71%   | 2         |
| ✅    |                            | Handlers.createArticleJSON                  | 91.30%   | 5         |
| ❌    |                            | Handlers.dashboardHTML                      | 0.00%    | 7         |
| ❌    |                            | Handlers.deleteArticleJSON                  | 78.57%   | 9         |
| ✅    |                            | Handlers.editArticleHTML                    | 18.18%   | 4         |
| ✅    |                            | Handlers.errorHandler                       | 78.57%   | 5         |
| ❌    |                            | Handlers.getArticleJSON                     | 0.00%    | 2         |
| ❌    |                            | Handlers.getSettingsJSON                    | 0.00%    | 1         |
| ❌    |                            | Handlers.getSettingsSchemaJSON              | 0.00%    | 1         |
| ❌    |                            | Handlers.listDraftsHTML                     | 0.00%    | 3         |
| ❌    |                            | Handlers.listDraftsJSON                     | 0.00%    | 2         |
| ❌    |                            | Handlers.listPublishedHTML                  | 0.00%    | 3         |
| ❌    |                            | Handlers.listPublishedJSON                  | 0.00%    | 2         |
| ❌    |                            | Handlers.listScheduledHTML                  | 0.00%    | 3         |
| ❌    |                            | Handlers.listScheduledJSON                  | 0.00%    | 2         |
| ❌    |                            | Handlers.newArticleHTML                     | 0.00%    | 1         |
| ✅    |                            | Handlers.publishArticleJSON                 | 83.33%   | 9         |
| ✅    |                            | Handlers.saveSettingsJSON                   | 84.62%   | 3         |
| ❌    |                            | Handlers.settingsHTML                       | 0.00%    | 2         |
| ✅    |                            | Handlers.updateArticleJSON                  | 83.33%   | 6         |
| ✅    |                            | NewError                                    | 100.00%  | 0         |
| ✅    |                            | NewHandlers                                 | 0.00%    | 0         |
| ✅    |                            | escapeYAML                                  | 100.00%  | 0         |
| ✅    |                            | isValidSlugAdmin                            | 100.00%  | 2         |
| ✅    |                            | parseDateTime                               | 85.71%   | 3         |
| ✅    |                            | parsePagination                             | 100.00%  | 8         |
| ✅    |                            | removeDraftFromFrontmatter                  | 95.24%   | 5         |
| ❌    |                            | requireLoginRedirect                        | 0.00%    | 3         |
| ✅    |                            | sanitizeContent                             | 100.00%  | 0         |
| ✅    |                            | validateSettings                            | 86.21%   | 8         |
| ✅    |                            | writeJSON                                   | 100.00%  | 0         |
| ✅    | blog/service/api           | ErrBadRequest                               | 100.00%  | 0         |
| ✅    |                            | ErrInternal                                 | 0.00%    | 0         |
| ✅    |                            | ErrNotFound                                 | 100.00%  | 0         |
| ✅    |                            | Error.Error                                 | 100.00%  | 0         |
| ✅    |                            | Error.Unwrap                                | 0.00%    | 0         |
| ✅    |                            | Handlers.GetArticleAdminJSON                | 100.00%  | 0         |
| ✅    |                            | Handlers.GetArticleJSON                     | 100.00%  | 0         |
| ✅    |                            | Handlers.ListArticlesAdminJSON              | 100.00%  | 0         |
| ✅    |                            | Handlers.ListArticlesJSON                   | 100.00%  | 0         |
| ✅    |                            | Handlers.Mount                              | 0.00%    | 0         |
| ✅    |                            | Handlers.SearchArticlesJSON                 | 100.00%  | 0         |
| ✅    |                            | Handlers.errorHandler                       | 78.57%   | 5         |
| ✅    |                            | Handlers.getArticleAdminJSON                | 87.50%   | 4         |
| ✅    |                            | Handlers.getArticleJSON                     | 92.86%   | 3         |
| ✅    |                            | Handlers.listArticlesAdminJSON              | 93.18%   | 11        |
| ✅    |                            | Handlers.listArticlesJSON                   | 88.24%   | 2         |
| ✅    |                            | Handlers.searchArticlesJSON                 | 90.91%   | 4         |
| ✅    |                            | NewError                                    | 100.00%  | 0         |
| ✅    |                            | NewHandlers                                 | 100.00%  | 0         |
| ✅    |                            | isValidSlug                                 | 100.00%  | 2         |
| ✅    | blog/service/web           | ErrBadRequest                               | 0.00%    | 0         |
| ✅    |                            | ErrInternal                                 | 0.00%    | 0         |
| ✅    |                            | ErrNotFound                                 | 100.00%  | 0         |
| ✅    |                            | Error.Error                                 | 0.00%    | 0         |
| ✅    |                            | Error.Unwrap                                | 0.00%    | 0         |
| ✅    |                            | Handlers.GetArticleHTML                     | 100.00%  | 0         |
| ✅    |                            | Handlers.GetAtomFeed                        | 0.00%    | 0         |
| ✅    |                            | Handlers.IndexHTML                          | 0.00%    | 0         |
| ✅    |                            | Handlers.ListArticlesAdminHTML              | 0.00%    | 0         |
| ✅    |                            | Handlers.ListArticlesHTML                   | 0.00%    | 0         |
| ✅    |                            | Handlers.Mount                              | 0.00%    | 0         |
| ✅    |                            | Handlers.Repository                         | 0.00%    | 0         |
| ✅    |                            | Handlers.Views                              | 0.00%    | 0         |
| ✅    |                            | Handlers.errorHandler                       | 64.29%   | 5         |
| ❌    |                            | Handlers.getArticleHTML                     | 52.38%   | 10        |
| ❌    |                            | Handlers.getAtomFeed                        | 0.00%    | 2         |
| ❌    |                            | Handlers.indexHTML                          | 0.00%    | 2         |
| ❌    |                            | Handlers.listArticlesAdminHTML              | 0.00%    | 11        |
| ❌    |                            | Handlers.listArticlesHTML                   | 0.00%    | 2         |
| ✅    |                            | Handlers.registerAssets                     | 0.00%    | 0         |
| ✅    |                            | NewError                                    | 100.00%  | 0         |
| ✅    |                            | NewHandlers                                 | 0.00%    | 0         |
| ✅    |                            | isValidSlug                                 | 66.67%   | 2         |
| ✅    | blog/storage               | CountArticles                               | 100.00%  | 0         |
| ✅    |                            | CountDraftArticles                          | 100.00%  | 0         |
| ✅    |                            | CountPublishedArticles                      | 0.00%    | 0         |
| ✅    |                            | CountScheduledArticles                      | 0.00%    | 0         |
| ✅    |                            | DB                                          | 0.00%    | 0         |
| ✅    |                            | DeleteArticle                               | 100.00%  | 0         |
| ✅    |                            | GetArticleByID                              | 90.00%   | 1         |
| ✅    |                            | GetArticleBySlug                            | 100.00%  | 1         |
| ✅    |                            | GetArticles                                 | 92.86%   | 1         |
| ✅    |                            | GetDraftArticles                            | 90.00%   | 1         |
| ✅    |                            | GetGlobalSettings                           | 100.00%  | 0         |
| ✅    |                            | GetPublishedArticleBySlug                   | 100.00%  | 1         |
| ✅    |                            | GetPublishedArticles                        | 90.00%   | 1         |
| ❌    |                            | GetScheduledArticles                        | 0.00%    | 1         |
| ✅    |                            | GetSettingByUserID                          | 90.00%   | 1         |
| ❌    |                            | GitFS.CopyFile                              | 0.00%    | 3         |
| ✅    |                            | GitFS.Open                                  | 100.00%  | 0         |
| ✅    |                            | GitFS.ReadDir                               | 100.00%  | 0         |
| ✅    |                            | GitFS.ReadFile                              | 100.00%  | 0         |
| ✅    |                            | GitFS.Remove                                | 83.33%   | 1         |
| ✅    |                            | GitFS.Rename                                | 82.76%   | 5         |
| ✅    |                            | GitFS.Root                                  | 100.00%  | 0         |
| ✅    |                            | GitFS.Stat                                  | 100.00%  | 0         |
| ✅    |                            | GitFS.WriteFile                             | 80.00%   | 2         |
| ✅    |                            | GitFS.commit                                | 71.43%   | 2         |
| ✅    |                            | GitFS.commitFile                            | 66.67%   | 2         |
| ✅    |                            | GitFS.initNewRepo                           | 77.78%   | 2         |
| ✅    |                            | GitFS.initRepo                              | 42.86%   | 2         |
| ✅    |                            | GitFS.normalizePath                         | 60.00%   | 2         |
| ✅    |                            | InsertArticle                               | 100.00%  | 1         |
| ✅    |                            | Migrate                                     | 85.71%   | 2         |
| ❌    |                            | New                                         | 0.00%    | 1         |
| ✅    |                            | NewGitFS                                    | 63.64%   | 4         |
| ✅    |                            | NewStorage                                  | 66.67%   | 1         |
| ✅    |                            | SaveSetting                                 | 100.00%  | 0         |
| ✅    |                            | SearchArticles                              | 95.00%   | 1         |
| ✅    |                            | SearchPublishedArticles                     | 94.12%   | 1         |
| ✅    |                            | Storage.CountArticles                       | 100.00%  | 1         |
| ✅    |                            | Storage.CountDraftArticles                  | 100.00%  | 1         |
| ❌    |                            | Storage.CountPublishedArticles              | 0.00%    | 1         |
| ❌    |                            | Storage.CountScheduledArticles              | 0.00%    | 1         |
| ✅    |                            | Storage.DeleteArticle                       | 100.00%  | 1         |
| ✅    |                            | Storage.GetArticleByID                      | 100.00%  | 1         |
| ✅    |                            | Storage.GetArticleBySlug                    | 100.00%  | 1         |
| ✅    |                            | Storage.GetArticles                         | 100.00%  | 1         |
| ✅    |                            | Storage.GetDraftArticles                    | 100.00%  | 1         |
| ✅    |                            | Storage.GetGlobalSettings                   | 100.00%  | 0         |
| ✅    |                            | Storage.GetPublishedArticleBySlug           | 100.00%  | 1         |
| ✅    |                            | Storage.GetPublishedArticles                | 100.00%  | 1         |
| ❌    |                            | Storage.GetScheduledArticles                | 0.00%    | 1         |
| ✅    |                            | Storage.GetSettingByUserID                  | 0.00%    | 0         |
| ✅    |                            | Storage.InsertArticle                       | 100.00%  | 1         |
| ✅    |                            | Storage.SaveSetting                         | 100.00%  | 0         |
| ✅    |                            | Storage.SearchArticles                      | 100.00%  | 1         |
| ✅    |                            | Storage.SearchPublishedArticles             | 100.00%  | 1         |
| ✅    |                            | Storage.UpdateArticle                       | 100.00%  | 1         |
| ✅    |                            | UpdateArticle                               | 100.00%  | 0         |
| ✅    |                            | escapeLike                                  | 100.00%  | 0         |
| ✅    | blog/view                  | AdminEditData.Map                           | 84.62%   | 4         |
| ✅    |                            | AdminListData.Map                           | 0.00%    | 0         |
| ✅    |                            | AdminListData.TotalPages                    | 100.00%  | 1         |
| ✅    |                            | AdminNavigation.WithActive                  | 100.00%  | 2         |
| ✅    |                            | AdminViews.Dashboard                        | 0.00%    | 0         |
| ✅    |                            | AdminViews.Edit                             | 0.00%    | 0         |
| ✅    |                            | AdminViews.List                             | 0.00%    | 0         |
| ✅    |                            | DefaultFeedConfig                           | 0.00%    | 0         |
| ❌    |                            | ExtractCustomYAML                           | 0.00%    | 9         |
| ✅    |                            | IndexData.Map                               | 0.00%    | 0         |
| ✅    |                            | LoadAdminMenuConfig                         | 81.82%   | 2         |
| ✅    |                            | LoadMenuConfig                              | 90.91%   | 2         |
| ✅    |                            | Loader.Load                                 | 100.00%  | 0         |
| ✅    |                            | MenuData.Map                                | 100.00%  | 0         |
| ❌    |                            | MenuItem.Show                               | 0.00%    | 2         |
| ✅    |                            | NewAdminDashboardData                       | 100.00%  | 0         |
| ✅    |                            | NewAdminEditData                            | 100.00%  | 2         |
| ✅    |                            | NewAdminListData                            | 100.00%  | 0         |
| ✅    |                            | NewAdminNavigation                          | 0.00%    | 0         |
| ✅    |                            | NewAdminViews                               | 0.00%    | 0         |
| ✅    |                            | NewIndexData                                | 0.00%    | 0         |
| ✅    |                            | NewLoader                                   | 100.00%  | 0         |
| ✅    |                            | NewMenuData                                 | 100.00%  | 0         |
| ✅    |                            | NewPostData                                 | 0.00%    | 0         |
| ✅    |                            | NewViews                                    | 0.00%    | 0         |
| ✅    |                            | PostData.Map                                | 0.00%    | 0         |
| ❌    |                            | StripFrontMatter                            | 0.00%    | 2         |
| ✅    |                            | Templates                                   | 100.00%  | 0         |
| ✅    |                            | Views.AtomFeed                              | 0.00%    | 0         |
| ✅    |                            | Views.AtomFeedWithConfig                    | 0.00%    | 0         |
| ✅    |                            | Views.Blog                                  | 0.00%    | 0         |
| ✅    |                            | Views.Index                                 | 0.00%    | 0         |
| ✅    |                            | Views.Post                                  | 0.00%    | 0         |
| ❌    |                            | Views.atomFeed                              | 0.00%    | 12        |
| ✅    |                            | escapeXML                                   | 0.00%    | 0         |
| ❌    | daily                      | Module.Mount                                | 0.00%    | 13        |
| ✅    |                            | Module.Name                                 | 0.00%    | 0         |
| ❌    |                            | Module.Start                                | 0.00%    | 3         |
| ✅    |                            | NewModule                                   | 0.00%    | 0         |
| ❌    | daily/model                | Migrations.Delete                           | 0.00%    | 1         |
| ✅    |                            | Migrations.GetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.GetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Insert                           | 0.00%    | 1         |
| ❌    |                            | Migrations.Select                           | 0.00%    | 4         |
| ✅    |                            | Migrations.SetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.SetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.SetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.SetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Update                           | 0.00%    | 5         |
| ❌    |                            | QueryConfig.Apply                           | 0.00%    | 13        |
| ✅    |                            | QueryConfig.WithColumns                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithLimit                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithOrderBy                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithStatement                   | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithTable                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithWhere                       | 0.00%    | 0         |
| ❌    |                            | Todo.Delete                                 | 0.00%    | 1         |
| ✅    |                            | Todo.GetCompleted                           | 0.00%    | 0         |
| ✅    |                            | Todo.GetCreatedAt                           | 0.00%    | 0         |
| ✅    |                            | Todo.GetDeletedAt                           | 0.00%    | 0         |
| ✅    |                            | Todo.GetID                                  | 0.00%    | 0         |
| ✅    |                            | Todo.GetTitle                               | 0.00%    | 0         |
| ✅    |                            | Todo.GetUpdatedAt                           | 0.00%    | 0         |
| ✅    |                            | Todo.GetUserID                              | 0.00%    | 0         |
| ❌    |                            | Todo.Insert                                 | 0.00%    | 1         |
| ❌    |                            | Todo.Select                                 | 0.00%    | 4         |
| ✅    |                            | Todo.SetCompleted                           | 0.00%    | 0         |
| ✅    |                            | Todo.SetCreatedAt                           | 0.00%    | 0         |
| ✅    |                            | Todo.SetDeletedAt                           | 0.00%    | 0         |
| ✅    |                            | Todo.SetID                                  | 0.00%    | 0         |
| ✅    |                            | Todo.SetTitle                               | 0.00%    | 0         |
| ✅    |                            | Todo.SetUpdatedAt                           | 0.00%    | 0         |
| ✅    |                            | Todo.SetUserID                              | 0.00%    | 0         |
| ❌    |                            | Todo.Update                                 | 0.00%    | 5         |
| ✅    |                            | WithColumns                                 | 0.00%    | 0         |
| ✅    |                            | WithLimit                                   | 0.00%    | 0         |
| ✅    |                            | WithOrderBy                                 | 0.00%    | 0         |
| ✅    |                            | WithStatement                               | 0.00%    | 0         |
| ✅    |                            | WithTable                                   | 0.00%    | 0         |
| ✅    |                            | WithWhere                                   | 0.00%    | 0         |
| ✅    | daily/storage              | DB                                          | 0.00%    | 0         |
| ❌    |                            | Migrate                                     | 0.00%    | 2         |
| ❌    |                            | New                                         | 0.00%    | 1         |
| ❌    |                            | NewStorage                                  | 0.00%    | 1         |
| ❌    |                            | Storage.Add                                 | 0.00%    | 2         |
| ❌    |                            | Storage.Complete                            | 0.00%    | 6         |
| ❌    |                            | Storage.Delete                              | 0.00%    | 6         |
| ❌    |                            | Storage.Get                                 | 0.00%    | 2         |
| ❌    |                            | Storage.List                                | 0.00%    | 2         |
| ❌    |                            | Storage.Update                              | 0.00%    | 6         |
| ❌    |                            | boolToInt                                   | 0.00%    | 1         |
| ✅    | daily/view                 | Loader.Load                                 | 100.00%  | 0         |
| ✅    |                            | NewLoader                                   | 100.00%  | 0         |
| ✅    |                            | NewViews                                    | 0.00%    | 0         |
| ✅    |                            | Views.Index                                 | 0.00%    | 0         |
| ✅    | email                      | DefaultServiceOptions                       | 0.00%    | 0         |
| ✅    |                            | Handler.Mount                               | 0.00%    | 0         |
| ✅    |                            | Handler.Name                                | 0.00%    | 0         |
| ❌    |                            | Handler.Start                               | 0.00%    | 2         |
| ❌    |                            | Handler.Stop                                | 0.00%    | 1         |
| ✅    |                            | NewModule                                   | 0.00%    | 0         |
| ❌    |                            | NewService                                  | 0.00%    | 13        |
| ❌    |                            | Service.AddEmail                            | 0.00%    | 2         |
| ✅    |                            | Service.Start                               | 0.00%    | 0         |
| ✅    |                            | Service.Stop                                | 0.00%    | 0         |
| ❌    |                            | Service.processPendingEmails                | 0.00%    | 4         |
| ❌    |                            | Service.processQueue                        | 0.00%    | 6         |
| ❌    |                            | Service.sendEmail                           | 0.00%    | 13        |
| ❌    | email/model                | Email.Delete                                | 0.00%    | 1         |
| ✅    |                            | Email.GetBody                               | 0.00%    | 0         |
| ✅    |                            | Email.GetCreatedAt                          | 0.00%    | 0         |
| ✅    |                            | Email.GetError                              | 0.00%    | 0         |
| ✅    |                            | Email.GetID                                 | 0.00%    | 0         |
| ✅    |                            | Email.GetRecipient                          | 0.00%    | 0         |
| ✅    |                            | Email.GetRetryAt                            | 0.00%    | 0         |
| ✅    |                            | Email.GetRetryCount                         | 0.00%    | 0         |
| ✅    |                            | Email.GetRetryError                         | 0.00%    | 0         |
| ✅    |                            | Email.GetSentAt                             | 0.00%    | 0         |
| ✅    |                            | Email.GetStatus                             | 0.00%    | 0         |
| ✅    |                            | Email.GetSubject                            | 0.00%    | 0         |
| ❌    |                            | Email.Insert                                | 0.00%    | 1         |
| ❌    |                            | Email.Select                                | 0.00%    | 4         |
| ✅    |                            | Email.SetCreatedAt                          | 0.00%    | 0         |
| ✅    |                            | Email.SetRetryAt                            | 0.00%    | 0         |
| ✅    |                            | Email.SetSentAt                             | 0.00%    | 0         |
| ❌    |                            | Email.Update                                | 0.00%    | 5         |
| ❌    |                            | EmailFailed.Delete                          | 0.00%    | 1         |
| ✅    |                            | EmailFailed.GetBody                         | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetCreatedAt                    | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetError                        | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetID                           | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetRecipient                    | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetRetryAt                      | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetRetryCount                   | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetRetryError                   | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetSentAt                       | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetStatus                       | 0.00%    | 0         |
| ✅    |                            | EmailFailed.GetSubject                      | 0.00%    | 0         |
| ❌    |                            | EmailFailed.Insert                          | 0.00%    | 1         |
| ❌    |                            | EmailFailed.Select                          | 0.00%    | 4         |
| ✅    |                            | EmailFailed.SetCreatedAt                    | 0.00%    | 0         |
| ✅    |                            | EmailFailed.SetRetryAt                      | 0.00%    | 0         |
| ✅    |                            | EmailFailed.SetSentAt                       | 0.00%    | 0         |
| ❌    |                            | EmailFailed.Update                          | 0.00%    | 5         |
| ❌    |                            | EmailSent.Delete                            | 0.00%    | 1         |
| ✅    |                            | EmailSent.GetBody                           | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetCreatedAt                      | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetError                          | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetID                             | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetRecipient                      | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetRetryAt                        | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetRetryCount                     | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetRetryError                     | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetSentAt                         | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetStatus                         | 0.00%    | 0         |
| ✅    |                            | EmailSent.GetSubject                        | 0.00%    | 0         |
| ❌    |                            | EmailSent.Insert                            | 0.00%    | 1         |
| ❌    |                            | EmailSent.Select                            | 0.00%    | 4         |
| ✅    |                            | EmailSent.SetCreatedAt                      | 0.00%    | 0         |
| ✅    |                            | EmailSent.SetRetryAt                        | 0.00%    | 0         |
| ✅    |                            | EmailSent.SetSentAt                         | 0.00%    | 0         |
| ❌    |                            | EmailSent.Update                            | 0.00%    | 5         |
| ❌    |                            | Migrations.Delete                           | 0.00%    | 1         |
| ✅    |                            | Migrations.GetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.GetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Insert                           | 0.00%    | 1         |
| ❌    |                            | Migrations.Select                           | 0.00%    | 4         |
| ❌    |                            | Migrations.Update                           | 0.00%    | 5         |
| ✅    |                            | NewEmail                                    | 0.00%    | 0         |
| ❌    |                            | QueryConfig.Apply                           | 0.00%    | 13        |
| ✅    |                            | QueryConfig.WithColumns                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithLimit                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithOrderBy                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithStatement                   | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithTable                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithWhere                       | 0.00%    | 0         |
| ✅    |                            | WithColumns                                 | 0.00%    | 0         |
| ✅    |                            | WithLimit                                   | 0.00%    | 0         |
| ✅    |                            | WithOrderBy                                 | 0.00%    | 0         |
| ✅    |                            | WithStatement                               | 0.00%    | 0         |
| ✅    |                            | WithTable                                   | 0.00%    | 0         |
| ✅    |                            | WithWhere                                   | 0.00%    | 0         |
| ✅    | email/schema               | Migrations                                  | 0.00%    | 0         |
| ✅    | email/smtp                 | ConfigFromEnv                               | 0.00%    | 0         |
| ✅    |                            | NewSMTPSender                               | 0.00%    | 0         |
| ❌    |                            | SMTPSender.Send                             | 0.00%    | 3         |
| ❌    |                            | getEnv                                      | 0.00%    | 1         |
| ❌    |                            | getEnvInt                                   | 0.00%    | 3         |
| ✅    | email/storage              | DB                                          | 0.00%    | 0         |
| ❌    |                            | EmailStorage.Create                         | 0.00%    | 1         |
| ❌    |                            | EmailStorage.Get                            | 0.00%    | 1         |
| ❌    |                            | EmailStorage.GetFailed                      | 0.00%    | 1         |
| ❌    |                            | EmailStorage.GetPending                     | 0.00%    | 1         |
| ❌    |                            | EmailStorage.GetSent                        | 0.00%    | 1         |
| ❌    |                            | EmailStorage.Update                         | 0.00%    | 7         |
| ❌    |                            | Migrate                                     | 0.00%    | 2         |
| ✅    |                            | NewEmailStorage                             | 0.00%    | 0         |
| ❌    |                            | NewEmailStorageErr                          | 0.00%    | 1         |
| ❌    | internal/cmd/structurelint | catchAllValue                               | 0.00%    | 10        |
| ❌    |                            | cellValue                                   | 0.00%    | 10        |
| ❌    |                            | columnWidths                                | 0.00%    | 9         |
| ❌    |                            | hasGoFiles                                  | 0.00%    | 5         |
| ❌    |                            | main                                        | 0.00%    | 23        |
| ✅    |                            | pad                                         | 0.00%    | 0         |
| ✅    | maillist                   | MailList.Create                             | 0.00%    | 0         |
| ✅    |                            | MailList.Index                              | 0.00%    | 0         |
| ✅    |                            | MailList.Mount                              | 0.00%    | 0         |
| ✅    |                            | MailList.Name                               | 0.00%    | 0         |
| ✅    |                            | MailList.Start                              | 0.00%    | 0         |
| ✅    |                            | NewMailList                                 | 0.00%    | 0         |
| ✅    | maillist/model             | Maillist.GetCreatedAt                       | 0.00%    | 0         |
| ✅    |                            | Maillist.GetDescription                     | 0.00%    | 0         |
| ✅    |                            | Maillist.GetID                              | 0.00%    | 0         |
| ✅    |                            | Maillist.GetName                            | 0.00%    | 0         |
| ✅    |                            | Maillist.GetSlug                            | 0.00%    | 0         |
| ✅    |                            | Maillist.GetUpdatedAt                       | 0.00%    | 0         |
| ✅    |                            | Maillist.SetCreatedAt                       | 0.00%    | 0         |
| ✅    |                            | Maillist.SetUpdatedAt                       | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetBody                    | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetCreatedAt               | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetID                      | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetMaillistID              | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetScheduledAt             | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetSentAt                  | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetStatus                  | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetSubject                 | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.GetUpdatedAt               | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.SetCreatedAt               | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.SetScheduledAt             | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.SetSentAt                  | 0.00%    | 0         |
| ✅    |                            | MaillistCampaign.SetUpdatedAt               | 0.00%    | 0         |
| ✅    |                            | MaillistMember.GetEmail                     | 0.00%    | 0         |
| ✅    |                            | MaillistMember.GetID                        | 0.00%    | 0         |
| ✅    |                            | MaillistMember.GetMaillistID                | 0.00%    | 0         |
| ✅    |                            | MaillistMember.GetName                      | 0.00%    | 0         |
| ✅    |                            | MaillistMember.GetSource                    | 0.00%    | 0         |
| ✅    |                            | MaillistMember.GetSubscribedAt              | 0.00%    | 0         |
| ✅    |                            | MaillistMember.GetUnsubscribedAt            | 0.00%    | 0         |
| ✅    |                            | MaillistMember.SetSubscribedAt              | 0.00%    | 0         |
| ✅    |                            | MaillistMember.SetUnsubscribedAt            | 0.00%    | 0         |
| ✅    |                            | MaillistMemberActivity.GetActivityAt        | 0.00%    | 0         |
| ✅    |                            | MaillistMemberActivity.GetActivityType      | 0.00%    | 0         |
| ✅    |                            | MaillistMemberActivity.GetCampaignID        | 0.00%    | 0         |
| ✅    |                            | MaillistMemberActivity.GetID                | 0.00%    | 0         |
| ✅    |                            | MaillistMemberActivity.GetMaillistMemberID  | 0.00%    | 0         |
| ✅    |                            | MaillistMemberActivity.SetActivityAt        | 0.00%    | 0         |
| ✅    |                            | Migrations.GetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.GetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatus                        | 0.00%    | 0         |
| ✅    | maillist/storage           | DB                                          | 0.00%    | 0         |
| ❌    |                            | Migrate                                     | 0.00%    | 4         |
| ❌    |                            | NewPermissions                              | 0.00%    | 1         |
| ✅    | pulse                      | NewModule                                   | 0.00%    | 0         |
| ❌    | pulse/client               | Client.EnsureToken                          | 0.00%    | 4         |
| ✅    |                            | Client.LoadToken                            | 76.92%   | 5         |
| ❌    |                            | Client.RefreshToken                         | 0.00%    | 5         |
| ✅    |                            | Client.SaveToken                            | 75.00%   | 4         |
| ❌    |                            | Client.SendPulse                            | 0.00%    | 6         |
| ✅    |                            | Client.ShouldRefresh                        | 100.00%  | 3         |
| ✅    |                            | Client.Token                                | 66.67%   | 1         |
| ✅    |                            | New                                         | 100.00%  | 0         |
| ✅    |                            | configPath                                  | 66.67%   | 2         |
| ❌    | pulse/cmd/pulse            | main                                        | 0.00%    | 1         |
| ❌    |                            | run                                         | 0.00%    | 5         |
| ✅    | pulse/cmd/pulse/login      | NewCommand                                  | 0.00%    | 0         |
| ❌    |                            | Options.Bind                                | 0.00%    | 1         |
| ❌    |                            | Run                                         | 0.00%    | 12        |
| ✅    | pulse/cmd/pulse/record     | NewCommand                                  | 0.00%    | 0         |
| ❌    |                            | Options.Bind                                | 0.00%    | 1         |
| ❌    |                            | Run                                         | 0.00%    | 20        |
| ✅    | pulse/cmd/pulse/register   | NewCommand                                  | 0.00%    | 0         |
| ❌    |                            | Options.Bind                                | 0.00%    | 1         |
| ❌    |                            | Run                                         | 0.00%    | 17        |
| ✅    | pulse/cmd/pulse/server     | NewCommand                                  | 0.00%    | 0         |
| ❌    |                            | Run                                         | 0.00%    | 1         |
| ✅    | pulse/cmd/pulse/version    | NewCommand                                  | 0.00%    | 0         |
| ❌    |                            | Run                                         | 0.00%    | 31        |
| ✅    | pulse/config               | ConfigFS                                    | 0.00%    | 0         |
| ❌    | pulse/model                | Migrations.Delete                           | 0.00%    | 1         |
| ✅    |                            | Migrations.GetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.GetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Insert                           | 0.00%    | 1         |
| ❌    |                            | Migrations.Select                           | 0.00%    | 4         |
| ✅    |                            | Migrations.SetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.SetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.SetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.SetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Update                           | 0.00%    | 5         |
| ❌    |                            | PulseDaily.Delete                           | 0.00%    | 1         |
| ✅    |                            | PulseDaily.GetCount                         | 0.00%    | 0         |
| ✅    |                            | PulseDaily.GetHostname                      | 0.00%    | 0         |
| ✅    |                            | PulseDaily.GetStamp                         | 0.00%    | 0         |
| ✅    |                            | PulseDaily.GetUserID                        | 0.00%    | 0         |
| ❌    |                            | PulseDaily.Insert                           | 0.00%    | 1         |
| ❌    |                            | PulseDaily.Select                           | 0.00%    | 4         |
| ✅    |                            | PulseDaily.SetCount                         | 0.00%    | 0         |
| ✅    |                            | PulseDaily.SetHostname                      | 0.00%    | 0         |
| ✅    |                            | PulseDaily.SetStamp                         | 0.00%    | 0         |
| ✅    |                            | PulseDaily.SetUserID                        | 0.00%    | 0         |
| ❌    |                            | PulseDaily.Update                           | 0.00%    | 5         |
| ❌    |                            | PulseHost.Delete                            | 0.00%    | 1         |
| ✅    |                            | PulseHost.GetCreatedAt                      | 0.00%    | 0         |
| ✅    |                            | PulseHost.GetHostname                       | 0.00%    | 0         |
| ✅    |                            | PulseHost.GetUserID                         | 0.00%    | 0         |
| ❌    |                            | PulseHost.Insert                            | 0.00%    | 1         |
| ❌    |                            | PulseHost.Select                            | 0.00%    | 4         |
| ✅    |                            | PulseHost.SetCreatedAt                      | 0.00%    | 0         |
| ✅    |                            | PulseHost.SetHostname                       | 0.00%    | 0         |
| ✅    |                            | PulseHost.SetUserID                         | 0.00%    | 0         |
| ❌    |                            | PulseHost.Update                            | 0.00%    | 5         |
| ❌    |                            | PulseHourly.Delete                          | 0.00%    | 1         |
| ✅    |                            | PulseHourly.GetCount                        | 0.00%    | 0         |
| ✅    |                            | PulseHourly.GetHostname                     | 0.00%    | 0         |
| ✅    |                            | PulseHourly.GetStamp                        | 0.00%    | 0         |
| ✅    |                            | PulseHourly.GetUserID                       | 0.00%    | 0         |
| ❌    |                            | PulseHourly.Insert                          | 0.00%    | 1         |
| ❌    |                            | PulseHourly.Select                          | 0.00%    | 4         |
| ✅    |                            | PulseHourly.SetCount                        | 0.00%    | 0         |
| ✅    |                            | PulseHourly.SetHostname                     | 0.00%    | 0         |
| ✅    |                            | PulseHourly.SetStamp                        | 0.00%    | 0         |
| ✅    |                            | PulseHourly.SetUserID                       | 0.00%    | 0         |
| ❌    |                            | PulseHourly.Update                          | 0.00%    | 5         |
| ❌    |                            | QueryConfig.Apply                           | 0.00%    | 13        |
| ✅    |                            | QueryConfig.WithColumns                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithLimit                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithOrderBy                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithStatement                   | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithTable                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithWhere                       | 0.00%    | 0         |
| ✅    |                            | WithColumns                                 | 0.00%    | 0         |
| ✅    |                            | WithLimit                                   | 0.00%    | 0         |
| ✅    |                            | WithOrderBy                                 | 0.00%    | 0         |
| ✅    |                            | WithStatement                               | 0.00%    | 0         |
| ✅    |                            | WithTable                                   | 0.00%    | 0         |
| ✅    |                            | WithWhere                                   | 0.00%    | 0         |
| ✅    | pulse/schema               | Migrations                                  | 100.00%  | 0         |
| ✅    | pulse/service              | FS                                          | 0.00%    | 0         |
| ✅    |                            | Handlers.IndexPage                          | 0.00%    | 0         |
| ✅    |                            | Handlers.Mount                              | 0.00%    | 0         |
| ✅    |                            | Handlers.PostIngest                         | 0.00%    | 0         |
| ✅    |                            | Handlers.UserPage                           | 0.00%    | 0         |
| ❌    |                            | Handlers.errorHandler                       | 0.00%    | 1         |
| ❌    |                            | Handlers.indexPage                          | 0.00%    | 4         |
| ❌    |                            | Handlers.postIngest                         | 0.00%    | 3         |
| ❌    |                            | Handlers.userPage                           | 0.00%    | 45        |
| ✅    |                            | NewHandlers                                 | 0.00%    | 0         |
| ✅    |                            | NewPulseModule                              | 0.00%    | 0         |
| ✅    |                            | PulseModule.Mount                           | 0.00%    | 0         |
| ✅    |                            | PulseModule.Name                            | 0.00%    | 0         |
| ❌    |                            | PulseModule.Start                           | 0.00%    | 1         |
| ❌    |                            | PulseModule.setupPulseStorage               | 0.00%    | 2         |
| ❌    |                            | PulseModule.setupStorage                    | 0.00%    | 2         |
| ❌    |                            | PulseModule.setupUserStorage                | 0.00%    | 1         |
| ❌    | pulse/service/keycounter   | KeyboardCounter                             | 0.00%    | 28        |
| ✅    |                            | NewOptions                                  | 0.00%    | 0         |
| ❌    |                            | Options.Flush                               | 0.00%    | 1         |
| ✅    |                            | bytesReader                                 | 0.00%    | 0         |
| ✅    |                            | reader.Read                                 | 0.00%    | 0         |
| ✅    | pulse/storage              | DB                                          | 0.00%    | 0         |
| ✅    |                            | Migrate                                     | 85.71%   | 2         |
| ✅    |                            | NewStorage                                  | 100.00%  | 0         |
| ✅    |                            | Storage.GetUserDaily                        | 80.00%   | 1         |
| ✅    |                            | Storage.GetUserHosts                        | 80.00%   | 1         |
| ❌    |                            | Storage.GetUserHourly                       | 0.00%    | 1         |
| ❌    |                            | Storage.GetUserHourlyAll                    | 0.00%    | 1         |
| ❌    |                            | Storage.GetUserHourlyByHost                 | 0.00%    | 1         |
| ❌    |                            | Storage.ListUserCounts                      | 0.00%    | 1         |
| ❌    |                            | Storage.Pulse                               | 0.00%    | 1         |
| ❌    |                            | Storage.pulseFn                             | 0.00%    | 6         |
| ✅    | pulse/view                 | Templates                                   | 0.00%    | 0         |
| ❌    | user                       | AuthCookie                                  | 0.00%    | 2         |
| ✅    |                            | AuthHeader                                  | 0.00%    | 0         |
| ✅    |                            | AuthOptional                                | 0.00%    | 0         |
| ✅    |                            | AuthQuery                                   | 0.00%    | 0         |
| ✅    |                            | GetSessionUser                              | 100.00%  | 0         |
| ✅    |                            | IsLoggedIn                                  | 0.00%    | 0         |
| ❌    |                            | Middleware.ServeHTTP                        | 0.00%    | 3         |
| ❌    |                            | Middleware.authorizeCookie                  | 0.00%    | 5         |
| ❌    |                            | Middleware.authorizeJWT                     | 0.00%    | 1         |
| ❌    |                            | Middleware.authorizeQuery                   | 0.00%    | 1         |
| ❌    |                            | Middleware.authorizeToken                   | 0.00%    | 8         |
| ❌    |                            | Middleware.authorizeUser                    | 0.00%    | 2         |
| ❌    |                            | Middleware.init                             | 0.00%    | 2         |
| ❌    |                            | Middleware.serveHTTP                        | 0.00%    | 12        |
| ❌    |                            | NewMiddleware                               | 0.00%    | 2         |
| ✅    |                            | NewModule                                   | 0.00%    | 0         |
| ✅    |                            | SetSessionUser                              | 0.00%    | 0         |
| ❌    |                            | SigningKey                                  | 0.00%    | 1         |
| ❌    | user/cmd/generate          | GeneratePlantUML                            | 0.00%    | 4         |
| ❌    |                            | GenerateRego                                | 0.00%    | 2         |
| ❌    |                            | GenerateRouteMap                            | 0.00%    | 5         |
| ❌    |                            | LoadConfig                                  | 0.00%    | 1         |
| ❌    |                            | ValidateFlowsYAML                           | 0.00%    | 7         |
| ❌    |                            | main                                        | 0.00%    | 1         |
| ❌    |                            | start                                       | 0.00%    | 5         |
| ❌    | user/model                 | Migrations.Delete                           | 0.00%    | 1         |
| ✅    |                            | Migrations.GetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.GetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.GetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Insert                           | 0.00%    | 1         |
| ❌    |                            | Migrations.Select                           | 0.00%    | 4         |
| ✅    |                            | Migrations.SetFilename                      | 0.00%    | 0         |
| ✅    |                            | Migrations.SetProject                       | 0.00%    | 0         |
| ✅    |                            | Migrations.SetStatementIndex                | 0.00%    | 0         |
| ✅    |                            | Migrations.SetStatus                        | 0.00%    | 0         |
| ❌    |                            | Migrations.Update                           | 0.00%    | 5         |
| ✅    |                            | NewUser                                     | 100.00%  | 0         |
| ✅    |                            | NewUserGroup                                | 0.00%    | 0         |
| ✅    |                            | NewUserSession                              | 0.00%    | 0         |
| ❌    |                            | QueryConfig.Apply                           | 0.00%    | 13        |
| ✅    |                            | QueryConfig.WithColumns                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithLimit                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithOrderBy                     | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithStatement                   | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithTable                       | 0.00%    | 0         |
| ✅    |                            | QueryConfig.WithWhere                       | 0.00%    | 0         |
| ✅    |                            | TransportJSON                               | 75.00%   | 1         |
| ❌    |                            | User.Delete                                 | 0.00%    | 1         |
| ✅    |                            | User.GetCreatedAt                           | 0.00%    | 0         |
| ✅    |                            | User.GetDeletedAt                           | 0.00%    | 0         |
| ✅    |                            | User.GetFullName                            | 0.00%    | 0         |
| ✅    |                            | User.GetID                                  | 0.00%    | 0         |
| ✅    |                            | User.GetSlug                                | 0.00%    | 0         |
| ✅    |                            | User.GetUpdatedAt                           | 0.00%    | 0         |
| ✅    |                            | User.GetUsername                            | 0.00%    | 0         |
| ❌    |                            | User.Insert                                 | 0.00%    | 1         |
| ✅    |                            | User.Ok                                     | 100.00%  | 0         |
| ❌    |                            | User.Select                                 | 0.00%    | 4         |
| ✅    |                            | User.SetCreatedAt                           | 0.00%    | 0         |
| ✅    |                            | User.SetDeletedAt                           | 100.00%  | 0         |
| ✅    |                            | User.SetFullName                            | 0.00%    | 0         |
| ✅    |                            | User.SetID                                  | 0.00%    | 0         |
| ✅    |                            | User.SetSlug                                | 0.00%    | 0         |
| ✅    |                            | User.SetUpdatedAt                           | 0.00%    | 0         |
| ✅    |                            | User.SetUsername                            | 0.00%    | 0         |
| ✅    |                            | User.String                                 | 100.00%  | 1         |
| ❌    |                            | User.Update                                 | 0.00%    | 5         |
| ✅    |                            | User.Validate                               | 28.57%   | 3         |
| ❌    |                            | UserAuth.Delete                             | 0.00%    | 1         |
| ✅    |                            | UserAuth.GetActivatedAt                     | 0.00%    | 0         |
| ✅    |                            | UserAuth.GetActivationSentAt                | 0.00%    | 0         |
| ✅    |                            | UserAuth.GetActivationToken                 | 0.00%    | 0         |
| ✅    |                            | UserAuth.GetCreatedAt                       | 0.00%    | 0         |
| ✅    |                            | UserAuth.GetEmail                           | 0.00%    | 0         |
| ✅    |                            | UserAuth.GetPassword                        | 0.00%    | 0         |
| ✅    |                            | UserAuth.GetUpdatedAt                       | 0.00%    | 0         |
| ✅    |                            | UserAuth.GetUserID                          | 0.00%    | 0         |
| ❌    |                            | UserAuth.Insert                             | 0.00%    | 1         |
| ❌    |                            | UserAuth.Select                             | 0.00%    | 4         |
| ✅    |                            | UserAuth.SetActivatedAt                     | 0.00%    | 0         |
| ✅    |                            | UserAuth.SetActivationSentAt                | 0.00%    | 0         |
| ✅    |                            | UserAuth.SetActivationToken                 | 0.00%    | 0         |
| ✅    |                            | UserAuth.SetCreatedAt                       | 0.00%    | 0         |
| ✅    |                            | UserAuth.SetEmail                           | 0.00%    | 0         |
| ✅    |                            | UserAuth.SetPassword                        | 0.00%    | 0         |
| ✅    |                            | UserAuth.SetUpdatedAt                       | 0.00%    | 0         |
| ✅    |                            | UserAuth.SetUserID                          | 0.00%    | 0         |
| ❌    |                            | UserAuth.Update                             | 0.00%    | 5         |
| ❌    |                            | UserAuth.Valid                              | 0.00%    | 2         |
| ❌    |                            | UserCreateRequest.User                      | 0.00%    | 1         |
| ✅    |                            | UserCreateRequest.UserAuth                  | 0.00%    | 0         |
| ✅    |                            | UserCreateRequest.Valid                     | 88.89%   | 5         |
| ✅    |                            | UserCreateRequest.ValidateUsername          | 100.00%  | 4         |
| ❌    |                            | UserGroup.Delete                            | 0.00%    | 1         |
| ✅    |                            | UserGroup.GetCreatedAt                      | 0.00%    | 0         |
| ✅    |                            | UserGroup.GetID                             | 0.00%    | 0         |
| ✅    |                            | UserGroup.GetTitle                          | 0.00%    | 0         |
| ✅    |                            | UserGroup.GetUpdatedAt                      | 0.00%    | 0         |
| ❌    |                            | UserGroup.Insert                            | 0.00%    | 1         |
| ❌    |                            | UserGroup.Select                            | 0.00%    | 4         |
| ✅    |                            | UserGroup.SetCreatedAt                      | 0.00%    | 0         |
| ✅    |                            | UserGroup.SetID                             | 0.00%    | 0         |
| ✅    |                            | UserGroup.SetTitle                          | 0.00%    | 0         |
| ✅    |                            | UserGroup.SetUpdatedAt                      | 0.00%    | 0         |
| ✅    |                            | UserGroup.String                            | 0.00%    | 0         |
| ❌    |                            | UserGroup.Update                            | 0.00%    | 5         |
| ❌    |                            | UserGroupMember.Delete                      | 0.00%    | 1         |
| ✅    |                            | UserGroupMember.GetJoinedAt                 | 0.00%    | 0         |
| ✅    |                            | UserGroupMember.GetUserGroupID              | 0.00%    | 0         |
| ✅    |                            | UserGroupMember.GetUserID                   | 0.00%    | 0         |
| ❌    |                            | UserGroupMember.Insert                      | 0.00%    | 1         |
| ❌    |                            | UserGroupMember.Select                      | 0.00%    | 4         |
| ✅    |                            | UserGroupMember.SetJoinedAt                 | 0.00%    | 0         |
| ✅    |                            | UserGroupMember.SetUserGroupID              | 0.00%    | 0         |
| ✅    |                            | UserGroupMember.SetUserID                   | 0.00%    | 0         |
| ❌    |                            | UserGroupMember.Update                      | 0.00%    | 5         |
| ❌    |                            | UserPasskey.Delete                          | 0.00%    | 1         |
| ✅    |                            | UserPasskey.GetAttestationType              | 0.00%    | 0         |
| ✅    |                            | UserPasskey.GetCreatedAt                    | 0.00%    | 0         |
| ✅    |                            | UserPasskey.GetCredentialID                 | 0.00%    | 0         |
| ✅    |                            | UserPasskey.GetID                           | 0.00%    | 0         |
| ✅    |                            | UserPasskey.GetPublicKey                    | 0.00%    | 0         |
| ✅    |                            | UserPasskey.GetSignCount                    | 0.00%    | 0         |
| ✅    |                            | UserPasskey.GetTransport                    | 0.00%    | 0         |
| ✅    |                            | UserPasskey.GetUserID                       | 0.00%    | 0         |
| ❌    |                            | UserPasskey.Insert                          | 0.00%    | 1         |
| ❌    |                            | UserPasskey.Select                          | 0.00%    | 4         |
| ✅    |                            | UserPasskey.SetAttestationType              | 0.00%    | 0         |
| ✅    |                            | UserPasskey.SetCreatedAt                    | 0.00%    | 0         |
| ✅    |                            | UserPasskey.SetCredentialID                 | 0.00%    | 0         |
| ✅    |                            | UserPasskey.SetID                           | 0.00%    | 0         |
| ✅    |                            | UserPasskey.SetPublicKey                    | 0.00%    | 0         |
| ✅    |                            | UserPasskey.SetSignCount                    | 0.00%    | 0         |
| ✅    |                            | UserPasskey.SetTransport                    | 0.00%    | 0         |
| ✅    |                            | UserPasskey.SetUserID                       | 0.00%    | 0         |
| ✅    |                            | UserPasskey.ToCredential                    | 100.00%  | 0         |
| ❌    |                            | UserPasskey.Update                          | 0.00%    | 5         |
| ❌    |                            | UserSession.Delete                          | 0.00%    | 1         |
| ✅    |                            | UserSession.GetCreatedAt                    | 0.00%    | 0         |
| ✅    |                            | UserSession.GetExpiresAt                    | 0.00%    | 0         |
| ✅    |                            | UserSession.GetID                           | 0.00%    | 0         |
| ✅    |                            | UserSession.GetUserID                       | 0.00%    | 0         |
| ❌    |                            | UserSession.Insert                          | 0.00%    | 1         |
| ✅    |                            | UserSession.Ok                              | 0.00%    | 0         |
| ❌    |                            | UserSession.Select                          | 0.00%    | 4         |
| ✅    |                            | UserSession.SetCreatedAt                    | 0.00%    | 0         |
| ✅    |                            | UserSession.SetExpiresAt                    | 0.00%    | 0         |
| ✅    |                            | UserSession.SetID                           | 0.00%    | 0         |
| ✅    |                            | UserSession.SetUserID                       | 0.00%    | 0         |
| ❌    |                            | UserSession.Update                          | 0.00%    | 5         |
| ❌    |                            | UserSession.Validate                        | 0.00%    | 5         |
| ❌    |                            | UserTokenRevoked.Delete                     | 0.00%    | 1         |
| ✅    |                            | UserTokenRevoked.GetCreatedAt               | 0.00%    | 0         |
| ✅    |                            | UserTokenRevoked.GetExpiresAt               | 0.00%    | 0         |
| ✅    |                            | UserTokenRevoked.GetJti                     | 0.00%    | 0         |
| ✅    |                            | UserTokenRevoked.GetUserID                  | 0.00%    | 0         |
| ❌    |                            | UserTokenRevoked.Insert                     | 0.00%    | 1         |
| ❌    |                            | UserTokenRevoked.Select                     | 0.00%    | 4         |
| ✅    |                            | UserTokenRevoked.SetCreatedAt               | 0.00%    | 0         |
| ✅    |                            | UserTokenRevoked.SetExpiresAt               | 0.00%    | 0         |
| ✅    |                            | UserTokenRevoked.SetJti                     | 0.00%    | 0         |
| ✅    |                            | UserTokenRevoked.SetUserID                  | 0.00%    | 0         |
| ❌    |                            | UserTokenRevoked.Update                     | 0.00%    | 5         |
| ✅    |                            | WebAuthnUser.WebAuthnCredentials            | 100.00%  | 1         |
| ✅    |                            | WebAuthnUser.WebAuthnDisplayName            | 100.00%  | 0         |
| ✅    |                            | WebAuthnUser.WebAuthnID                     | 100.00%  | 0         |
| ✅    |                            | WebAuthnUser.WebAuthnName                   | 100.00%  | 0         |
| ✅    |                            | WithColumns                                 | 0.00%    | 0         |
| ✅    |                            | WithLimit                                   | 0.00%    | 0         |
| ✅    |                            | WithOrderBy                                 | 0.00%    | 0         |
| ✅    |                            | WithStatement                               | 0.00%    | 0         |
| ✅    |                            | WithTable                                   | 0.00%    | 0         |
| ✅    |                            | WithWhere                                   | 0.00%    | 0         |
| ✅    | user/schema                | Migrations                                  | 0.00%    | 0         |
| ✅    | user/service               | FS                                          | 0.00%    | 0         |
| ✅    |                            | NewUserModule                               | 0.00%    | 0         |
| ✅    |                            | UserModule.Mount                            | 0.00%    | 0         |
| ✅    |                            | UserModule.Name                             | 0.00%    | 0         |
| ❌    |                            | UserModule.Start                            | 0.00%    | 7         |
| ✅    | user/service/api           | Handlers.ActivateEmail                      | 0.00%    | 0         |
| ✅    |                            | Handlers.CreateToken                        | 100.00%  | 0         |
| ✅    |                            | Handlers.Mount                              | 0.00%    | 0         |
| ✅    |                            | Handlers.PasskeyLoginBegin                  | 0.00%    | 0         |
| ✅    |                            | Handlers.PasskeyLoginFinish                 | 100.00%  | 0         |
| ✅    |                            | Handlers.PasskeyRegisterBegin               | 100.00%  | 0         |
| ✅    |                            | Handlers.PasskeyRegisterFinish              | 100.00%  | 0         |
| ✅    |                            | Handlers.RefreshToken                       | 100.00%  | 0         |
| ✅    |                            | Handlers.Register                           | 100.00%  | 0         |
| ✅    |                            | Handlers.ResendActivation                   | 0.00%    | 0         |
| ✅    |                            | Handlers.RevokeToken                        | 100.00%  | 0         |
| ❌    |                            | Handlers.activateEmail                      | 0.00%    | 6         |
| ❌    |                            | Handlers.createToken                        | 18.52%   | 8         |
| ✅    |                            | Handlers.errorHandler                       | 80.00%   | 2         |
| ❌    |                            | Handlers.passkeyLoginBegin                  | 0.00%    | 1         |
| ✅    |                            | Handlers.passkeyLoginFinish                 | 20.00%   | 3         |
| ❌    |                            | Handlers.passkeyRegisterBegin               | 35.71%   | 6         |
| ✅    |                            | Handlers.passkeyRegisterFinish              | 20.00%   | 3         |
| ❌    |                            | Handlers.refreshToken                       | 70.83%   | 11        |
| ❌    |                            | Handlers.register                           | 50.00%   | 7         |
| ❌    |                            | Handlers.registerActivated                  | 0.00%    | 2         |
| ❌    |                            | Handlers.registerPending                    | 0.00%    | 2         |
| ❌    |                            | Handlers.resendActivation                   | 0.00%    | 3         |
| ✅    |                            | Handlers.revokeToken                        | 57.14%   | 5         |
| ✅    |                            | NewHandlers                                 | 100.00%  | 2         |
| ✅    |                            | RequestError.Error                          | 0.00%    | 0         |
| ❌    |                            | mapRegisterError                            | 0.00%    | 3         |
| ✅    | user/service/auth          | JWT.Claims                                  | 84.00%   | 14        |
| ✅    |                            | JWT.Create                                  | 100.00%  | 0         |
| ✅    |                            | JWT.CreateWithJTI                           | 95.00%   | 1         |
| ✅    |                            | JWT.IsUser                                  | 100.00%  | 0         |
| ✅    |                            | JWT.UserID                                  | 75.00%   | 1         |
| ✅    |                            | JWT.Validate                                | 75.00%   | 1         |
| ✅    |                            | NewJWT                                      | 100.00%  | 0         |
| ✅    | user/service/passkey       | Error.Error                                 | 100.00%  | 0         |
| ✅    |                            | New                                         | 100.00%  | 0         |
| ✅    |                            | Service.BeginLogin                          | 94.44%   | 1         |
| ✅    |                            | Service.BeginRegistration                   | 95.83%   | 1         |
| ❌    |                            | Service.FinishLogin                         | 12.00%   | 7         |
| ✅    |                            | Service.FinishRegistration                  | 11.11%   | 4         |
| ✅    |                            | Service.consumeSession                      | 80.00%   | 3         |
| ✅    | user/service/web           | Handlers.Error                              | 83.33%   | 2         |
| ✅    |                            | Handlers.GetError                           | 100.00%  | 0         |
| ✅    |                            | Handlers.Login                              | 0.00%    | 0         |
| ✅    |                            | Handlers.LoginView                          | 0.00%    | 0         |
| ✅    |                            | Handlers.Logout                             | 0.00%    | 0         |
| ✅    |                            | Handlers.LogoutView                         | 0.00%    | 0         |
| ✅    |                            | Handlers.Mount                              | 0.00%    | 0         |
| ✅    |                            | Handlers.Register                           | 100.00%  | 0         |
| ✅    |                            | Handlers.RegisterView                       | 100.00%  | 0         |
| ✅    |                            | Handlers.errorHandler                       | 25.00%   | 1         |
| ❌    |                            | Handlers.login                              | 0.00%    | 5         |
| ❌    |                            | Handlers.loginView                          | 0.00%    | 9         |
| ❌    |                            | Handlers.logout                             | 0.00%    | 2         |
| ✅    |                            | Handlers.logoutView                         | 0.00%    | 0         |
| ❌    |                            | Handlers.register                           | 51.06%   | 6         |
| ✅    |                            | Handlers.registerView                       | 100.00%  | 0         |
| ✅    |                            | NewHandlers                                 | 100.00%  | 0         |
| ✅    |                            | NewRenderer                                 | 100.00%  | 1         |
| ✅    |                            | Renderer.Load                               | 100.00%  | 0         |
| ✅    |                            | Renderer.Login                              | 100.00%  | 0         |
| ✅    |                            | Renderer.Logout                             | 100.00%  | 0         |
| ✅    |                            | Renderer.Register                           | 100.00%  | 0         |
| ✅    | user/storage               | DB                                          | 0.00%    | 0         |
| ❌    |                            | Migrate                                     | 0.00%    | 2         |
| ✅    |                            | NewPasskeyStorage                           | 100.00%  | 0         |
| ✅    |                            | NewRevokedTokenStorage                      | 0.00%    | 0         |
| ✅    |                            | NewSessionStorage                           | 0.00%    | 0         |
| ✅    |                            | NewUserStorage                              | 100.00%  | 0         |
| ❌    |                            | NewUserStorageErr                           | 0.00%    | 1         |
| ❌    |                            | PasskeyStorage.Create                       | 0.00%    | 1         |
| ❌    |                            | PasskeyStorage.Delete                       | 0.00%    | 1         |
| ❌    |                            | PasskeyStorage.GetByCredentialID            | 0.00%    | 1         |
| ❌    |                            | PasskeyStorage.ListByUser                   | 0.00%    | 1         |
| ❌    |                            | PasskeyStorage.UpdateSignCount              | 0.00%    | 1         |
| ❌    |                            | RevokedTokenStorage.IsRevoked               | 0.00%    | 4         |
| ❌    |                            | RevokedTokenStorage.PurgeExpired            | 0.00%    | 3         |
| ❌    |                            | RevokedTokenStorage.Revoke                  | 0.00%    | 4         |
| ❌    |                            | SessionStorage.Create                       | 0.00%    | 1         |
| ❌    |                            | SessionStorage.Delete                       | 0.00%    | 1         |
| ❌    |                            | SessionStorage.Get                          | 0.00%    | 4         |
| ❌    |                            | UserStorage.Activate                        | 0.00%    | 4         |
| ❌    |                            | UserStorage.Authenticate                    | 0.00%    | 8         |
| ❌    |                            | UserStorage.Create                          | 0.00%    | 5         |
| ❌    |                            | UserStorage.CreatePending                   | 0.00%    | 6         |
| ❌    |                            | UserStorage.Get                             | 0.00%    | 1         |
| ❌    |                            | UserStorage.GetByStub                       | 0.00%    | 1         |
| ❌    |                            | UserStorage.GetByUsername                   | 0.00%    | 1         |
| ❌    |                            | UserStorage.GetGroups                       | 0.00%    | 1         |
| ❌    |                            | UserStorage.IsActivated                     | 0.00%    | 1         |
| ❌    |                            | UserStorage.List                            | 0.00%    | 1         |
| ❌    |                            | UserStorage.ResetActivation                 | 0.00%    | 4         |
| ❌    |                            | UserStorage.Update                          | 0.00%    | 1         |
| ❌    |                            | UserStorage.activationBody                  | 0.00%    | 1         |
| ❌    |                            | UserStorage.insertUserAndAuth               | 0.00%    | 12        |
| ❌    |                            | UserStorage.sendActivationEmail             | 0.00%    | 1         |
| ❌    |                            | UserStorage.validateCreate                  | 0.00%    | 5         |
| ✅    |                            | init                                        | 100.00%  | 0         |
| ❌    |                            | newActivationToken                          | 0.00%    | 1         |
| ✅    | user/view                  | Templates                                   | 100.00%  | 0         |

