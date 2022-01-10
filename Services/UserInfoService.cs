using Microsoft.Extensions.Logging;
using System.Threading.Tasks;
using System.Text.Json;
using System.Collections.Generic;
using Microsoft.Data.SqlClient;
using Microsoft.Identity.Client;

namespace aspnet_webapp.Services {
    public interface IUserInfoService
    {
        Task<UserInfo> GetUserInfo();
    }

    public class UserInfo
    {
        public string Name { get; set; }

    }

    public class UserInfoService : IUserInfoService
    {
        private readonly ILogger _logger;
        

        public UserInfoService(ILogger<UserInfoService> logger)
        {
            _logger = logger;
        }

        public async Task<UserInfo> GetUserInfo()
        {
            await Task.CompletedTask;
            return new UserInfo {Name="Nils"};
        }
    }
}