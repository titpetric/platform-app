# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be covered by invocation.

The storage packages have integration tests.

## Packages

| Status | Package                                                | Coverage | Cognitive | Lines |
| ------ | ------------------------------------------------------ | -------- | --------- | ----- |
| ✅      | titpetric/platform-app                      | 0.00%    | 0         | 0     |
| ✅      | titpetric/platform-app/autoload             | 100.00%  | 0         | 3     |
| ✅      | titpetric/platform-app/cmd                  | 100.00%  | 0         | 21    |
| ✅      | titpetric/platform-app/cmd/app              | 0.00%    | 0         | 3     |
| ✅      | titpetric/platform-app/internal             | 80.00%   | 2         | 32    |
| ✅      | titpetric/platform-app/modules/expvar       | 95.23%   | 1         | 19    |
| ✅      | titpetric/platform-app/modules/theme        | 100.00%  | 0         | 11    |
| ❌      | titpetric/platform-app/modules/user         | 54.76%   | 13        | 104   |
| ✅      | titpetric/platform-app/modules/user/model   | 10.53%   | 3         | 88    |
| ❌      | titpetric/platform-app/modules/user/service | 22.28%   | 27        | 273   |
| ❌      | titpetric/platform-app/modules/user/storage | 44.01%   | 24        | 211   |

## Functions

| Status | Package                                                | Function                       | Coverage | Cognitive |
| ------ | ------------------------------------------------------ | ------------------------------ | -------- | --------- |
| ✅      | titpetric/platform-app/autoload             | init                           | 100.00%  | 0         |
| ✅      | titpetric/platform-app/cmd                  | Main                           | 100.00%  | 0         |
| ✅      | titpetric/platform-app/cmd                  | Register                       | 100.00%  | 0         |
| ✅      | titpetric/platform-app/cmd/app              | main                           | 0.00%    | 0         |
| ✅      | titpetric/platform-app/internal             | ContextValue[T].Get            | 100.00%  | 1         |
| ✅      | titpetric/platform-app/internal             | ContextValue[T].Set            | 100.00%  | 0         |
| ✅      | titpetric/platform-app/internal             | NewContextValue                | 100.00%  | 0         |
| ✅      | titpetric/platform-app/internal             | NewTemplate                    | 100.00%  | 0         |
| ❌      | titpetric/platform-app/internal             | Template.Render                | 0.00%    | 1         |
| ✅      | titpetric/platform-app/modules/expvar       | Handler.Mount                  | 100.00%  | 0         |
| ✅      | titpetric/platform-app/modules/expvar       | Handler.Start                  | 85.70%   | 1         |
| ✅      | titpetric/platform-app/modules/expvar       | NewHandler                     | 100.00%  | 0         |
| ✅      | titpetric/platform-app/modules/theme        | NewOptions                     | 100.00%  | 0         |
| ❌      | titpetric/platform-app/modules/user         | GetSessionUser                 | 0.00%    | 8         |
| ✅      | titpetric/platform-app/modules/user         | Handler.Mount                  | 100.00%  | 0         |
| ✅      | titpetric/platform-app/modules/user         | Handler.Name                   | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user         | Handler.Start                  | 83.30%   | 2         |
| ✅      | titpetric/platform-app/modules/user         | Handler.Stop                   | 100.00%  | 0         |
| ❌      | titpetric/platform-app/modules/user         | IsLoggedIn                     | 0.00%    | 3         |
| ✅      | titpetric/platform-app/modules/user         | NewHandler                     | 100.00%  | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | NewUser                        | 100.00%  | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | NewUserGroup                   | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.GetCreatedAt              | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.GetDeletedAt              | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.GetFirstName              | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.GetID                     | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.GetLastName               | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.GetUpdatedAt              | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.IsActive                  | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.SetCreatedAt              | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.SetDeletedAt              | 100.00%  | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.SetUpdatedAt              | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | User.String                    | 100.00%  | 1         |
| ✅      | titpetric/platform-app/modules/user/model   | UserAuth.GetCreatedAt          | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserAuth.GetEmail              | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserAuth.GetPassword           | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserAuth.GetUpdatedAt          | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserAuth.GetUserID             | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserAuth.SetCreatedAt          | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserAuth.SetUpdatedAt          | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserAuth.Valid                 | 100.00%  | 2         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroup.GetCreatedAt         | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroup.GetID                | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroup.GetTitle             | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroup.GetUpdatedAt         | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroup.SetCreatedAt         | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroup.SetUpdatedAt         | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroup.String               | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroupMember.GetJoinedAt    | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroupMember.GetUserGroupID | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroupMember.GetUserID      | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserGroupMember.SetJoinedAt    | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserSession.GetCreatedAt       | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserSession.GetExpiresAt       | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserSession.GetID              | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserSession.GetUserID          | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserSession.SetCreatedAt       | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/model   | UserSession.SetExpiresAt       | 0.00%    | 0         |
| ✅      | titpetric/platform-app/modules/user/service | NewService                     | 75.00%   | 1         |
| ✅      | titpetric/platform-app/modules/user/service | Service.Close                  | 100.00%  | 0         |
| ❌      | titpetric/platform-app/modules/user/service | Service.Error                  | 0.00%    | 2         |
| ✅      | titpetric/platform-app/modules/user/service | Service.GetError               | 0.00%    | 0         |
| ❌      | titpetric/platform-app/modules/user/service | Service.Login                  | 0.00%    | 5         |
| ❌      | titpetric/platform-app/modules/user/service | Service.LoginView              | 0.00%    | 8         |
| ❌      | titpetric/platform-app/modules/user/service | Service.Logout                 | 0.00%    | 2         |
| ✅      | titpetric/platform-app/modules/user/service | Service.LogoutView             | 0.00%    | 0         |
| ❌      | titpetric/platform-app/modules/user/service | Service.Register               | 0.00%    | 4         |
| ✅      | titpetric/platform-app/modules/user/service | Service.RegisterView           | 0.00%    | 0         |
| ❌      | titpetric/platform-app/modules/user/service | Service.View                   | 0.00%    | 3         |
| ✅      | titpetric/platform-app/modules/user/service | Service.initTemplates          | 92.30%   | 2         |
| ✅      | titpetric/platform-app/modules/user/storage | DB                             | 100.00%  | 0         |
| ✅      | titpetric/platform-app/modules/user/storage | NewSessionStorage              | 100.00%  | 0         |
| ✅      | titpetric/platform-app/modules/user/storage | NewUserStorage                 | 100.00%  | 0         |
| ❌      | titpetric/platform-app/modules/user/storage | SessionStorage.Create          | 0.00%    | 1         |
| ✅      | titpetric/platform-app/modules/user/storage | SessionStorage.Delete          | 85.70%   | 1         |
| ✅      | titpetric/platform-app/modules/user/storage | SessionStorage.Get             | 63.60%   | 4         |
| ❌      | titpetric/platform-app/modules/user/storage | UserStorage.Authenticate       | 34.80%   | 8         |
| ❌      | titpetric/platform-app/modules/user/storage | UserStorage.Create             | 0.00%    | 7         |
| ❌      | titpetric/platform-app/modules/user/storage | UserStorage.Get                | 0.00%    | 1         |
| ❌      | titpetric/platform-app/modules/user/storage | UserStorage.GetGroups          | 0.00%    | 1         |
| ❌      | titpetric/platform-app/modules/user/storage | UserStorage.Update             | 0.00%    | 1         |

